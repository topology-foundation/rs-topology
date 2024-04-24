use std::sync::Arc;

use ramd_db::storage::Storage;
use tracing::{error, info};

pub enum Action {
    CreateLiveObject(CreateLiveObjectAction),
    ExecuteLiveObject(ExecuteLiveObjectAction),
}

impl Action {
    pub(crate) fn perform<S>(&self, cache: Arc<S>) -> eyre::Result<()>
    where
        S: Storage<Vec<u8>, Vec<u8>>,
    {
        match self {
            Action::CreateLiveObject(action) => action.perform(cache),
            Action::ExecuteLiveObject(action) => action.perform(cache),
        }
    }
}

pub struct CreateLiveObjectAction {
    pub wasm_bytes: Vec<u8>,
}

impl CreateLiveObjectAction {
    fn perform<S>(&self, cache: Arc<S>) -> eyre::Result<()>
    where
        S: Storage<Vec<u8>, Vec<u8>>,
    {
        // TODO: use some cryptographic hash as a key.
        if let Err(e) = cache.set(vec![0], self.wasm_bytes.clone()) {
            error!(target: "ramd::processor", "Failed to set wasm bytes to cache with error `{}`", e.to_string());
            return Err(e);
        }

        info!(target: "ramd::processor", "Successfully performed create action");
        Ok(())
    }
}

pub struct ExecuteLiveObjectAction {
    pub live_object_id: [u8; 32],
    pub method: String,
    pub args: Vec<u8>,
}

impl ExecuteLiveObjectAction {
    fn perform<S>(&self, _cache: Arc<S>) -> eyre::Result<()>
    where
        S: Storage<Vec<u8>, Vec<u8>>,
    {
        info!(target: "ramd::processor", "Successfully performed execute action");
        Ok(())
    }
}
