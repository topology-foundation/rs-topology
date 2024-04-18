use async_trait::async_trait;
use jsonrpsee::core::RpcResult;
use ramd_jsonrpc_api::server::LiveObjectApiServer;
use ramd_jsonrpc_types::live_object::CreateLiveObject;
use tracing::info;

#[derive(Default)]
pub struct LiveObjectApi {}

impl LiveObjectApi {
    pub fn new() -> Self {
        Self {}
    }
}

#[async_trait]
impl LiveObjectApiServer for LiveObjectApi {
    async fn create_live_object(&self, request: CreateLiveObject) -> RpcResult<()> {
        info!("Received create live object data - {}", request.data);
        Ok(())
    }
}
