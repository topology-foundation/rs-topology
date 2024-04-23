use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct CreateLiveObject {
    pub wasm_bytes: String, // Base64 encoded wasm bytes.
}
