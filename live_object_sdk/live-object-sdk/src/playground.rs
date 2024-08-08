use std::fs::File;
use std::io::Read;
use std::path::Path;
use wasmer::{imports, Instance, Module, Store, Value, WasmPtr};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let path = Path::new("../live-object/res/native/live_object.wasm");
    let mut file = File::open(path)?;
    let mut wasm_bytes = Vec::new();
    file.read_to_end(&mut wasm_bytes)?;

    // Create a Store.
    let mut store = Store::default();

    // Let's compile the Wasm module.
    let module = Module::new(&store, wasm_bytes)?;

    // Create an empty import object.
    let import_object = imports! {};

    // Let's instantiate the Wasm module.
    let instance = Instance::new(&mut store, &module, &import_object)?;

    // Allocate begins
    // Client
    let json_args = r#"{"x": 3,  "y": 4}"#;
    println!("{:?}", json_args);

    // VM
    let bytes_args = json_args.as_bytes();

    let allocate_func = instance.exports.get_function("allocate")?;
    let memory_slice_ptr =
        allocate_func.call(&mut store, &[u32::try_from(bytes_args.len())?.into()])?;

    assert_eq!(memory_slice_ptr.len(), 1);

    let memory_slice_ptr = memory_slice_ptr[0].clone();
    let memory_slice_ptr: u32 = memory_slice_ptr.try_into().unwrap();
    // memory_slice_ptr is a pointer to a struct MemorySlice in WASM memory.

    println!("{:?}", memory_slice_ptr);
    // Allocate ends

    // WriteMemory begins
    // Rust compiler auto-generates `(memory $memory (export "memory") 16)`.
    let memory = instance.exports.get_memory("memory")?;
    {
        let memory_view = memory.view(&store);

        let memory_slice = MemorySlice::new(&memory_view, memory_slice_ptr).unwrap();
        memory_slice.write(&memory_view, bytes_args).unwrap();
    }

    let params_ptr: Vec<Value> = vec![memory_slice_ptr.into()];
    // WriteMemory ends

    // Call function begins
    let sum = instance.exports.get_function("sum")?;
    let result_ptr = sum.call(&mut store, &params_ptr)?;
    // Call function ends

    // Parse result begins
    let memory_view = memory.view(&store);
    let memory_slice =
        MemorySlice::new(&memory_view, result_ptr[0].clone().try_into().unwrap()).unwrap();
    let result = memory_slice.read(&memory_view, 9999).unwrap();
    let result: String = String::from_utf8(result).unwrap();
    println!("Results: {:?}", result);
    // Parse result ends

    // Deallocate begins
    let deallocate_func = instance.exports.get_function("deallocate")?;
    deallocate_func
        .call(&mut store, &[result_ptr[0].clone()])
        .unwrap();
    // Deallocate ends

    Ok(())
}

use std::mem::size_of;

pub struct MemorySlice {
    /// A pointer to the start of this memory slice,
    /// measured in bytes from the beginning of the WASM (guest) memory.
    pub ptr: u32,
    /// The number of bytes in this memory slice.
    pub len: u32,
}

type MemorySlicePtrBytes = [u8; size_of::<MemorySlice>()];

impl MemorySlice {
    /// Read in a `MemorySlice` from the WASM (guest) memory and return it.
    pub fn new(memory: &wasmer::MemoryView, ptr: u32) -> core::result::Result<MemorySlice, ()> {
        let wasm_ptr = WasmPtr::<MemorySlicePtrBytes>::new(ptr);
        let memory_slice_ptr_bytes = wasm_ptr.deref(memory).read().map_err(|_err| ())?;
        let memory_slice = MemorySlice::from_memory_slice_ptr_bytes(memory_slice_ptr_bytes);

        MemorySlice::validate(&memory_slice)?;

        Ok(memory_slice)
    }

    /// Write the given data to the memory slice.
    pub fn write(self, memory: &wasmer::MemoryView, data: &[u8]) -> core::result::Result<(), ()> {
        if data.len() > self.len as usize {
            return Err(());
        }

        memory.write(self.ptr as u64, data).map_err(|_err| ())?;

        Ok(())
    }

    /// Read the memory slice.
    pub fn read(
        self,
        memory: &wasmer::MemoryView,
        max_len: usize,
    ) -> core::result::Result<Vec<u8>, ()> {
        if self.len as usize > max_len {
            return Err(());
        }

        let mut data = vec![0u8; self.len as usize];
        memory.read(self.ptr as u64, &mut data).map_err(|_err| ())?;

        Ok(data)
    }

    /// Convert a `MemorySlicePtrBytes` to a `MemorySlice`.
    fn from_memory_slice_ptr_bytes(bytes: MemorySlicePtrBytes) -> Self {
        MemorySlice {
            ptr: u32::from_le_bytes(bytes[0..4].try_into().unwrap()),
            len: u32::from_le_bytes(bytes[4..8].try_into().unwrap()),
        }
    }

    /// Validate the memory slice.
    fn validate(memory_slice: &MemorySlice) -> core::result::Result<(), ()> {
        if memory_slice.ptr == 0 {
            return Err(());
        }

        if memory_slice.len > (u32::MAX - memory_slice.ptr) {
            return Err(());
        }

        Ok(())
    }
}
