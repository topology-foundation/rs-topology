use std::path::PathBuf;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
pub struct NodeConfig {
    pub root_path: PathBuf,
    pub config_path: PathBuf,
}

impl Default for NodeConfig {
    fn default() -> Self {
        Self {
            root_path: [env!("HOME"), ".ramd"].iter().collect(),
            config_path: [env!("HOME"), ".ramd", "config", "ramd.toml"]
                .iter()
                .collect(),
        }
    }
}
