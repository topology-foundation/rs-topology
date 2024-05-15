use std::io;

use ramd_config::configs::tracing::TracingConfig;
use rolling_file::{RollingConditionBasic, RollingFileAppender};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt, EnvFilter, Layer};

/// Init basic fmt logger that writes to console
///
/// Note:
/// RUST_LOG env must be set in order for tracing to start working
pub fn init(config: &TracingConfig) {
    println!("{}", config.path.as_os_str().to_str().unwrap());
    let (file_appender, _guard) = tracing_appender::non_blocking(
        RollingFileAppender::new(
            config.path.as_os_str().to_str().unwrap(),
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
        // tracing_journald::layer()
        //     .expect("Failed to init journald layer")
        //     .with_filter(EnvFilter::builder().from_env_lossy())
        //     .boxed(),
    ];

    let _ = tracing_subscriber::registry().with(layers).try_init();
}
