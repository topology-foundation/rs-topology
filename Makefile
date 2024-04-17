.PHONY: build-workspace clippy lint-check

build-workspace:
	cargo build

clippy:
	cargo clippy -- -D warnings

lint-check:
	cargo fmt --all -- --check