use std::thread::park;

use ramd_config::config::RamdConfig;
use ramd_jsonrpc_server::launch;
use ramd_tracing::init as init_tracing;
use tracing::info;

#[tokio::main]
async fn main() -> eyre::Result<()> {
    // Init tracing logger
    init_tracing();

    info!("Topology is a community-driven technology that brings random access memory to the world computer to power lock-free asynchronous decentralized applications.");

    // Init or read ramd config
    let ramd_config = RamdConfig::init_or_read()?;

    // Launch jsonrpc server
    let handle = launch(&ramd_config.json_rpc).await?;

    // // TODO: for now we don't care about server, simply start it and forget
    // // Revisit once proper server handle handling will be required
    tokio::spawn(handle.stopped());

    // TODO: implement proper process handler, for now simply park the main thread until ctrl+c
    park();

    Ok(())
}
