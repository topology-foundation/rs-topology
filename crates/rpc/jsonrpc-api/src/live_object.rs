use jsonrpsee::{core::RpcResult, proc_macros::rpc};
use ramd_jsonrpc_types::live_object::CreateLiveObject;

#[rpc(server, client, namespace = "live_object")]
pub trait LiveObjectApi {
    #[method(name = "create")]
    async fn create_live_object(&self, request: CreateLiveObject) -> RpcResult<()>;
}
