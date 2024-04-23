use std::sync::Arc;

use crate::config::NodeConfig;
use ramd_db::rocks::RocksStorage;
use ramd_db::storage::Storage;

pub struct Node {
    storage: Arc<dyn Storage<Vec<u8>, Vec<u8>>>,
}

impl Node {
    pub fn with_config(config: NodeConfig) -> eyre::Result<Self> {
        let storage = Arc::new(RocksStorage::new(&config.rocks)?);

        Ok(Node {
            storage: storage.clone(),
        })
    }
}
