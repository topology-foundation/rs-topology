use libp2p::{multiaddr::Protocol, Multiaddr};
use serde::{Deserialize, Serialize};
use std::{str::FromStr, time::Duration};

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
#[serde(default)]
pub struct P2pConfig {
    pub port: u16,
    pub idle_connection_timeout_secs: u64,
    pub boot_nodes: Vec<String>,
    pub peers: Option<Vec<String>>,
    pub topic: String,
}

impl P2pConfig {
    pub fn idle_connection_timeout(&self) -> Duration {
        Duration::from_secs(self.idle_connection_timeout_secs)
    }

    pub fn peer_addresses(&self) -> eyre::Result<Option<Vec<Multiaddr>>> {
        if let Some(peers) = self.peers.clone() {
            let mut multi_addrs = vec![];

            for peer in peers.iter() {
                let mut addr = Multiaddr::from_str(peer)?;

                let last = addr.pop();
                match last {
                    // for a multiaddr that ends with a peer id, this strips this suffix. Rust-libp2p
                    // only supports dialing to an address without providing the peer id.
                    Some(Protocol::P2p(_peer_id)) => {}
                    // if its another protocol appened suffix back
                    Some(other) => addr.push(other),
                    _ => {}
                }

                multi_addrs.push(addr);
            }

            Ok(Some(multi_addrs))
        } else {
            Ok(None)
        }
    }
}

impl Default for P2pConfig {
    fn default() -> Self {
        Self {
            port: 1211,
            idle_connection_timeout_secs: 60,
            boot_nodes: vec![
                "QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN".to_owned(), // TODO: set default values to our node once done
                "QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa".to_owned(),
                "QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb".to_owned(),
                "QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt".to_owned(),
            ],
            peers: None,
            topic: "ramd-topic".to_owned(),
        }
    }
}
