use std::sync::Arc;

use crate::message::Message;
use ramd_db::storage::Storage;
use tracing::info;

pub struct Processor {
    storage: Arc<dyn Storage<Vec<u8>, Vec<u8>>>,
}

impl Processor {
    pub fn new(storage: Arc<dyn Storage<Vec<u8>, Vec<u8>>>) -> Self {
        Processor { storage }
    }

    pub fn process_messages(&self, messages: &[Message]) {
        // TODO: use cache that wraps around storage.
        let cache = self.storage.clone();

        // TODO: add to messsage pool and then process messages.

        for message in messages {
            if let Err(_) = message.process(cache.clone()) {
                // TODO: log message ID.
                info!(target: "ramd::processor", "Failed to process a message");
                return;
            }
        }

        // TODO: cache.commit();
    }
}
