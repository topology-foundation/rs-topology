pub trait Storage<K, V>: Send + Sync
where
    K: AsRef<[u8]>,
    V: AsRef<[u8]>,
{
    fn has(&self, key: K) -> eyre::Result<bool>;
    fn get(&self, key: K) -> eyre::Result<Vec<u8>>;
    fn get_opt(&self, key: K) -> eyre::Result<Option<Vec<u8>>>;
    fn set(&self, key: K, value: V) -> eyre::Result<()>;
    fn delete(&self, key: K) -> eyre::Result<()>;
}
