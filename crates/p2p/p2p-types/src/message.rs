use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize)]
pub enum P2pMessage {
    Noop { data: String },
}
