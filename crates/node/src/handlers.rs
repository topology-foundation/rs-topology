pub trait LiveObjectHandler: Send + Sync {
    fn create_live_object(&self, wasm_bytes: Vec<u8>);
}
