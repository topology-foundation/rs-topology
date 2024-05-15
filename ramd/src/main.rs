use clap::Parser;
use cli::Subcommand;
use commands::NodeCmd;
use dotenv::dotenv;
use eyre::{eyre, Result};
use ramd_config::{
    configs::{
        network::P2pConfig, node::NodeConfig, rpc::JsonRpcServerConfig, storage::RocksConfig,
        tracing::TracingConfig,
    },
    RamdConfig,
};
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
async fn main() -> Result<()> {
    let cli: Cli = Cli::parse();

    match cli.subcommand {
        Some(Subcommand::Bootnode(_)) => Err(eyre!("Bootnode not implemented!")),
        Some(Subcommand::Node(flags)) => {
            let config: RamdConfig = parse_flags(flags)?;

            // parse .env faile
            dotenv().ok();

            if let Err(e) = start(config).await {
                return Err(eyre!("Failed to start ramd node. Reason: {}", e));
            }

            // TODO: implement proper process handler, for now simply park the main thread until ctrl+c
            park();

            Ok(())
        }
        Some(Subcommand::Relayer(_)) => Err(eyre!("Relayer node not implemented!")),
        // Handled by #[command(arg_required_else_help = true)]
        None => Ok(()),
    }
}

/// This is a temp solution to properly log received error during start-up process
async fn start(config: RamdConfig) -> eyre::Result<()> {
    // Init or read ramd config
    // let ramd_config = RamdConfig::init_or_read()?;

    // Init tracing logger
    init_tracing(&config.tracing);

    tracing::info!("Topology is a community-driven technology that brings random access memory to the world computer to power lock-free asynchronous decentralized applications.");

    // Construct RocksDB
    let rocks = Arc::new(RocksStorage::new(&config.rocks)?);

    // Construct a RAM node
    let node = Arc::new(Node::new(&config.node, rocks.clone())?);

    // Launch p2p server
    let (mut p2p, _p2p_msg_sender) = P2pServer::new(&config.p2p, rocks.clone())?;
    tokio::spawn(async move { p2p.launch().await });

    // Launch jsonrpc server
    // TODO: for now we don't care about server, simply start it and forget
    // Revisit once proper server handle handling will be required
    let handle = launch(&config.json_rpc, node.clone()).await?;
    tokio::spawn(handle.stopped());

    Ok(())
}

fn parse_flags(mut flags: NodeCmd) -> Result<RamdConfig> {
    // set pathbufs for non-default `ram_dir_name` and default paths
    if flags
        .clone()
        .node
        .ramd_dir_name
        .into_os_string()
        .into_string()
        .unwrap()
        .starts_with("$HOME")
    {
        flags.node.ramd_dir_name = [
            env!("HOME").to_string(),
            flags
                .node
                .ramd_dir_name
                .into_os_string()
                .into_string()
                .unwrap()
                .replace("$HOME/", ""),
        ]
        .iter()
        .collect();
    }

    if !flags
        .clone()
        .db
        .db_rocks_path
        .into_os_string()
        .into_string()
        .unwrap()
        .starts_with('/')
    {
        flags.db.db_rocks_path = [flags.node.ramd_dir_name.clone(), flags.db.db_rocks_path]
            .iter()
            .collect()
    }

    if !flags
        .clone()
        .network
        .network_config_path
        .into_os_string()
        .into_string()
        .unwrap()
        .starts_with('/')
    {
        flags.network.network_config_path = [
            flags.node.ramd_dir_name.clone(),
            flags.network.network_config_path,
        ]
        .iter()
        .collect()
    }

    if !flags
        .clone()
        .node
        .ramd_config_file
        .into_os_string()
        .into_string()
        .unwrap()
        .starts_with('/')
    {
        flags.node.ramd_config_file = [
            flags.node.ramd_dir_name.clone(),
            flags.node.ramd_config_file,
        ]
        .iter()
        .collect()
    }

    if !flags
        .clone()
        .tracing
        .tracing_path
        .into_os_string()
        .into_string()
        .unwrap()
        .starts_with('/')
    {
        flags.tracing.tracing_path = [flags.node.ramd_dir_name, flags.tracing.tracing_path]
            .iter()
            .collect()
    }

    Ok(RamdConfig {
        node: NodeConfig {},
        rocks: RocksConfig {
            path: flags.db.db_rocks_path,
        },
        json_rpc: JsonRpcServerConfig {
            port: flags.rpc.json_rpc_port,
        },
        p2p: P2pConfig {
            boot_nodes: flags.network.network_boot_nodes,
            config_path: flags.network.network_config_path,
            idle_connection_timeout_secs: flags.network.network_idle_connection_timeout,
            max_peers_limit: flags.network.network_max_peers_limit,
            network_key: flags.network.network_key,
            port: flags.network.network_port,
        },
        tracing: TracingConfig {
            path: flags.tracing.tracing_path,
            max_files: flags.tracing.tracing_max_files,
            max_size_bytes: flags.tracing.tracing_max_size_bytes,
        },
    })
}
