use crate::commands::{BootnodeCmd, NodeCmd, RelayerCmd};
use clap::Parser;

#[derive(Debug, Parser)]
#[command(version, about, arg_required_else_help = true)]
pub struct Cli {
    #[command(subcommand)]
    pub subcommand: Option<Subcommand>,
}

#[derive(Debug, clap::Subcommand)]
pub enum Subcommand {
    /// Runs ramd as bootnode mode, where the only functionalities are peer discovery
    Bootnode(BootnodeCmd),

    /// Runs ramd node with full functionalities
    Node(NodeCmd),

    // NOTE: This one might make sense to have a separate implementation/repo
    /// Runs ramd relayer node. The only functionality is relaying messages
    Relayer(RelayerCmd),
}
