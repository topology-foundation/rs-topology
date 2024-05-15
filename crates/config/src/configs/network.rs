// use libp2p::{multiaddr::Protocol, Multiaddr};
use serde::{Deserialize, Serialize};
use std::{path::PathBuf, time::Duration};

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
#[serde(default)]
pub struct P2pConfig {
    pub boot_nodes: Option<Vec<String>>,
    pub config_path: PathBuf,
    pub idle_connection_timeout_secs: u64,
    pub max_peers_limit: usize,
    pub network_key: Option<PathBuf>,
    pub port: u16,
}

impl P2pConfig {
    pub fn idle_connection_timeout(&self) -> Duration {
        Duration::from_secs(self.idle_connection_timeout_secs)
    }
}

impl Default for P2pConfig {
    fn default() -> Self {
        Self {
            boot_nodes: None,
            config_path: PathBuf::new(),
            idle_connection_timeout_secs: 60,
            max_peers_limit: 10,
            network_key: None,
            port: 1211,
        }
    }
}
