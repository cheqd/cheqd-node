#!/usr/bin/env python3


# Python package imports
import datetime
import os
import subprocess
import sys
import urllib.error
from pathlib import Path
import urllib.request as request
import json
import re
import functools
import pwd
import shutil
import signal
import time
import threading


# Installation Parameters
LAST_N_RELEASES = 5
DEFAULT_HOME = "/home/cheqd"
DEFAULT_INSTALL_PATH = "/usr/bin"
DEFAULT_CHEQD_USER = "cheqd"
DEFAULT_BINARY_NAME = "cheqd-noded"
DEFAULT_CHAINS = ['testnet', 'mainnet']
DEFAULT_CHAIN = "mainnet"
PRINT_PREFIX = "********* "


### Cosmovisor Config
COSMOVISOR_BINARY = "https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz"
DEFAULT_USE_COSMOVISOR = "yes"


### Genesis and Seeds
GENESIS_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/networks/{}/genesis.json"
SEEDS_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/networks/{}/seeds.txt"

###############################################################
###     				Node snapshots      				###
###############################################################
DEFAULT_SNAPSHOT_SERVER = "https://snapshots.cheqd.net"
DEFAULT_INIT_FROM_SNAPSHOT = "yes"
TESTNET_SNAPSHOT = "https://cheqd-node-backups.ams3.cdn.digitaloceanspaces.com/testnet/latest/cheqd-testnet-4_{}.tar.gz"
MAINNET_SNAPSHOT = "https://cheqd-node-backups.ams3.cdn.digitaloceanspaces.com/mainnet/latest/cheqd-mainnet-1_{}.tar.gz"

###############################################################
###     				Systemd Config      				###
###############################################################
SERVICE_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/tools/build/cheqd-noded.service"
SERVICE_FILE_PATH = "/lib/systemd/system/cheqd-noded.service"
DEFAULT_LOGROTATE_FILE = "/etc/logrotate.d/cheqd-node"
DEFAULT_RSYSLOG_FILE = "/etc/rsyslog.d/cheqd-node.conf"


def sigint_handler(signal, frame):
    print ('Exiting from cheqd-node installer')
    sys.exit(0)

signal.signal(signal.SIGINT, sigint_handler)

class Release:
    def __init__(self, release_map):
        self.version = release_map['tag_name']
        self.url = release_map['html_url']
        self.assets = release_map['assets']

    def get_tar_gz_url(self):
        archive_urls = [ a['browser_download_url'] for a in self.assets if a['browser_download_url'].find("tar.gz") > 0]
        if len(archive_urls) == 0:
            failure_exit(f"No tar.gz in release: {self.version}")
        return archive_urls[0]

    def get_binary_url(self):
        binary_urls = [ a['browser_download_url'] for a in self.assets if a['name'] == DEFAULT_BINARY_NAME]
        if len(binary_urls) == 0:
            failure_exit(f"No binaries in release: {self.version}")
        return binary_urls[0]

    def __str__(self):
        return f"Name: {self.version}, Tar URL: {self.get_tar_gz_url()}"


def failure_exit(reason):
    print(f"Reason of failure: {reason}")
    print("Exiting....")
    sys.exit(1)


def post_process(func):
    @functools.wraps(func)
    def wrapper(*args, **kwds):
        _allow_error = kwds.get('allow_error', False)
        try:
            value = func(*args)
        except subprocess.CalledProcessError as err:
            if err.returncode and _allow_error:
                return err
            failure_exit(err)
        return value

    return wrapper


def default_answer(func):
    @functools.wraps(func)
    def wrapper(*args, **kwds):
        _default = kwds.get('default', "")
        if _default:
            args = list(args)
            args[-1] += f"[{_default}] {os.linesep}"
        value = func(*args)
        return value if value != "" else _default

    return wrapper


class Installer():
    def __init__(self, interviewer):
        self.version = interviewer.release.version
        self.release = interviewer.release
        self.verbose = True
        self.interviewer = interviewer

    @property
    def binary_path(self):
        return self.get_binary_path()

    def get_binary_path(self):
        return os.path.join(os.path.realpath(os.path.curdir), DEFAULT_BINARY_NAME)

    @property
    def cosmovisor_service_cfg(self):
        return f"""
[Unit]
Description=Service for running cheqd-node daemon
After=network.target
Documentation=https://docs.cheqd.io/node

[Service]
Environment="DAEMON_HOME={self.cheqd_root_dir}"
Environment="DAEMON_NAME={DEFAULT_BINARY_NAME}"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=true"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="UNSAFE_SKIP_BACKUP=true"
Type=simple
User=cheqd
ExecStart=/usr/bin/cosmovisor run start
Restart=on-failure
RestartSec=30
StartLimitBurst=5
StartLimitInterval=60
TimeoutSec=120
StandardOutput=syslog
StandardError=syslog
SyslogFacility=syslog
SyslogIdentifier=cosmovisor
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
"""

    @property
    def default_logrotate_cfg(self):
        return """
%s/stdout.log {
  rotate 7
  daily
  maxsize 100M
  notifempty
  copytruncate
  compress
  maxage 7
}
""" % self.cheqd_log_dir

    @property
    def default_rsyslog_cfg(self):
        binary_name = "cosmovisor" if self.interviewer.is_cosmo_needed else "cheqd-noded"
        return f"""
if $programname == '{binary_name}' then {self.cheqd_log_dir}/stdout.log
& stop
"""

    @property
    def cheqd_root_dir(self):
        return os.path.join(self.interviewer.home_dir, ".cheqdnode")

    @property
    def cheqd_config_dir(self):
        return os.path.join(self.cheqd_root_dir, "config")

    @property
    def cheqd_data_dir(self):
        return os.path.join(self.cheqd_root_dir, "data")

    @property
    def cheqd_log_dir(self):
        return os.path.join(self.cheqd_root_dir, "log")

    @property
    def cosmovisor_root_dir(self):
        return os.path.join(self.cheqd_root_dir, "cosmovisor")

    @property
    def cosmovisor_cheqd_bin_path(self):
        return os.path.join(self.cosmovisor_root_dir, f"current/bin/{DEFAULT_BINARY_NAME}")

    def log(self, msg):
        if self.verbose:
            print(f"{PRINT_PREFIX} {msg}")

    @post_process
    def exec(self, cmd):
        self.log(f"Executing command: {cmd}")
        return subprocess.run(cmd, shell=True, check=True, capture_output=True)

    def get_binary(self):
        if self.release.version <= LAST_VERSION_WITH_TARBALL:
            self.exec(f"wget -qO - {self.release.get_tar_gz_url()}  | tar xz")
        else:
            self.exec(f"wget -qo cheqd-noded {self.release.get_binary_url()}")

    def is_user_exists(self, username) -> bool:
        try:
            pwd.getpwnam(username)
        except KeyError:
            self.log(f"User {username} does not exist")
            return False
        self.log(f"User {username} already exists")
        return True

    def is_already_installed(self) -> bool:
        return self.is_user_exists(DEFAULT_CHEQD_USER)

    def install(self):
        if not self.is_already_installed():
            self.log("Download the binary")
            self.get_binary()
        else:
            failure_exit("Looks like installation already exists.")

        if not self.interviewer.is_cosmo_needed:
            self.log(f"Moving binary from {self.binary_path} to {DEFAULT_INSTALL_PATH}")
            self.exec("sudo mv {} {}".format(self.binary_path, DEFAULT_INSTALL_PATH))

        if not self.is_user_exists(DEFAULT_CHEQD_USER):
            self.log(f"Create a user {DEFAULT_CHEQD_USER} cause it's not created yet")
            self.prepare_cheqd_user()

        if not os.path.exists(self.cheqd_root_dir):
            self.log("Make root directory for cheqd-node")
            self.mkdir_p(self.cheqd_root_dir)

            self.log(f"Chown to default cheqd user: {DEFAULT_CHEQD_USER}")
            self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cheqd_root_dir}")

        if not os.path.exists(self.cheqd_log_dir):
            self.log("Setup log directory")
            self.setup_log_dir()

        if os.path.exists("/var/log/cheqd-node") and not os.path.islink("/var/log/cheqd-node"):
            self.log("Make a link to /var/log/cheqd-node")
            os.symlink(self.cheqd_log_dir, "/var/log/cheqd-node", target_is_directory=True)

        if os.path.exists("/etc/rsyslog.d/"):
            if not os.path.exists(DEFAULT_RSYSLOG_FILE):
                self.log("Configure rsyslog")
                with open(DEFAULT_RSYSLOG_FILE, mode="w") as fd:
                    fd.write(self.default_rsyslog_cfg)
                # Sometimes it can take a lot of time: https://github.com/rsyslog/rsyslog/issues/3133
                self.exec("systemctl restart rsyslog")

        if not os.path.exists(DEFAULT_LOGROTATE_FILE):
            if not os.path.exists(DEFAULT_LOGROTATE_FILE):
                self.log("Add config for logrotation")
                with open(DEFAULT_LOGROTATE_FILE, mode="w") as fd:
                    fd.write(self.default_logrotate_cfg)
                # Sometimes it can take a lot of time: https://github.com/rsyslog/rsyslog/issues/3133
                self.exec("systemctl restart rsyslog")

        self.log("Restart logrotate services")
        self.exec("systemctl restart logrotate.service")
        self.exec("systemctl restart logrotate.timer")

        self.log("Setup systemctl service config")
        if self.interviewer.is_cosmo_needed:
            with open(SERVICE_FILE_PATH, mode="w") as fd:
                fd.write(self.cosmovisor_service_cfg)
        else:
            self.exec(f"curl -s {SERVICE_FILE} > {SERVICE_FILE_PATH}")

        self.log("Enable systemctl service")
        self.exec("systemctl enable cheqd-noded")

        if self.interviewer.is_cosmo_needed:
            self.log("Setup the cosmovisor")
            self.setup_cosmovisor()

        if self.interviewer.is_setup_needed:
            self.post_install()

        if self.interviewer.init_from_snapshot:
            self.log("Going to download the archive and untar it on a fly. It can take a really LONG TIME")
            self.untar_from_snapshot()

    def post_install(self):
        # Init the node with provided moniker
        if not os.path.exists(os.path.join(self.cheqd_config_dir, 'genesis.json')):
            self.exec(f"sudo -u {DEFAULT_CHEQD_USER} cheqd-noded init {self.interviewer.moniker}")

            # Downloading genesis file
            self.exec(f"curl -s {GENESIS_FILE.format(self.interviewer.chain)} > {os.path.join(self.cheqd_config_dir, 'genesis.json')}")
            shutil.chown(os.path.join(self.cheqd_config_dir, 'genesis.json'),
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)

        # Setting up the external_address
        if self.interviewer.external_address:
            self.exec(f"sudo -u {DEFAULT_CHEQD_USER} cheqd-noded configure p2p external-address {self.interviewer.external_address}")

        # Setting up the seeds
        seeds = self.exec(f"curl -s {SEEDS_FILE.format(self.interviewer.chain)}").stdout.decode("utf-8").strip()
        self.exec(f"sudo -u {DEFAULT_CHEQD_USER} cheqd-noded configure p2p seeds {seeds}")

    def prepare_cheqd_user(self):
        self.log(f"Create group, {DEFAULT_CHEQD_USER} by default")
        self.exec(f"addgroup {DEFAULT_CHEQD_USER} --quiet")

        self.log(f"Create user, {DEFAULT_CHEQD_USER} by default")
        self.exec(
            f"adduser --system {DEFAULT_CHEQD_USER} --home {self.interviewer.home_dir} --shell /bin/bash --ingroup {DEFAULT_CHEQD_USER} --quiet")

    def mkdir_p(self, dir_name):
        try:
            os.mkdir(dir_name)
        except FileExistsError as err:
            self.log(f"Directory {dir_name} already exists")

    def setup_log_dir(self):
        self.mkdir_p(self.cheqd_log_dir)
        Path(os.path.join(self.cheqd_log_dir, "stdout.log")).touch()
        self.exec(f"chown -R syslog:syslog {self.cheqd_log_dir}")

    def setup_cosmovisor(self):
        self.exec(f"wget -qO - {COSMOVISOR_BINARY}  | tar xz")
        self.mkdir_p(self.cosmovisor_root_dir)
        self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "genesis"))
        self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "genesis/bin"))
        self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "upgrades"))
        if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, "cosmovisor")):
            self.log(f"Moving cosmovisor binary to default installation directory")
            shutil.move("./cosmovisor", DEFAULT_INSTALL_PATH)

        self.log(f"Moving binary from {self.binary_path} to {self.cosmovisor_cheqd_bin_path}")
        self.exec("sudo mv {} {}".format(self.binary_path, self.cosmovisor_cheqd_bin_path))

        if not os.path.exists(os.path.join(self.cosmovisor_root_dir, "current")):
            self.log(f"Making symlink current -> genesis")
            os.symlink(os.path.join(self.cosmovisor_root_dir, "genesis"),
                       os.path.join(self.cosmovisor_root_dir, "current"))

        if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME)):
            self.log(f"Making symlink to {self.cosmovisor_cheqd_bin_path}")
            os.symlink(self.cosmovisor_cheqd_bin_path,
                       os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))

        self.log(f"Changing owner to {DEFAULT_CHEQD_USER} user")
        self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cosmovisor_root_dir}")

    def untar_from_snapshot(self):
        self.mkdir_p(self.cheqd_data_dir)
        cmd = f"wget -c -O - {self.interviewer.snapshot_url}  | sudo -u {DEFAULT_CHEQD_USER} tar xzf - -C {os.path.join(self.cheqd_root_dir, 'data')}"
        thread = threading.Thread(target=functools.partial(self.exec, cmd))
        thread.start()
        sec_counter = 0

        # wait small period of time for waiting the command running
        time.sleep(3)
        while thread.is_alive():
            time.sleep(60)
            sec_counter += 60
            self.log(f"Downloading is alive, it already took: {str(datetime.timedelta(seconds=sec_counter))}")

        # Some kind of hacks cause cosmovisor expects upgrade-info.json file in cosmovisor/current directory also
        if self.interviewer.is_cosmo_needed:
            if os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")):
                self.log(f"Copying upgrade-info.json file to cosmovisor/current/")
                shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                            os.path.join(self.cosmovisor_root_dir, "current"))

            self.log(f"Changing owner to {DEFAULT_CHEQD_USER} user")
            self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cosmovisor_root_dir}")

        self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cheqd_data_dir}")


class Interviewer:
    def __init__(self,
                 home_dir=DEFAULT_HOME,
                 chain=DEFAULT_CHAIN):
        self._home_dir = home_dir
        self._is_cosmo_needed = True
        self._init_from_snapshot = False
        self._release = None
        self._chain = chain
        self.verbose = True
        self._snapshot_url = self.prepare_url_for_latest()
        self._is_setup_needed = False
        self._moniker = ""
        self._external_address = ""

    @property
    def release(self) -> Release:
        return self._release

    @property
    def home_dir(self) -> str:
        return self._home_dir

    @property
    def is_cosmo_needed(self) -> bool:
        return self._is_cosmo_needed

    @property
    def init_from_snapshot(self) -> bool:
        return self._init_from_snapshot

    @property
    def chain(self) -> str:
        return self._chain

    @property
    def snapshot_url(self) -> str:
        return self._snapshot_url

    @property
    def is_setup_needed(self) -> bool:
        return self._is_setup_needed

    @property
    def moniker(self) -> str:
        return self._moniker

    @property
    def external_address(self) -> str:
        return self._external_address

    @release.setter
    def release(self, release):
        self._release = release

    @home_dir.setter
    def home_dir(self, hd):
        self._home_dir = hd

    @is_cosmo_needed.setter
    def is_cosmo_needed(self, icn):
        self._is_cosmo_needed = icn

    @init_from_snapshot.setter
    def init_from_snapshot(self, ifs):
        self._init_from_snapshot = ifs

    @chain.setter
    def chain(self, chain):
        self._chain = chain

    @snapshot_url.setter
    def snapshot_url(self, su):
        self._snapshot_url = su

    @is_setup_needed.setter
    def is_setup_needed(self, is_setup_needed):
        self._is_setup_needed = is_setup_needed

    @moniker.setter
    def moniker(self, moniker):
        self._moniker = moniker

    @external_address.setter
    def external_address(self, external_address):
        self._external_address = external_address

    def log(self, msg):
        if self.verbose:
            print(f"{PRINT_PREFIX} {msg}")

    def get_releases(self):
        req = request.Request("https://api.github.com/repos/cheqd/cheqd-node/releases")
        req.add_header("Accept", "application/vnd.github.v3+json")

        with request.urlopen(req) as response:
            r_list = json.loads(response.read().decode("utf-8"))
            return [Release(r) for r in r_list]

    def r_filter(self, releases):
        return [r for r in releases if ONLY_DIGIT_VERSIONS.match(r.version)]

    def ask_for_version(self):
        all_releases = self.get_releases()
        last_n_releases = self.r_filter(all_releases)[:LAST_N_RELEASES]
        default = last_n_releases[0].version
        answer = self.ask(
            f"Which version do you want to install? Or type 'list' for get the list of releases: ", default=default)
        if answer == default:
            self.release = last_n_releases[0]
        elif answer == "list":
            for i, release in enumerate(last_n_releases):
                print(f"{i + 1}) {release.version}")
            try:
                num = int(input("Please insert the number for picking up the version: "))
            except ValueError as err:
                failure_exit("Version number should be integer value.")
            self.release = last_n_releases[num - 1]
        else:
            if answer[0] != "v":
                answer = f"v{answer}"
            _t = [a for a in all_releases if a.version == answer]
            if len(_t) > 0:
                self.release = _t[0]
            else:
                failure_exit(f"Version: {answer} does not exist")

    @default_answer
    def ask(self, question, **kwargs):
        return str(input(question))

    def ask_for_home_directory(self, default) -> str:
        answer = self.ask(
            f"Please, type here the path to home directory for user cheqd. For keeping default value, just type "
            f"'Enter': ", default=default)
        self.home_dir = answer

    def ask_for_cosmovisor(self, default) -> str:
        answer = self.ask(
            f"Do you want to use Cosmovisor? "
            f"Please type any kind of variants: yes, no, y, n. ", default=default)
        self.is_cosmo_needed = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_init_from_snapshot(self):
        answer = self.ask(
            f"Do you want to deploy the latest snapshot? "
            f"Please type any kind of variants: yes, no, y, n. ",
            default="No"
        )
        self.init_from_snapshot = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_chain(self):
        answer = self.ask(
            f"Which chain do you want to use? Possible variants are: {', '.join(DEFAULT_CHAINS)} ",
            default="testnet"
        )
        self.chain = answer if answer in DEFAULT_CHAINS else failure_exit(f"Possible chains are: {DEFAULT_CHAINS}")

    def ask_for_snapshot_url(self):
        answer = self.ask(
            f"Which snapshot do you want to use? Please type the full URL to archive or press return to use the latest ",
            default=self.prepare_url_for_latest()
        )
        self.snapshot_url = answer

    def ask_for_setup(self):
        answer = self.ask(
            f"Do you want to setup node after installation? "
            f"Please type any kind of variants: yes, no, y, n. ",
            default="No"
        )
        self.is_setup_needed = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_moniker(self):
        answer = self.ask(
            f"Please, type the moniker for your node: {os.linesep}",
            default=""
        )
        self.moniker = answer

    def ask_for_external_address(self):
        answer = self.ask(
            f"What are the external IP address and the P2P port (default is 26656) for your node? Please type in format: <ip_address>:<port>{os.linesep}",
            default=""
        )
        self.external_address = answer

    def prepare_url_for_latest(self) -> str:
        template = TESTNET_SNAPSHOT if self.chain == "testnet" else MAINNET_SNAPSHOT
        _date = datetime.date.today()
        _url = template.format(_date.strftime("%Y-%m-%d"))
        while not self.is_url_exists(_url):
            _date -= datetime.timedelta(days=1)
            _url = template.format(_date.strftime("%Y-%m-%d"))
        return _url

    def is_url_exists(self, url):
        try:
            request.urlopen(request.Request(url))
        except urllib.error.HTTPError:
            return False
        return True

if __name__ == '__main__':
    # Ask user for information
    interviewer = Interviewer()
    interviewer.ask_for_version()
    interviewer.ask_for_home_directory(default=DEFAULT_HOME)
    interviewer.ask_for_cosmovisor(default=DEFAULT_USE_COSMOVISOR)
    interviewer.ask_for_chain()
    interviewer.ask_for_init_from_snapshot()
    if interviewer.init_from_snapshot:
        interviewer.ask_for_snapshot_url()

    interviewer.ask_for_setup()
    if interviewer.is_setup_needed:
        interviewer.ask_for_moniker()
        interviewer.ask_for_external_address()

    # Install
    installer = Installer(interviewer)
    installer.install()
