# Guide to Debian packages for `cheqd-node`

## Context

We provide [pre-compiled Debian packages](https://github.com/cheqd/cheqd-node/releases) for ease-of-installation and configuration of cheqd-node on standalone virtual machines / hosts.

These are tested to work with Ubuntu 20.04 LTS, which is the current long-term support release for Ubuntu Linux.

It may be possible to use the same package on other Ubuntu / Debian distributions based on the same LTS version; however, this method is not officially supported.


## How-to instructions for `cheqd-node` Debian packages

* [Install (or uninstall) a new node](deb-package-install.md) using Debian package
* If you already have an existing `cheqd-node` installation that was done using the .deb, find out how to [upgrade your node using the Debian package](deb-package-upgrade.md).


## Pre/post-install action details

This section describes the system changes that our Debian packages attempt to carry out, including:

* Pre-install system configuration actions.
* Installing pre-compiled `cheqd-node` binary.
* Post-install actions, such as setting up the `cheqd-node` daemon as a `systemctl` service.

### System user creation

By default, Cosmos SDK creates all requisite directories in the `HOME` directory of the user that initiates installation.

Our package creates a new system user called `cheqd` with home directory set to `/var/lib/cheqd`. This allows node operators to keep sysadmin / standard users separate from the service user.

### Application directory and symbolic link creation

To keep `cheqd-node` configuration data in 

#### App data directories

* `/etc/cheqd-node`
  * Configuration files location
  * Permissions: `cheqd:cheqd`
* `/var/lib/cheqd/data`
  * Place for blockchain data
  * Permissions: cheqd:cheqd
* `/var/log/cheqd-node`
  * Place for logs
  * Permissions: syslog:adm \(set by rsyslog\)

The following symlinks will be created:

* `/etc/cheqd-node/` -&gt; `/var/lib/cheqd/.cheqdnode/config`
  * For configs
* `/var/lib/cheqd/data` -&gt; `/var/lib/cheqd/.cheqdnode/data`
  * For data

### Rsyslog configuration

The next config for rsyslog will be created:

```text
if \$programname == 'cheqd-noded' then /var/log/cheqd-node/stdout.log
& stop
```

It redirects all the logs into the file.

### Logrotate config

For rotating log file will be used `logrotate` - the general approach for Linux/systemd with the following config:

```text
/var/log/cheqd-node/stdout.log {
  rotate 30
  maxsize 100M
  notifempty
  copytruncate
  compress
  maxage 30
}
```

It means, that log will be rotated after achieving 100 Mb size of `stdout.log` and compressed. All the archives will be stored for a month \(30 days\). Also, the main file will truncated instead of removing. It needs for continue logging process in terms of file pointers.

Once a day by crontab will be called a small script for running logrotate logic.

## Systemd

The main part of post-installation process is making our binary as a service. The following systemd service file will be created:

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

The main thing here is that it will restart on binary failures and put all output to the `rsyslog`.

## Post-remove actions

For now, all files created during installation process will be removed from the system, like:

```text
/etc/rsyslog.d/cheqd-node.conf
/etc/logrotate.d/cheqd-node
/etc/cron.daily/cheqd-node
/etc/systemd/system/cheqd-noded.service
```