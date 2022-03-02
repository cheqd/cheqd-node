# Guide to Debian packages for cheqd-node

## Context

We provide [pre-compiled Debian packages](https://github.com/cheqd/cheqd-node/releases) for ease-of-installation and configuration of cheqd-node on standalone virtual machines / hosts.

These are tested to work with Ubuntu 20.04 LTS, which is the current long-term support release for Ubuntu Linux.

It may be possible to use the same package on other Ubuntu / Debian distributions based on the same LTS version; however, this method is not officially supported.

## How-to instructions for `cheqd-node` Debian packages

* [Install (or uninstall) a new node](deb-package-install.md) using Debian package
* If you already have an existing `cheqd-node` installation that was done using the .deb, find out how to [upgrade your node using the Debian package](deb-package-upgrade.md).

## Pre-install actions executed by Debian package

This section describes the system changes that our Debian packages attempt to carry out, including:

* Pre-install system configuration actions.
* Installing pre-compiled `cheqd-node` binary.
* Post-install actions, such as setting up the `cheqd-noded` daemon as a `systemctl` service.

### System user creation

By default, Cosmos SDK creates all requisite directories in the `HOME` directory of the user that initiates installation.

Our package creates a new system user called `cheqd` with home directory set to `/home/cheqd`. This allows node operators to keep sysadmin / standard users separate from the service user. Home directory can be changed by passing a variable called `CHEQD_HOME_DIR` before executing installation.

### Directory location configuration

To keep `cheqd-node` configuration data in segregated from userspace home directories, the installer creates new application data directories and symbolic links.

#### App data directories

* `$HOME/.cheqdnode/config`
  * Location for configuration files
  * Ownership permission set to: `cheqd:cheqd`
* `$HOME/.cheqdnode/data`
  * Location for ledger data
  * Ownership permission set to: `cheqd:cheqd`

### Logging configuration

The default logging location and permissions are as follows:

* `$HOME/.cheqdnode/log`
  * Location for app logs
  * Ownership permissions set to: `syslog:cheqd` (set by rsyslog)

The log location can be overridden by passing the variable `CHEQD_LOG_DIR` before executing the installation proceess.

`rsyslog` is configured to redirect logs from the `cheqd-node` daemon to the log directory defined above.


```bash
if \$programname == 'cheqd-noded' then ${CHEQD_LOG_DIR}/stdout.log
& stop
```

#### Log rotation configuration

Log rotation is achieved using the system `logrotate` service.

Our installer makes the following changes:

* Sets the maximum filesize of the `stdout.log` file is set to 100 MB, after which the log file is compressed and stored separately. The `stdout.log` file then continues with storing newer log entries.
* Archives logs are deleted after 30 days.
* Once a day, uses the `logrotate.timer` service to rotate logs.

```bash
${CHEQD_LOG_DIR}/stdout.log {
  rotate 30
  daily
  maxsize 100M
  notifempty
  copytruncate
  compress
  maxage 30
}
```

## Post-install actions executed by Debian package

The main part of post-installation process is to make the `cheqd-node` binary run as a `systemctl` daemon.

This ensures the service is restarted after any failures and output sent to `rsyslog`.

```text
[Unit]
Description=Service for running Cheqd node
After=network.target
[Service]
Type=simple
User=cheqd
ExecStart=/bin/bash -c '/usr/bin/cheqd-noded start'
Restart=on-failure
RestartSec=10
StartLimitBurst=10
StartLimitInterval=200
TimeoutSec=300
StandardOutput=syslog
StandardError=syslog
SyslogFacility=syslog
SyslogIdentifier=cheqd-noded
[Install]
WantedBy=multi-user.target
```

## Uninstalling the Debian package

| :warning: WARNING |
| :--- |
| Please make sure any accounts keys are backed up or exported before attempting uninstallation |

To uninstall `cheqd-node` when it has been installed using the Debian package release, execute the following (with `sudo` or as the `root` user):

```bash
apt remove cheqd-node
```

This will remove all configuration files created during installation process from the system, such as:

```text
/etc/rsyslog.d/cheqd-node.conf
/etc/logrotate.d/cheqd-node
/lib/systemd/system/cheqd-noded.service
```
