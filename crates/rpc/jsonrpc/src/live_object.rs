use std::sync::Arc;

use async_trait::async_trait;
use base64::prelude::{Engine, BASE64_STANDARD};
use jsonrpsee::core::RpcResult;
use jsonrpsee::types::{error::ErrorObject, ErrorCode};
use ramd_jsonrpc_api::server::LiveObjectApiServer;
use ramd_jsonrpc_types::live_object::CreateLiveObject;
use ramd_node::LiveObjectHandler;
use tracing::info;

pub struct LiveObjectApi {
    node: Arc<dyn LiveObjectHandler>,
}

impl LiveObjectApi {
    pub fn new(node: Arc<dyn LiveObjectHandler>) -> Self {
        Self { node }
    }
}

#[async_trait]
impl LiveObjectApiServer for LiveObjectApi {
    async fn create_live_object(&self, request: CreateLiveObject) -> RpcResult<()> {
        info!(target: "ramd::jsonrpc", "Request to create a live object with wasm bytes {}", request.wasm_bytes);

        let wasm_bytes = match BASE64_STANDARD.decode(request.wasm_bytes) {
            Ok(bytes) => bytes,
            Err(e) => {
                info!(target: "ramd::jsonrpc", "Failed to decode wasm bytes with error `{}`", e.to_string());

                return Err(ErrorObject::from(ErrorCode::InvalidParams));
            }
        };

        self.node.create_live_object(wasm_bytes);

        Ok(())
    }
}
