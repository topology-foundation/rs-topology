use std::sync::Arc;

use ramd_db::storage::Storage;
use tracing::info;

pub enum Action {
    Create(CreateAction),
    Execute(ExecuteAction),
}

impl Action {
    pub(crate) fn perform(&self, cache: Arc<dyn Storage<Vec<u8>, Vec<u8>>>) -> eyre::Result<()> {
        match self {
            Action::Create(action) => self.perform_create(action, cache),
            Action::Execute(action) => self.perform_execute(action, cache),
        }
    }

    fn perform_create(
        &self,
        action: &CreateAction,
        cache: Arc<dyn Storage<Vec<u8>, Vec<u8>>>,
    ) -> eyre::Result<()> {
        // TODO: use some cryptographic hash as a key.
        if let Err(e) = cache.set(vec![0], action.wasm_bytes.clone()) {
            info!(target: "ramd::processor", "Failed to set wasm bytes to cache with error `{}`", e.to_string());
            return Err(e);
        }

        info!(target: "ramd::processor", "Successfully performed create action");

        Ok(())
    }

    fn perform_execute(
        &self,
        action: &ExecuteAction,
        cache: Arc<dyn Storage<Vec<u8>, Vec<u8>>>,
    ) -> eyre::Result<()> {
        info!(target: "ramd::processor", "Successfully performed execute action");

        Ok(())
    }
}

pub struct CreateAction {
    wasm_bytes: Vec<u8>,
}

impl CreateAction {
    pub fn new(wasm_bytes: Vec<u8>) -> Self {
        CreateAction { wasm_bytes }
    }
}

impl From<CreateAction> for Action {
    fn from(action: CreateAction) -> Self {
        Action::Create(action)
    }
}

pub struct ExecuteAction {
    live_object_id: [u8; 32],
    method: String,
    args: Vec<u8>,
}

impl ExecuteAction {
    pub fn new(live_object_id: [u8; 32], method: String, args: Vec<u8>) -> Self {
        ExecuteAction {
            live_object_id,
            method,
            args,
        }
    }
}

impl From<ExecuteAction> for Action {
    fn from(action: ExecuteAction) -> Self {
        Action::Execute(action)
    }
}
