use crate::storage::Storage;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
pub struct RocksConfig {
    pub path: String,
}

impl RocksConfig {
    // TODO: there might be a better way to build/read db path, implement one needed
    pub fn db_name() -> String {
        "ramd_db".to_owned()
    }
}

pub struct RocksStorage {
    db: rocksdb::DB,
}

impl RocksStorage {
    pub fn new(config: &RocksConfig) -> eyre::Result<Self> {
        let db = rocksdb::DB::open_default(&config.path)?;

        Ok(Self { db })
    }
}

impl<K: AsRef<[u8]>, V: AsRef<[u8]>> Storage<K, V> for RocksStorage {
    fn has(&self, key: K) -> eyre::Result<bool> {
        let v = self.db.get(key)?;
        Ok(v.is_some())
    }

    fn get(&self, key: K) -> eyre::Result<Vec<u8>> {
        let v = self.db.get(key)?;
        if let Some(v) = v {
            Ok(v)
        } else {
            Err(eyre::eyre!("Key not found"))
        }
    }

    fn get_opt(&self, key: K) -> eyre::Result<Option<Vec<u8>>> {
        let v = self.db.get(key)?;
        Ok(v)
    }

    fn set(&self, key: K, value: V) -> eyre::Result<()> {
        self.db.put(key, value)?;
        Ok(())
    }

    fn delete(&self, key: K) -> eyre::Result<()> {
        self.db.delete(key)?;
        Ok(())
    }
}