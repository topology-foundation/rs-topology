use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
pub struct TracingConfig {}

/// Init basic fmt logger that writes to console
///
/// Note:
/// RUST_LOG env must be set in order for tracing to start working
pub fn init() {
    tracing_subscriber::fmt::init();
}
