use std::path::PathBuf;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
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
