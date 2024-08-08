.PHONY: build
build:
	@mkdir -p ./output

	rustup target add wasm32-unknown-unknown
	cargo build --all --target wasm32-unknown-unknown --release
	cp target/wasm32-unknown-unknown/release/*.wasm ./output
