# Overview of changes that deb package does

Debian package is the most common way for Ubuntu OS and also it can provide post-install steps which can help to make our application to run as a service.

The package consists of:

* Binary, named `cheqd-noded`;
* Script with post-install and post-remove actions. 

## Post-install actions

### System user creation

By default, cosmos-sdk create all needed directories in the `HOME` directory. That's why package creates a special `cheqd` user with home directory set to `/var/lib/cheqd`.

### Directories and symlinks

According to general filesystem hierarchy standard \(FHS\), the next directories will be created:

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

