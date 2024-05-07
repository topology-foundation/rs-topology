use futures::prelude::*;
use libp2p::{
    identify, identity,
    kad::{self, Mode},
    noise,
    swarm::{NetworkBehaviour, SwarmEvent},
    tcp, yamux,
};
use ramd_db::{
    config::RocksConfig, keys::RAMD_P2P_KEYPAIR_KEY, rocks::RocksStorage, storage::Storage,
};
use ramd_tracing::{config::TracingConfig, init as init_tracing};
use serde::{Deserialize, Serialize};
use std::time::Duration;
use tracing::{error, info};

const RAM_PROTOCOL_VERSION: &str = "ram/0.1.0";

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
pub struct BootstrapNodeInfo {
    pub peer_id: String,
    pub address: String,
}

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
pub struct Config {
    pub port: u16,
    pub bootstrap_interval_secs: u64,
    pub idle_connection_timeout_secs: u64,
    pub boot_nodes: Option<Vec<BootstrapNodeInfo>>,
}

impl Config {
    pub fn load() -> eyre::Result<Self> {
        let config = std::fs::read_to_string("./boot.toml").map_err(|_| {
            eyre::eyre!("File doesn't exist at project root. Make sure boot.toml is present")
        })?;

        let config: Self = toml::from_str(&config)?;
        Ok(config)
    }
}

#[tokio::main(flavor = "multi_thread")]
async fn main() {
    // Init tracing logger
    init_tracing(&TracingConfig::default());

    let boot_config = match Config::load() {
        Ok(cfg) => cfg,
        Err(e) => {
            panic!("Failed to configure p2p server for boot node. Error: {e:?}");
        }
    };

    // configure p2p server for boot node
    let mut swarm = match setup_swarm(&boot_config).await {
        Ok(swarm) => swarm,
        Err(e) => {
            panic!("Failed to configure p2p server for boot node. Error: {e:?}");
        }
    };

    let peer_id = swarm.local_peer_id();
    info!("Launching ram network bootstrap node with peer id: {peer_id}");

    let mut interval = tokio::time::interval(std::time::Duration::from_secs(
        boot_config.bootstrap_interval_secs,
    ));

    // listen for incomming events
    loop {
        tokio::select! {
            _ = interval.tick() => {
                if let Err(e) = swarm.behaviour_mut().kademlia.bootstrap() {
                    error!(target: "bootnode", "Bootstrap step has failed, waiting for next iteration. Reason: {}", e.to_string());
                }
            }
            event = swarm.select_next_some() => match event {
                SwarmEvent::NewListenAddr { address, .. } => {
                    info!(target: "bootnode", "One of our listeners has reported a new local listening address. Listening on {address:?}");
                }
                SwarmEvent::Behaviour(BootNodeBehaviorEvent::Identify(identify::Event::Received { peer_id, info, })) => {
                    // verify that peer has right protocol version
                    // and supports kademlia, if so add received addresses to the routing table
                    if info.protocol_version.as_str() == RAM_PROTOCOL_VERSION && info.protocols.iter().any(|p| *p == kad::PROTOCOL_NAME) {
                        info!(target: "bootnode", "Peer {peer_id} has been discovered, adding to the DHT");
                        for addr in info.listen_addrs.into_iter() {
                            swarm
                                .behaviour_mut()
                                .kademlia
                                .add_address(&peer_id, addr);
                        }
                    }
                }
                _ => {}
            }
        }
    }
}

#[derive(NetworkBehaviour)]
struct BootNodeBehavior {
    identify: identify::Behaviour,
    kademlia: kad::Behaviour<kad::store::MemoryStore>,
}

async fn setup_swarm(boot_config: &Config) -> eyre::Result<libp2p::Swarm<BootNodeBehavior>> {
    // setup storage, we need it only for storing node's private key
    let rocks_cfg = RocksConfig::new("./".into());
    let db: Box<dyn Storage<Vec<u8>, Vec<u8>>> = Box::new(RocksStorage::new(&rocks_cfg)?);

    // get node's private key
    let node_key = if let Some(pk) = db.get_opt(RAMD_P2P_KEYPAIR_KEY.into())? {
        identity::Keypair::from_protobuf_encoding(&pk)?
    } else {
        let pk = identity::Keypair::generate_ed25519();
        db.set(RAMD_P2P_KEYPAIR_KEY.into(), pk.to_protobuf_encoding()?)?;

        pk
    };

    let mut swarm = libp2p::SwarmBuilder::with_existing_identity(node_key)
        .with_tokio()
        .with_tcp(
            tcp::Config::default(),
            noise::Config::new,
            yamux::Config::default,
        )?
        .with_dns()?
        .with_behaviour(|key| {
            // Configure kademlia behavior
            let peer_id = key.public().to_peer_id();
            let kademlia =
                kad::Behaviour::new(peer_id.clone(), kad::store::MemoryStore::new(peer_id));

            // Configure identify protocol so that this node cane be discovered
            let identify = identify::Behaviour::new(identify::Config::new(
                RAM_PROTOCOL_VERSION.to_string(),
                key.public(),
            ));

            Ok(BootNodeBehavior { identify, kademlia })
        })?
        .with_swarm_config(|c| {
            c.with_idle_connection_timeout(Duration::from_secs(
                boot_config.idle_connection_timeout_secs,
            ))
        })
        .build();

    // Adding boot node addresses for initial peer discovery
    if let Some(boot_nodes) = boot_config.boot_nodes.clone() {
        for boot_node in boot_nodes.iter() {
            let peer_id = boot_node.peer_id.parse()?;
            let address = boot_node.address.parse()?;

            swarm
                .behaviour_mut()
                .kademlia
                .add_address(&peer_id, address);
        }
    }

    swarm.behaviour_mut().kademlia.set_mode(Some(Mode::Server));
    swarm.listen_on(format!("/ip4/0.0.0.0/tcp/{}", boot_config.port).parse()?)?;

    Ok(swarm)
}
