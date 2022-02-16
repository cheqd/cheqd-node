# Prepare environment

For now, we need to separate storing DID in the wallet and send the full DID-Doc to the cheqd network.
Propagating DID-doc to the cheqd network can be dividd into 2 steps:
- DID manipulation, like create and update operations locally can be implemented by using VDR tools from [here](https://gitlab.com/evernym/verity/vdr-tools)
- DID sending to the network can be made by creating a json string from DID stored in the VDR wallet and passed it to the `cheqd-noded` binary.

# VDR tools installation

In general this process is described on the main page of VDR tools [repository](https://gitlab.com/evernym/verity/vdr-tools#installing) but to be short let's make the next steps inside Ubuntu 20.04:
1. Install additional packages:
   ```
   apt install curl libsodium23 libzmq5 libncursesw5-dev -y
   ```
2. Download and install `libvdrtools`:
    ```
    curl https://gitlab.com/evernym/verity/vdr-tools/-/package_files/27311917/download --output libvdrtools_0.8.4-focal_amd64.deb && dpkg -i libvdrtools_0.8.4-focal_amd64.deb
    ```
3. Download and install `vdrtools-cli`:
    ```
    curl https://gitlab.com/evernym/verity/vdr-tools/-/package_files/27311922/download --output vdrtools-cli_0.8.4-focal_amd64.deb && dpkg -i vdrtools-cli_0.8.4-focal_amd64.deb
    ```