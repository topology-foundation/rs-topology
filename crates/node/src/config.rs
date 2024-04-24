use ramd_db::config::RocksConfig;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
pub struct NodeConfig {
    pub rocks: RocksConfig,
    // TODO: add VMConfig.
}
