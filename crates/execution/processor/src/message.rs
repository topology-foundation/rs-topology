use std::sync::Arc;

use crate::Action;
use ramd_db::storage::Storage;

pub struct Message {
    // TODO: add ID and dependencies.
    // pub id: some_cryptographic_hash,
    // pub predecessors: Vec<some_cryptographic_hash>,
    pub action: Action,
}

impl Message {
    pub(crate) fn process<S>(&self, cache: Arc<S>) -> eyre::Result<()>
    where
        S: Storage<Vec<u8>, Vec<u8>>,
    {
        self.action.perform(cache)
    }
}
