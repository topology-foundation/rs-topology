use std::path::PathBuf;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
pub struct RocksConfig {
    pub path: PathBuf,
}

impl RocksConfig {
    pub fn new(root_path: PathBuf) -> Self {
        let db_path = root_path.join(Self::db_name());
        Self { path: db_path }
    }

    fn db_name() -> PathBuf {
        "ramd_db".into()
    }
}
