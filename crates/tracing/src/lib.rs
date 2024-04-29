use std::io;
use std::path::PathBuf;

use rolling_file::{RollingConditionBasic, RollingFileAppender};
use serde::{Deserialize, Serialize};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt, EnvFilter, Layer};

#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
#[serde(default)]
pub struct TracingConfig {
    pub path: PathBuf,
    pub max_size_bytes: u64,
    pub max_files: usize,
}

impl TracingConfig {
    pub fn new(root_path: PathBuf) -> Self {
        let log_dir = root_path.join(Self::log_file_dir());
        Self {
            path: log_dir,
            max_size_bytes: 200,
            max_files: 5,
        }
    }

    fn log_file_dir() -> PathBuf {
        "logs".into()
    }

    pub fn log_file_name(&self) -> PathBuf {
        self.path.join("ramd.log")
    }
}

/// Init basic fmt logger that writes to console
///
/// Note:
/// RUST_LOG env must be set in order for tracing to start working
pub fn init(config: &TracingConfig) {
    let (file_appender, _guard) = tracing_appender::non_blocking(
        RollingFileAppender::new(
            config.log_file_name(),
            RollingConditionBasic::new().max_size(config.max_size_bytes),
            config.max_files,
        )
        .expect("Failed to initialize file appender"),
    );

    let layers = vec![
        tracing_subscriber::fmt::layer()
            .with_writer(io::stdout)
            .with_filter(EnvFilter::builder().from_env_lossy())
            .boxed(),
        tracing_subscriber::fmt::layer()
            .with_target(true)
            .with_writer(file_appender)
            .with_filter(EnvFilter::builder().from_env_lossy())
            .boxed(),
        tracing_journald::layer()
            .expect("Faiiled to init journald layer")
            .with_filter(EnvFilter::builder().from_env_lossy())
            .boxed(),
    ];

    let _ = tracing_subscriber::registry().with(layers).try_init();
}
