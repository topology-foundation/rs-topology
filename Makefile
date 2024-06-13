.PHONY: build-workspace clippy lint-check run-ramd clean-ramd-dir

build-workspace:
	cargo build

clippy:
	cargo clippy -- -D warnings

lint-check:
	cargo fmt --all -- --check

run-ramd:
	RUST_LOG=info cargo run --bin ramd node

clean-ramd-dir:
	rm -r ${HOME}/.ramd
