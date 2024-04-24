use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
pub struct JsonRpcServerConfig {
    pub port: u16,
}

impl Default for JsonRpcServerConfig {
    fn default() -> Self {
        Self { port: 1319 }
    }
}
