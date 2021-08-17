# Overview
This document describes how to make package upgrading for the Node with `cheqd-node` package inside.
## Get the package
The latest package can be found in [releases](https://github.com/cheqd/cheqd-node/releases);
## Steps to upgrade package locally
1. First of all, please make sure that service is stopped and make backups of data and keys before package upgrading.
For stopping service can be used this command:
```
systemctl stop cheqd-noded
```
and 
```
systemctl status cheqd-noded
```
for checking that service was stopped.
2. After that new package can be installed by calling:
```
apt install <path/to/package>
```
## Needed checks after installation
* Start `cheqd-noded` service by calling:
  ```
  systemctl start cheqd-noded
  ```
  and check the state of it:
  ```
  systemctl status cheqd-noded
  ```
* After installation complete it would be great to check that all data and configs were kept the same. 
Package installation process is idempotent and it should not change service files, configs or any other user data.

P.S. In case of using just binary installation before, we recommend to copy all `config` into `/etc/cheqd-node` and `data` into the `/var/lib/cheqd/data`.
Also, additional information about debian package can be find [here](deb-package-overview.md).