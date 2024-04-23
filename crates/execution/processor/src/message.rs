use std::sync::Arc;

use crate::action::Action;
use ramd_db::storage::Storage;

pub struct Message {
    // TODO: add ID and dependencies.
    // pub id: some_cryptographic_hash,
    // pub predecessors: Vec<some_cryptographic_hash>,
    pub action: Action,
}

impl Message {
    pub fn from_action(action: Action) -> Self {
        Message { action }
    }

    pub(crate) fn process(&self, cache: Arc<dyn Storage<Vec<u8>, Vec<u8>>>) -> eyre::Result<()> {
        self.action.perform(cache)?;
        Ok(())
    }
}
