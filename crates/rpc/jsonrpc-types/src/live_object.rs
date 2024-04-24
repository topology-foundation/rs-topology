use base64::prelude::{Engine, BASE64_STANDARD};
use jsonrpsee::core::RpcResult;
use jsonrpsee::types::{error::ErrorObject, ErrorCode};
use serde::{Deserialize, Serialize};
use tracing::error;

#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct CreateLiveObject {
    pub wasm_bytes: String, // Base64 encoded wasm bytes.
}

impl CreateLiveObject {
    pub fn decode_wasm_bytes(&self) -> RpcResult<Vec<u8>> {
        match BASE64_STANDARD.decode(self.wasm_bytes.clone()) {
            Ok(bytes) => Ok(bytes),
            Err(e) => {
                error!(target: "ramd::jsonrpc-types", "Failed to decode wasm bytes with error `{}`", e.to_string());

                return Err(ErrorObject::from(ErrorCode::InvalidParams));
            }
        }
    }
}
