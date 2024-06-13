# Rust Topology

This monorepo implements Topology Protocol in Rust.

## Running `ramd`

To build `ramd`, Rust (version 1.77.2 or later) is required. Please make sure Rust is installed and then run:

```
make run-ramd
```

It will run `ramd` and generate default config and necessary files. To view these files, use:

```
cd $HOME/.ramd
ls
```

(Optional) If you wish to change the location of these files, you can set `RAMD_DIR_NAME` environment variable in .env file as demonstrated [here](./.env.example). After configuring .env and running `ramd`, you can access the files with:

```
cd $HOME/{RAMD_DIR_NAME}
ls
```

If you encounter any issues, please feel free to [reach out](#contact) to us.

## Contributing

We are committed to community-driven development and welcome feedback and contributions from anyone on the internet!

If you're interested in collaborating with us, please refer to [CONTRIBUTING.md](./CONTRIBUTING.md) for more details.

## Contact

You can join our [Discord](https://discord.gg/hMsQas3Vw9) to ask questions or engage in discussions.

## License

RAM monorepo is licensed under the MIT License.
