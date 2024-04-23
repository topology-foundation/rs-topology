use std::sync::Arc;

use crate::config::NodeConfig;
use crate::handlers::LiveObjectHandler;
use ramd_db::rocks::RocksStorage;
use ramd_processor::{CreateAction, Message, Processor};
use tracing::info;

pub struct Node {
    processor: Processor,
}

impl Node {
    pub fn with_config(config: NodeConfig) -> eyre::Result<Self> {
        let storage = Arc::new(RocksStorage::new(&config.rocks)?);

        Ok(Node {
            processor: Processor::new(storage.clone()),
        })
    }
}

impl LiveObjectHandler for Node {
    fn create_live_object(&self, wasm_bytes: Vec<u8>) {
        let create_action = CreateAction::new(wasm_bytes);
        let messages = vec![Message::from_action(create_action.into())];

        // TODO: log message ID.
        info!(target: "ramd::node", "New message with create action");

        self.processor.process_messages(&messages);
    }
}
