<h1 align="center">Project Eudico</h1>

Eudico is a modularised implementation of [Lotus](https://github.com/filecoin-project/lotus), itself an implementation of the Filecoin Distributed Storage Network. For more details about Filecoin, check out the [Filecoin Spec](https://spec.filecoin.io). This is a work-in-progress, intended to enable easier experimentation with future protocol features, and is not meant to be used in the production network.

[![consensus-shipyard](https://circleci.com/gh/consensus-shipyard/lotus.svg?style=svg)](https://app.circleci.com/pipelines/github/consensus-shipyard/lotus)

## Building & Documentation

> Note: The `dev` branch is the cutting-edge implementation, recommended for subnets, while the `spacenet` branch guarantees Spacenet compatibility.
 
For complete instructions on how to build, install and setup eudico/lotus, please visit [https://lotus.filecoin.io](https://lotus.filecoin.io/lotus/install/prerequisites/#supported-platforms). Basic build instructions can be found further down in this readme.

## Reporting a Vulnerability

Please send an email to security@filecoin.org. See our [security policy](SECURITY.md) for more details.

## Basic Build Instructions
**System-specific Software Dependencies**:

Building eudico requires some system dependencies, usually provided by your distribution.

Ubuntu/Debian:
```
sudo apt install mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl clang build-essential hwloc libhwloc-dev wget -y && sudo apt upgrade -y
```

Fedora:
```
sudo dnf -y install gcc make git bzr jq pkgconfig mesa-libOpenCL mesa-libOpenCL-devel opencl-headers ocl-icd ocl-icd-devel clang llvm wget hwloc hwloc-devel
```

For other distributions you can find the required dependencies [here.](https://lotus.filecoin.io/lotus/install/prerequisites/#supported-platforms) For instructions specific to macOS, you can find them [here.](https://lotus.filecoin.io/lotus/install/macos/)

#### Go

To build eudico, you need a working installation of [Go 1.19.7 or higher](https://golang.org/dl/):

```bash
wget -c https://golang.org/dl/go1.19.7.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
```

**TIP:**
You'll need to add `/usr/local/go/bin` to your path. For most Linux distributions you can run something like:

```shell
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc && source ~/.bashrc
```

See the [official Golang installation instructions](https://golang.org/doc/install) if you get stuck.

### Build and install Lotus

Once all the dependencies are installed, you can build and install eudico.

   ```sh
   git clone https://github.com/consensus-shipyard/lotus.git
   cd lotus/
   make spacenet
   ```
 
Eudico will use the `$HOME/.lotus` folder by default for storage (configuration, chain data, wallets, etc). See [advanced options](https://lotus.filecoin.io/lotus/configure/defaults/#environment-variables) for information on how to customize the Lotus folder.

You should now have Lotus installed. You can now [start the Lotus daemon and sync the chain](https://lotus.filecoin.io/lotus/install/linux/#start-the-lotus-daemon-and-sync-the-chain).

## License

Dual-licensed under [MIT](https://github.com/filecoin-project/lotus/blob/master/LICENSE-MIT) + [Apache 2.0](https://github.com/filecoin-project/lotus/blob/master/LICENSE-APACHE)
