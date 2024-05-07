use async_channel::{Receiver, Sender};
use futures::prelude::*;
use libp2p::{
    gossipsub, gossipsub::IdentTopic, identify, identity, kad, kad::Mode, noise,
    swarm::NetworkBehaviour, swarm::SwarmEvent, tcp, yamux, PeerId,
};
use ramd_db::{keys::RAMD_P2P_KEYPAIR_KEY, storage::Storage};
use ramd_p2p_types::message::P2pMessage;
use std::{
    collections::hash_map::DefaultHasher,
    hash::{Hash, Hasher},
    sync::Arc,
    time::Duration,
};
use tracing::{debug, error, info, warn};

use crate::config::{P2pConfig, RAM_PROTOCOL_VERSION};

#[derive(NetworkBehaviour)]
struct RamdBehavior {
    gossipsub: gossipsub::Behaviour,
    kademlia: kad::Behaviour<kad::store::MemoryStore>,
    identify: identify::Behaviour,
}

pub struct Server<S>
where
    S: Storage<Vec<u8>, Vec<u8>>,
{
    _storage: Arc<S>,
    swarm: libp2p::Swarm<RamdBehavior>,
    boot_nodes: Vec<PeerId>,
    topic: IdentTopic,
    max_peers_limit: usize,
    msg_receiver: Receiver<P2pMessage>,
    bootstrap_interval: tokio::time::Interval,
}

impl<S> Server<S>
where
    S: Storage<Vec<u8>, Vec<u8>>,
{
    pub fn new(p2p_cfg: &P2pConfig, storage: Arc<S>) -> eyre::Result<(Self, Sender<P2pMessage>)> {
        let node_key = Self::get_node_key(storage.clone())?;

        let mut swarm = libp2p::SwarmBuilder::with_existing_identity(node_key)
            .with_tokio()
            .with_tcp(
                tcp::Config::default(),
                noise::Config::new,
                yamux::Config::default,
            )?
            .with_dns()?
            .with_behaviour(|key| {
                // To content-address message, we can take the hash of message and use it as an ID.
                let message_id_fn = |message: &gossipsub::Message| {
                    let mut s = DefaultHasher::new();
                    message.data.hash(&mut s);
                    gossipsub::MessageId::from(s.finish().to_string())
                };

                // Configure gossipsub behavior
                let gossipsub_config = gossipsub::ConfigBuilder::default()
                    .heartbeat_interval(Duration::from_secs(10)) // This is set to aid debugging by not cluttering the log space
                    .validation_mode(gossipsub::ValidationMode::Strict)
                    .message_id_fn(message_id_fn)
                    .build()
                    .map_err(|msg| std::io::Error::new(std::io::ErrorKind::Other, msg))?;

                let gossipsub = gossipsub::Behaviour::new(
                    gossipsub::MessageAuthenticity::Signed(key.clone()),
                    gossipsub_config,
                )?;

                // Configure kademlia behavior
                let peer_id = key.public().to_peer_id();
                let kademlia = kad::Behaviour::new(peer_id, kad::store::MemoryStore::new(peer_id));

                // Configure identify protocol so that this node cane be discovered
                let identify = identify::Behaviour::new(identify::Config::new(
                    RAM_PROTOCOL_VERSION.to_string(),
                    key.public(),
                ));

                Ok(RamdBehavior {
                    gossipsub,
                    kademlia,
                    identify,
                })
            })?
            .with_swarm_config(|c| {
                c.with_idle_connection_timeout(p2p_cfg.idle_connection_timeout())
            })
            .build();

        // Subscribe to configured topic
        let topic = IdentTopic::new(&p2p_cfg.topic);
        swarm.behaviour_mut().gossipsub.subscribe(&topic)?;

        // Adding boot node addresses for initial peer discovery
        let boot_nodes = p2p_cfg
            .boot_nodes
            .iter()
            .map(|boot| {
                let peer_id = boot.peer_id.parse()?;
                let address = boot.address.parse()?;

                swarm
                    .behaviour_mut()
                    .kademlia
                    .add_address(&peer_id, address);

                Ok(peer_id)
            })
            .collect::<eyre::Result<Vec<_>>>()?;

        // Try to deal with peers from config
        if let Some(known_peers) = p2p_cfg.peer_addresses()? {
            for known_peer in known_peers.into_iter() {
                if let Err(e) = swarm.dial(known_peer) {
                    error!(
                        target: "p2p",
                        "Failed to dial with provided peer. Reason: {}",
                        e.to_string()
                    );
                }
            }
        }

        swarm.behaviour_mut().kademlia.set_mode(Some(Mode::Server));
        swarm.listen_on(format!("/ip4/0.0.0.0/tcp/{}", p2p_cfg.port).parse()?)?;

        // Create channel for communicating with p2p module
        let (msg_sender, msg_receiver) = async_channel::unbounded();

        // Create bootstrap interval timer
        let bootstrap_interval = tokio::time::interval(std::time::Duration::from_secs(
            p2p_cfg.bootstrap_interval_secs,
        ));

        Ok((
            Self {
                _storage: storage,
                swarm,
                boot_nodes,
                topic,
                max_peers_limit: p2p_cfg.max_peers_limit,
                msg_receiver,
                bootstrap_interval,
            },
            msg_sender,
        ))
    }

    pub async fn launch(&mut self) {
        loop {
            tokio::select! {
                // periodically running bootstrap to refresh peers dht
                _ = self.bootstrap_interval.tick() => {
                    // as long as peer limit is reached we do not need to search for more peers
                    if self.is_peer_limit_reached() {
                        continue;
                    }

                    if let Err(e) = self.swarm.behaviour_mut().kademlia.bootstrap() {
                        error!(target: "p2p", "Bootstrap step has failed, waiting for next iteration. Reason: {}", e.to_string());
                    }
                }
                // ramd request for broadcasting a message
                Ok(ramd_msg) = self.msg_receiver.recv() => {
                    let Ok(msg) = serde_json::to_string(&ramd_msg) else {
                        error!(target: "p2p", "Failed to serialize P2pMessage struct. Received message: {:?}", ramd_msg);
                        continue;
                    };

                    // Try to broadcast message to connected nodes
                    if let Err(e) = self.swarm.behaviour_mut().gossipsub.publish(
                        self.topic.clone(),
                        msg.as_bytes(), // TODO: for now just forward the message
                    ) {
                        error!("Failed to broadcast due to: {}", e.to_string());
                    }
                }
                // libp2p related events
                event = self.swarm.select_next_some() => match event {
                    // Event from local server, logging locally assigned
                    SwarmEvent::NewListenAddr { address, .. } => {
                        info!(target: "p2p", "One of our listeners has reported a new local listening address. Listening on {address:?}");
                    }
                    SwarmEvent::ConnectionEstablished { peer_id, .. } => {
                        info!(target: "p2p", "Connection established with peer: {}", peer_id);

                        // if peers limit reached on a new connection established event - disconnect from it
                        if self.is_peer_limit_reached() {
                            debug!("Peer limit is reached. Disconnecting from peer {peer_id}");
                            self.disconnect_peer(&peer_id);
                        }
                    }
                    SwarmEvent::ConnectionClosed { peer_id, .. } => {
                        info!(target: "p2p", "Connection was closed with peer: {}", peer_id);
                    }
                    // Handle kademlia behavior events
                    SwarmEvent::Behaviour(RamdBehaviorEvent::Kademlia(kad::Event::OutboundQueryProgressed { result, ..})) => {
                        match result {
                            kad::QueryResult::GetProviders(Ok(kad::GetProvidersOk::FoundProviders { key, providers, .. })) => {
                                for peer in providers {
                                    debug!(
                                        "KAD: Peer {peer:?} provides key {:?}",
                                        std::str::from_utf8(key.as_ref()).unwrap()
                                    );
                                }
                            }
                            kad::QueryResult::GetProviders(Err(err)) => {
                                debug!("KAD: Failed to get providers: {err:?}");
                            }
                            kad::QueryResult::GetRecord(Ok(
                                kad::GetRecordOk::FoundRecord(kad::PeerRecord {
                                    record: kad::Record { key, value, .. },
                                    ..
                                })
                            )) => {
                                debug!(
                                    "KAD: Got record {:?} {:?}",
                                    std::str::from_utf8(key.as_ref()).unwrap(),
                                    std::str::from_utf8(&value).unwrap(),
                                );
                            }
                            kad::QueryResult::GetRecord(Ok(_)) => {}
                            kad::QueryResult::GetRecord(Err(err)) => {
                                debug!("KAD: Failed to get record: {err:?}");
                            }
                            kad::QueryResult::PutRecord(Ok(kad::PutRecordOk { key })) => {
                                debug!(
                                    "KAD: Successfully put record {:?}",
                                    std::str::from_utf8(key.as_ref()).unwrap()
                                );
                            }
                            kad::QueryResult::PutRecord(Err(err)) => {
                                debug!("KAD: Failed to put record: {err:?}");
                            }
                            kad::QueryResult::StartProviding(Ok(kad::AddProviderOk { key })) => {
                                debug!(
                                    "KAD: Successfully put provider record {:?}",
                                    std::str::from_utf8(key.as_ref()).unwrap()
                                );
                            }
                            kad::QueryResult::StartProviding(Err(err)) => {
                                debug!("KAD: Failed to put provider record: {err:?}");
                            }
                            _ => {}
                        }
                    }
                    // Handle gossipsub behavior events
                    SwarmEvent::Behaviour(RamdBehaviorEvent::Gossipsub(gossipsub::Event::Message {
                        propagation_source: peer_id,
                        message_id: id,
                        message,
                    })) => {
                        info!("GOSSIP: Received gossipsub message. peer {}, id {}, data {}", peer_id, id, std::str::from_utf8(&message.data).unwrap());

                        // first validate that received message is received from the right topic
                        if message.topic != self.topic.hash() {
                            self.disconnect_peer(&peer_id);
                        }
                    }
                    SwarmEvent::Behaviour(RamdBehaviorEvent::Gossipsub(gossipsub::Event::Subscribed {
                        peer_id,
                        topic,
                    })) => {
                        debug!("GOSSIP: New peer subscribed to topic. peer {}, topic {}", peer_id, topic);

                        // if connected peer is not a boot node and the topic is wrong - disconnect from it
                        if !self.is_boot_node(&peer_id) && topic != self.topic.hash() {
                            self.disconnect_peer(&peer_id);
                        }
                    }
                    SwarmEvent::Behaviour(RamdBehaviorEvent::Gossipsub(gossipsub::Event::Unsubscribed {
                        peer_id,
                        topic,
                    })) => {
                        warn!("GOSSIP: Peer unsubscribed from topic. peer {}, topic {}", peer_id, topic);
                    }
                    SwarmEvent::Behaviour(RamdBehaviorEvent::Gossipsub(gossipsub::Event::GossipsubNotSupported {
                        peer_id,
                    })) => {
                        warn!("GOSSIP: Peer with not supporting gossipsub has connected. peer {}.", peer_id);

                        if !self.is_boot_node(&peer_id) {
                            warn!("Disconnecting from not supporting gossipsub peer {}.", peer_id);
                            self.disconnect_peer(&peer_id);
                        }
                    }
                    _ => {}
                }
            }
        }
    }

    /// Checks does peer id is one of the boot nodes from the config
    fn is_boot_node(&self, peer_id: &PeerId) -> bool {
        self.boot_nodes.iter().any(|peer| peer == peer_id)
    }

    /// Checks is the current amount of connected peers exceed configured limit or not
    fn is_peer_limit_reached(&self) -> bool {
        let network_info = self.swarm.network_info();
        network_info.num_peers() >= self.max_peers_limit
    }

    /// Fully disconnect from peer and remove it from the routing table
    fn disconnect_peer(&mut self, peer_id: &PeerId) {
        // disconnect from peer
        if let Err(e) = self.swarm.disconnect_peer_id(*peer_id) {
            error!("Failed to disconnect from peer. Reason: {e:?}");
        }

        // remove it from kademlia table
        let _ = self.swarm.behaviour_mut().kademlia.remove_peer(peer_id);
    }

    /// If private key was already created then recover it from the storage,
    /// otherwise create a new pair and store it
    fn get_node_key(storage: Arc<S>) -> eyre::Result<identity::Keypair> {
        if let Some(pk) = storage.get_opt(RAMD_P2P_KEYPAIR_KEY.into())? {
            // pk is already stored, recover it
            Ok(identity::Keypair::from_protobuf_encoding(&pk)?)
        } else {
            // pk doesn't exists yet, create a new one and store it
            let pk = identity::Keypair::generate_ed25519();
            storage.set(RAMD_P2P_KEYPAIR_KEY.into(), pk.to_protobuf_encoding()?)?;

            Ok(pk)
        }
    }
}
