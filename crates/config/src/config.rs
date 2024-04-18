use std::path::PathBuf;

use eyre::OptionExt;
use ramd_db::rocks::RocksConfig;
use ramd_jsonrpc_server::JsonRpcServerConfig;
use ramd_tracing::TracingConfig;
use serde::{Deserialize, Serialize};

/// Directory path for storing all ramd related data
const RAMD_DIR: &str = ".ramd";

/// Directory path for storing ramd config information
const CONFIG_DIR: &str = "config";

/// This struct gathers all config values used across ramd node
#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
pub struct RamdConfig {
    /// Configuration for tracing/logging
    pub tracing: TracingConfig,
    /// Configuration for rocksdb storage
    pub rocks: RocksConfig,
    /// Configuration for jsonrpc server
    pub json_rpc: JsonRpcServerConfig,
}

impl RamdConfig {
    /// Reads config from default path, returns error if config doesn't exists
    pub fn read() -> eyre::Result<Self> {
        let home_path = std::env::var("HOME")?;
        let ramd_config_path: PathBuf = [home_path.as_str(), RAMD_DIR, CONFIG_DIR, "ramd.toml"]
            .iter()
            .collect();

        let config = std::fs::read_to_string(ramd_config_path)
            .map_err(|_| eyre::eyre!("Path doesn't exist"))?;

        let config: RamdConfig = toml::from_str(&config)?;
        Ok(config)
    }

    /// Creates default config if not exists otherwise reads it
    pub fn init_or_read() -> eyre::Result<Self> {
        let config_maybe = RamdConfig::read();
        if let Ok(config) = config_maybe {
            return Ok(config);
        };

        let home_path = std::env::var("HOME")?;

        let root_dir: PathBuf = [home_path.as_str(), RAMD_DIR].iter().collect();
        std::fs::create_dir_all(&root_dir)?;

        let config_dir = root_dir.join(CONFIG_DIR);
        std::fs::create_dir(&config_dir)?;

        let db_data_dir = root_dir.join(RocksConfig::db_name());
        std::fs::create_dir(&db_data_dir)?;

        let config = RamdConfig {
            rocks: RocksConfig {
                path: db_data_dir
                    .to_str()
                    .ok_or_eyre("Failed to get rocksdb data directory")?
                    .to_owned(),
            },
            ..Default::default()
        };

        let config_path = config_dir.join("ramd.toml");

        let toml_config = toml::to_string(&config)?;
        std::fs::write(config_path, toml_config)?;

        Ok(config)
    }
}
