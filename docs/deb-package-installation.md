## Overview
One of the possible distribution way is to use debian package. 
It's the most common way for Ubuntu OS and also it can provide post-install steps which can help to make our application to run as a service.
Because we have to run init procedure for creating tree of directories, `configs` and `data`, not all the steps can be automated.
By the way, debian package consists of binary, named `cheqd-noded` and script with post-install and post-remove actions. 

## Post-install actions
### Create a special user "cheqd"
By default, cosmos-sdk create all needed directories in the `HOME` directory. 
That's why package creates a special user with home directory `/var/lib/cheqd`. Also, this user will use for setting permissions to data and configs.

### Dividing configs, data and logs
#### Directories
According to general filesystem hierarchy standard (FHS), the next directories will be created:
```
/etc/cheqd-node                - configs, permissions cheqd:cheqd
/var/lib/cheqd/data            - data   , permissions cheqd:cheqd
/var/log/cheqd-node            - logs   , permissions syslog:adm (set by rsyslog)
```

After setting up the node, it's expected, then configs and data will be symlinked to the corresponded system directories.
For this purposes will be created the next symlinks to configs and data:
```
sudo ln -s /etc/cheqd-node/ /var/lib/cheqd/.cheqdnode/config   - for configs
sudo ln -s /var/lib/cheqd/data /var/lib/cheqd/.cheqdnode/      - for data
```

After this preparation, it would be possible to set up cheqd node in general but under `cheqd` user.

#### Rsyslog config
The next config for rsyslog will be created:
```
if \$programname == 'cheqd-noded' then /var/log/cheqd-node/stdout.log
& stop
```
It redirects all the logs into the file.
#### Logrotate config
For rotating log file will be used `logrotate` - the general approach for Linux/systemd with the next config:
```
/var/log/cheqd-node/stdout.log {
  rotate 30
  maxsize 100M
  notifempty
  copytruncate
  compress
  maxage 30
}
```
It means, that log will be rotated after achieving 100 Mb size of `stdout.log` and compressed. 
All the archives will be stored for a month (30 days). Also, the main file will truncated instead of removing. It needs for continue logging process in terms of file pointers.

Once a day by crontab will be called a small script for running logrotate logic.

### Systemd
The main part of post-installation process is making our binary as a service. Systemd service file can help with it:
```
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

## Exposing port


## Post-remove actions
For now, all files created during installation process will be removed from the system, like:
```
/etc/rsyslog.d/cheqd-node.conf
/etc/logrotate.d/cheqd-node
/etc/cron.daily/cheqd-node
/etc/systemd/system/cheqd-noded.service
```
