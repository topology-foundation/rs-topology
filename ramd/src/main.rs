use clap::Parser;
use dotenv::dotenv;
use ramd_config::RamdConfig;
use ramd_db::rocks::RocksStorage;
use ramd_jsonrpc_server::launch;
use ramd_node::Node;
use ramd_p2p_server::Server as P2pServer;
use ramd_tracing::init as init_tracing;
use std::{sync::Arc, thread::park};

mod cli;
mod commands;
use crate::cli::Cli;

/// Note: I think ideally inside of a main function we should create a ramd instance, with builder pattern to configure everything needed and then call some
/// sort of a blocking run function, so that all the modules we have like p2p, jsonrpc etc. are configured outside of the main function.
///
/// Example:
/// #[tokio::main]
/// async fn main() -> eyre::Result<()> {
///     let ramd = RamdBuilder::new().with_a().with_b().build().await;
///     ramd.run().await;
/// }
#[tokio::main(flavor = "multi_thread")]
async fn main() -> eyre::Result<()> {
    let cli = Cli::parse();

    // parse .env faile
    dotenv().ok();

    if let Err(e) = start().await {
        return Err(eyre::eyre!("Failed to start ramd node. Reason: {}", e));
    }

    // TODO: implement proper process handler, for now simply park the main thread until ctrl+c
    park();

    Ok(())
}

/// This is a temp solution to properly log received error during start-up process
async fn start() -> eyre::Result<()> {
    // Init or read ramd config
    let ramd_config = RamdConfig::init_or_read()?;

    // Init tracing logger
    init_tracing(&ramd_config.tracing);

    tracing::info!("Topology is a community-driven technology that brings random access memory to the world computer to power lock-free asynchronous decentralized applications.");

    // Construct RocksDB
    let rocks = Arc::new(RocksStorage::new(&ramd_config.rocks)?);

    // Construct a RAM node
    let node = Arc::new(Node::new(&ramd_config.node, rocks.clone())?);

    // Launch p2p server
    let (mut p2p, _p2p_msg_sender) = P2pServer::new(&ramd_config.p2p, rocks.clone())?;
    tokio::spawn(async move { p2p.launch().await });

    // Launch jsonrpc server
    // TODO: for now we don't care about server, simply start it and forget
    // Revisit once proper server handle handling will be required
    let handle = launch(&ramd_config.json_rpc, node.clone()).await?;
    tokio::spawn(handle.stopped());

    Ok(())
}
