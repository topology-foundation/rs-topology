use std::net::SocketAddr;

pub use jsonrpsee::server::ServerBuilder;
use jsonrpsee::{server::ServerHandle, RpcModule};
use ramd_jsonrpc::live_object::LiveObjectApi;
use ramd_jsonrpc_api::server::LiveObjectApiServer;
use serde::{Deserialize, Serialize};
use tracing::info;

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
#[serde(default)]
pub struct JsonRpcServerConfig {
    pub port: u16,
}

impl Default for JsonRpcServerConfig {
    fn default() -> Self {
        Self { port: 1111 }
    }
}

/// Launch configured jsonrpc server
pub async fn launch(config: &JsonRpcServerConfig) -> eyre::Result<ServerHandle> {
    let mut module = RpcModule::new(());

    let live_object_api = LiveObjectApi::new();
    module
        .merge(live_object_api.into_rpc())
        .map_err(|_| eyre::eyre!("Live object API has conflicting methods"))?;

    let socket_addr = format!("0.0.0.0:{}", config.port).parse::<SocketAddr>()?;
    let server = ServerBuilder::new()
        .build(socket_addr)
        .await
        .map_err(|_| eyre::eyre!("Failed to build jsonrpc server"))?;

    let local_addr = server
        .local_addr()
        .map_err(|_| eyre::eyre!("Failed to get server local address"))?;

    info!(target: "ramd::jsonrpc-server", "Launching jsonrpc server. Address: {}", local_addr.to_string());
    let handle = server.start(module);

    Ok(handle)
}
