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
COSMOVISOR_BINARY_URL = "https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz"
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
STANDALONE_SERVICE_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/b0e5ee4de6daf44f2dd8e49d5ed4a38b3c299873/tools/build/node-standalone.service"
COSMOVISOR_SERVICE_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/57978fe6dcac0634ad22936c9cc98bf968e573c6/tools/build/node-cosmovisor.service"

LOGROTATE_TEMPLATE = "https://raw.githubusercontent.com/cheqd/cheqd-node/57978fe6dcac0634ad22936c9cc98bf968e573c6/tools/build/logrotate.conf"
RSYSLOG_TEMPLATE = "https://raw.githubusercontent.com/cheqd/cheqd-node/57978fe6dcac0634ad22936c9cc98bf968e573c6/tools/build/rsyslog.conf"

DEFAULT_STANDALONE_SERVICE_NAME = 'cheqd-noded'
DEFAULT_COSMOVISOR_SERVICE_NAME = 'cheqd-cosmovisor'

DEFAULT_STANDALONE_SERVICE_FILE_PATH = f"/lib/systemd/system/{DEFAULT_STANDALONE_SERVICE_NAME}.service"
DEFAULT_COSMOVISOR_SERVICE_FILE_PATH = f"/lib/systemd/system/{DEFAULT_COSMOVISOR_SERVICE_NAME}.service"
DEFAULT_LOGROTATE_FILE = "/etc/logrotate.d/cheqd-node"
DEFAULT_RSYSLOG_FILE = "/etc/rsyslog.d/cheqd-node.conf"

### Default parameters

DEFAULT_RPC_PORT = "26657"
DEFAULT_P2P_PORT = "26656"
DEFAULT_GAS_PRICE = "25ncheq"


def sigint_handler(signal, frame):
    print ('Exiting from cheqd-node installer')
    sys.exit(0)

signal.signal(signal.SIGINT, sigint_handler)


class Release:
    def __init__(self, release_map):
        self.version = release_map['tag_name']
        self.url = release_map['html_url']
        self.assets = release_map['assets']
        self.is_prerelease = release_map['prerelease']

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
        _allow_error = kwds.pop('allow_error', False)
        try:
            value = func(*args, **kwds)
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
            args[-1] += f"[{_default}]:{os.linesep}"
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
        fname = os.path.basename(COSMOVISOR_SERVICE_FILE)
        self.exec(f"wget -c {COSMOVISOR_SERVICE_FILE}")
        with open(fname) as f:
            s = re.sub(
                r'({CHEQD_ROOT_DIR}|{CHEQD_BINARY_NAME})',
                lambda m:{'{CHEQD_ROOT_DIR}': self.cheqd_root_dir,
                        '{CHEQD_BINARY_NAME}': DEFAULT_BINARY_NAME}[m.group()],
                f.read()
            )
        self.remove_safe(fname)
        return s 

    @property
    def logrotate_cfg(self):
        fname = os.path.basename(LOGROTATE_TEMPLATE)
        self.exec(f"wget -c {LOGROTATE_TEMPLATE}")
        with open(fname) as f:
            s = re.sub(
                r'({CHEQD_LOG_DIR})',
                lambda m:{'{CHEQD_LOG_DIR}': self.cheqd_log_dir}[m.group()],
                f.read()
            )
        self.remove_safe(fname)
        return s

    @property
    def rsyslog_cfg(self):
        binary_name = "cosmovisor" if self.interviewer.is_cosmo_needed else DEFAULT_BINARY_NAME
        fname = os.path.basename(RSYSLOG_TEMPLATE)
        self.exec(f"wget -c {RSYSLOG_TEMPLATE}")
        with open(fname) as f:
            s =re.sub(
                r'({BINARY_FOR_LOGGING}|{CHEQD_LOG_DIR})',
                lambda m:{'{BINARY_FOR_LOGGING}': binary_name,
                        '{CHEQD_LOG_DIR}': self.cheqd_log_dir}[m.group()],
                f.read()
            )
        self.remove_safe(fname)
        return s

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
    def exec(self, cmd, use_stdout=True, suppress_err=False):
        self.log(f"Executing command: {cmd}")
        kwargs = {
            "shell": True,
            "check": True,
        }
        if use_stdout:
            kwargs["stdout"] = subprocess.PIPE
        else:
            kwargs["capture_output"] = True

        if suppress_err:
            kwargs["stderr"] = subprocess.DEVNULL
        return subprocess.run(cmd, **kwargs)

    def get_binary(self):
        self.log("Download the binary")
        tar_url = self.release.get_tar_gz_url()
        fname= os.path.basename(tar_url)
        self.exec(f"wget -c {tar_url}")
        self.exec(f"tar xzf {fname}")
        self.remove_safe(fname)

    def is_user_exists(self, username) -> bool:
        try:
            pwd.getpwnam(username)
        except KeyError:
            self.log(f"User {username} does not exist")
            return False
        self.log(f"User {username} already exists")
        return True

    def remove_safe(self, path, is_dir=False):
        if is_dir and os.path.exists(path):
            shutil.rmtree(path)
        if os.path.exists(path):
            os.remove(path)


    def pre_install(self):

        if self.interviewer.is_from_scratch:
            self.log("Removing user's data and configs")
            self.remove_safe(self.cheqd_root_dir, is_dir=True)

            self.remove_safe(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH)
            self.remove_safe(DEFAULT_STANDALONE_SERVICE_FILE_PATH)
            self.remove_safe(DEFAULT_RSYSLOG_FILE)
            self.remove_safe(DEFAULT_LOGROTATE_FILE)

    def prepare_directory_tree(self):
        """
        Needs only in case of Clean installation

        """
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

    def is_service_file_exists(self) -> bool:
        return os.path.exists(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH) or \
            os.path.exists(DEFAULT_STANDALONE_SERVICE_FILE_PATH)

    def setup_systemctl_services(self):
        self.log("Setup systemctl service config")
        if not self.interviewer.is_upgrade or \
                self.interviewer.rewrite_systemd or \
                not self.is_service_file_exists():
            self.remove_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME)
            self.remove_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME)

            if self.interviewer.is_cosmo_needed:
                with open(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH, mode="w") as fd:
                    fd.write(self.cosmovisor_service_cfg)
            else:
                self.exec(f"curl -s {STANDALONE_SERVICE_FILE} > {DEFAULT_STANDALONE_SERVICE_FILE_PATH}")

            self.log("Enable systemctl service")
            self.exec(f"systemctl enable {DEFAULT_COSMOVISOR_SERVICE_NAME if self.interviewer.is_cosmo_needed else DEFAULT_STANDALONE_SERVICE_NAME}")

    def check_systemd_service_on(self, service_name) -> bool:
        stat = self.exec(f'systemctl is-active {service_name}', suppress_err=True, allow_error=True).returncode
        if stat != 0:
            stat = self.exec(f'systemctl is-enabled {service_name}', suppress_err=True, allow_error=True).returncode
        return stat == 0

    def remove_systemd_service(self, service_name):
        if self.check_systemd_service_on(service_name):
            self.log(f"Stopping systemd service: {service_name}")
            self.exec(f"systemctl stop {service_name}")

            self.log(f"Disable systemd service: {service_name}")
            self.exec(f"systemctl disable {service_name}")

            self.log("Reset failed services")
            self.exec("systemctl reset-failed")

            self.log(f"Remove service file")
            self.remove_safe(DEFAULT_STANDALONE_SERVICE_FILE_PATH if service_name == DEFAULT_STANDALONE_SERVICE_NAME else DEFAULT_COSMOVISOR_SERVICE_FILE_PATH)

            self.log("Daemon-reload")
            self.exec('systemctl daemon-reload')

    def setup_system_configs(self):
        if os.path.exists("/etc/rsyslog.d/"):
            if not os.path.exists(DEFAULT_RSYSLOG_FILE) or self.interviewer.rewrite_rsyslog:
                self.log("Configure rsyslog")
                with open(DEFAULT_RSYSLOG_FILE, mode="w") as fd:
                    fd.write(self.rsyslog_cfg)
                # Sometimes it can take a lot of time: https://github.com/rsyslog/rsyslog/issues/3133
                self.exec("systemctl restart rsyslog")

        if os.path.exists("/etc/logrotate.d"):
            if not os.path.exists(DEFAULT_LOGROTATE_FILE) or self.interviewer.rewrite_logrotate:
                self.log("Add config for logrotation")
                with open(DEFAULT_LOGROTATE_FILE, mode="w") as fd:
                    fd.write(self.logrotate_cfg)
                # Sometimes it can take a long period of time: https://github.com/rsyslog/rsyslog/issues/3133
                self.exec("systemctl restart rsyslog")

        self.log("Restart logrotate services")
        self.exec("systemctl restart logrotate.service")
        self.exec("systemctl restart logrotate.timer")

        self.setup_systemctl_services()


    def install(self):

        """
        Steps:
        - remove all data and configs if needed
        - download cheqd-noded binary
        - prepare cheqd user
        - prepare directory tree
        - setup systemctl configs
        - setup cosmovisor if needed
        - copy cheqd-noded binary to cosmovisor or /usr/bin
        - postinstall if needed
        - snapshot if needed
        """
        self.pre_install()

        self.get_binary()

        self.prepare_cheqd_user()

        self.prepare_directory_tree()

        self.setup_system_configs()

        if self.interviewer.is_cosmo_needed:
            self.log("Setup the cosmovisor")
            self.setup_cosmovisor()

        if not self.interviewer.is_cosmo_needed:
            self.log(f"Moving binary from {self.binary_path} to {DEFAULT_INSTALL_PATH}")
            self.exec("sudo mv {} {}".format(self.binary_path, DEFAULT_INSTALL_PATH))

        if self.interviewer.is_setup_needed:
            self.post_install()

        if self.interviewer.init_from_snapshot:
            self.log("Going to download the archive and untar it. It can take a really LONG TIME")
            self.untar_from_snapshot()

    def post_install(self):
        # Init the node with provided moniker
        if not os.path.exists(os.path.join(self.cheqd_config_dir, 'genesis.json')):
            self.exec(f"sudo su -c 'cheqd-noded init {self.interviewer.moniker}' {DEFAULT_CHEQD_USER}")

            # Downloading genesis file
            self.exec(f"curl -s {GENESIS_FILE.format(self.interviewer.chain)} > {os.path.join(self.cheqd_config_dir, 'genesis.json')}")
            shutil.chown(os.path.join(self.cheqd_config_dir, 'genesis.json'),
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)

        # Setting up the external_address
        if self.interviewer.external_address:
            self.exec(f"sudo su -c 'cheqd-noded configure p2p external-address {self.interviewer.external_address}:{self.interviewer.p2p_port}' {DEFAULT_CHEQD_USER}")

        # Setting up the seeds
        seeds = self.exec(f"curl -s {SEEDS_FILE.format(self.interviewer.chain)}").stdout.decode("utf-8").strip()
        self.exec(f"sudo su -c 'cheqd-noded configure p2p seeds {seeds}' {DEFAULT_CHEQD_USER}")

        # Setting up the RPC port
        self.exec(f"sudo su -c 'cheqd-noded configure rpc-laddr \"tcp://0.0.0.0:{self.interviewer.rpc_port}\"' {DEFAULT_CHEQD_USER}")

        # Setting up the P2P port
        self.exec(f"sudo su -c 'cheqd-noded configure p2p laddr \"tcp://0.0.0.0:{self.interviewer.p2p_port}\"' {DEFAULT_CHEQD_USER}")

        # Setting up min gas-price
        self.exec(f"sudo su -c 'cheqd-noded configure min-gas-prices {self.interviewer.gas_price}' {DEFAULT_CHEQD_USER}")

    def prepare_cheqd_user(self):
        if not self.is_user_exists(DEFAULT_CHEQD_USER):
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
        fname= os.path.basename(COSMOVISOR_BINARY_URL)
        self.exec(f"wget -c {COSMOVISOR_BINARY_URL}")
        self.exec(f"tar xzf {fname}")
        self.remove_safe(fname)
        # Remove cosmovisor artifacts...
        self.remove_safe("CHANGELOG.md")
        self.remove_safe("README.md")
        self.mkdir_p(self.cosmovisor_root_dir)
        self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "genesis"))
        self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "genesis/bin"))
        self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "upgrades"))
        if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, "cosmovisor")):
            self.log(f"Moving cosmovisor binary to default installation directory")
            shutil.move("./cosmovisor", DEFAULT_INSTALL_PATH)

        if not os.path.exists(os.path.join(self.cosmovisor_root_dir, "current")):
            self.log(f"Making symlink current -> genesis")
            os.symlink(os.path.join(self.cosmovisor_root_dir, "genesis"),
                       os.path.join(self.cosmovisor_root_dir, "current"))

        self.log(f"Moving binary from {self.binary_path} to {self.cosmovisor_cheqd_bin_path}")
        self.exec("sudo mv {} {}".format(self.binary_path, self.cosmovisor_cheqd_bin_path))

        if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME)):
            self.log(f"Making symlink to {self.cosmovisor_cheqd_bin_path}")
            os.symlink(self.cosmovisor_cheqd_bin_path,
                       os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))

        self.log(f"Changing owner to {DEFAULT_CHEQD_USER} user")
        self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cosmovisor_root_dir}")

    def untar_from_snapshot(self):
        archive_name = os.path.basename(self.interviewer.snapshot_url)

        self.mkdir_p(self.cheqd_data_dir)
        self.log("Install additional tool for showing the progress")
        self.exec("apt install pv")
        self.exec(f"wget -c {self.interviewer.snapshot_url}")
        self.exec(f"sudo su -c 'pv {archive_name} | tar xzf - -C {os.path.join(self.cheqd_root_dir, 'data')}'")
        self.exec(f"rm {archive_name}")

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
        self._is_upgrade = False
        self._is_cosmo_needed = True
        self._init_from_snapshot = False
        self._release = None
        self._chain = chain
        self.verbose = True
        self._snapshot_url = self.prepare_url_for_latest()
        self._is_setup_needed = False
        self._moniker = ""
        self._external_address = ""
        self._rpc_port = DEFAULT_RPC_PORT
        self._p2p_port = DEFAULT_P2P_PORT
        self._gas_price = DEFAULT_GAS_PRICE
        self._is_from_scratch = False
        self._rewrite_systemd = False
        self._rewrite_rsyslog = False
        self._rewrite_logrotate = False

    @property
    def cheqd_root_dir(self):
        return os.path.join(self.home_dir, ".cheqdnode")

    
    @property
    def cheqd_config_dir(self):
        return os.path.join(self.cheqd_root_dir, "config")

    @property
    def cheqd_data_dir(self):
        return os.path.join(self.cheqd_root_dir, "data")

    @property
    def release(self) -> Release:
        return self._release

    @property
    def home_dir(self) -> str:
        return self._home_dir

    @property
    def is_upgrade(self) -> bool:
        return self._is_upgrade

    @property
    def is_from_scratch(self) -> bool:
        return self._is_from_scratch

    @property
    def rewrite_systemd(self) -> bool:
        return self._rewrite_systemd

    @property
    def rewrite_rsyslog(self) -> bool:
        return self._rewrite_rsyslog

    @property
    def rewrite_logrotate(self) -> bool:
        return self._rewrite_logrotate

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

    @property
    def rpc_port(self) -> str:
        return self._rpc_port

    @property
    def p2p_port(self) -> str:
        return self._p2p_port

    @property
    def gas_price(self) -> str:
        return self._gas_price

    @release.setter
    def release(self, release):
        self._release = release

    @home_dir.setter
    def home_dir(self, hd):
        self._home_dir = hd

    @is_upgrade.setter
    def is_upgrade(self, iu):
        self._is_upgrade = iu

    @is_from_scratch.setter
    def is_from_scratch(self, ifs):
        self._is_from_scratch = ifs

    @rewrite_systemd.setter
    def rewrite_systemd(self, rs):
        self._rewrite_systemd = rs

    @rewrite_rsyslog.setter
    def rewrite_rsyslog(self, rr):
        self._rewrite_rsyslog = rr

    @rewrite_logrotate.setter
    def rewrite_logrotate(self, rl):
        self._rewrite_logrotate = rl

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

    @rpc_port.setter
    def rpc_port(self, rpc_port):
        self._rpc_port = rpc_port

    @p2p_port.setter
    def p2p_port(self, p2p_port):
        self._p2p_port = p2p_port

    @gas_price.setter
    def gas_price(self, gas_price):
        self._gas_price = gas_price

    def log(self, msg):
        if self.verbose:
            print(f"{PRINT_PREFIX} {msg}")

    def is_already_installed(self) -> bool:
        return os.path.exists(self.home_dir) and \
            os.path.exists(self.cheqd_root_dir)

    def is_systemd_config_exists(self) -> bool:
        return os.path.exists(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH) or \
            os.path.exists(DEFAULT_STANDALONE_SERVICE_FILE_PATH)


    def get_releases(self):
        req = request.Request("https://api.github.com/repos/cheqd/cheqd-node/releases")
        req.add_header("Accept", "application/vnd.github.v3+json")

        with request.urlopen(req) as response:
            r_list = json.loads(response.read().decode("utf-8"))
            return [Release(r) for r in r_list]

    def get_last_prerelease(self, r_list) -> Release:
        for r in r_list:
            if not r.is_prerelease:
                return r

    def ask_for_version(self):
        all_releases = self.get_releases()
        default = self.get_last_prerelease(all_releases)
        self.log(f"Default version is: {default}")
        d_index = all_releases.index(default)
        last_n_releases = all_releases[d_index:LAST_N_RELEASES]
        print(f"Which version below do you want to install?")
        for i, release in enumerate(last_n_releases):
            print(f"{i + 1}) {release.version}")
        num = self.ask("Please insert the number for picking up the version ",
                       default=1)
        self.release = last_n_releases[num - 1]

    @default_answer
    def ask(self, question, **kwargs):
        return str(input(question))

    def ask_for_home_directory(self, default) -> str:
        answer = self.ask(
            f"Please, type here the path to home directory for user cheqd. For keeping default value, just type "
            f"'Enter'", default=default)
        self.home_dir = answer

    def ask_for_upgrade(self):
        answer = self.ask(
            f"Looks like installation already exists. Do you want to upgrade it? yes/no ", default="No")
        self.is_upgrade = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_install_from_scratch(self):
        answer = self.ask(
            f"Do you want to install from scratch (clean installation)? "
            f"In case of yes it will remove all your configs and data. "
            f"Please make sure that you copied all your configs and private keys. "
            f"Typing no means exit. yes/no ", default="No")
        self.is_from_scratch = True if answer.lower() in ['yes', 'y'] else failure_exit("Aborting...")

    def ask_for_rewrite_systemd(self):
        answer = self.ask(
            f"Do you want to rewrite current systemd configs? yes/no ", default="No")
        self.rewrite_systemd = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_rewrite_logrotate(self):
        answer = self.ask(
            f"Do you want to rewrite current logrotate config? yes/no ", default="No")
        self.rewrite_logrotate = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_rewrite_rsyslog(self):
        answer = self.ask(
            f"Do you want to rewrite current rsyslog config? yes/no ", default="No")
        self.rewrite_rsyslog = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_cosmovisor(self, text, default) -> str:
        answer = self.ask(text, default=default)
        self.is_cosmo_needed = True if answer.lower() in ['yes', 'y'] else False

    def ask_for_init_from_snapshot(self):
        answer = self.ask(
            f"Do you want to deploy the latest snapshot? "
            f"Please type any kind of variants: yes/no. ",
            default="No"
        )
        self.init_from_snapshot = True if answer.lower() in ['yes', 'y'] else False
        if self.init_from_snapshot:
            self.snapshot_url = self.prepare_url_for_latest()

    def ask_for_chain(self):
        answer = self.ask(
            f"Which chain do you want to use? Possible variants are: {', '.join(DEFAULT_CHAINS)} ",
            default=DEFAULT_CHAIN
        )
        self.chain = answer if answer in DEFAULT_CHAINS else failure_exit(f"Possible chains are: {DEFAULT_CHAINS}")

    def ask_for_setup(self):
        answer = self.ask(
            f"Do you want to setup node after installation? "
            f"Please type any kind of variants: yes/no. ",
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
            f"What are the external IP address for your node? {os.linesep}",
            default=""
        )
        self.external_address = answer

    def ask_for_rpc_port(self):
        answer = self.ask(
            f"What is the RPC port? ",
            default=DEFAULT_RPC_PORT
        )
        self.rpc_port = answer

    def ask_for_p2p_port(self):
        answer = self.ask(
            f"What is the P2P port? ",
            default=DEFAULT_P2P_PORT
        )
        self.p2p_port = answer

    def ask_for_gas_price(self):
        answer = self.ask(
            f"What is the gas-price? ",
            default=DEFAULT_GAS_PRICE
        )
        self.gas_price = answer

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
    if interviewer.is_already_installed():
        interviewer.ask_for_upgrade()
        if interviewer.is_upgrade:
            if os.path.exists(DEFAULT_LOGROTATE_FILE):
                interviewer.ask_for_rewrite_logrotate()
            if os.path.exists(DEFAULT_RSYSLOG_FILE):
                interviewer.ask_for_rewrite_rsyslog()
            if interviewer.is_systemd_config_exists():
                interviewer.ask_for_rewrite_systemd()
            interviewer.ask_for_cosmovisor(f"Do you use Cosmovisor now? Please type any kind of variants: yes/no ", default=DEFAULT_USE_COSMOVISOR)
        elif not interviewer.is_upgrade:
            interviewer.ask_for_install_from_scratch()
    else:
        interviewer.ask_for_cosmovisor(f"Do you want to use Cosmovisor? Please type any kind of variants: yes/no ", default=DEFAULT_USE_COSMOVISOR)
    if not interviewer.is_upgrade:
        interviewer.ask_for_chain()
        interviewer.ask_for_init_from_snapshot()
    if not interviewer.is_upgrade:
        interviewer.ask_for_setup()
    if interviewer.is_setup_needed:
        interviewer.ask_for_moniker()
        interviewer.ask_for_external_address()
        interviewer.ask_for_rpc_port()
        interviewer.ask_for_p2p_port()
        interviewer.ask_for_gas_price()

    # Install
    installer = Installer(interviewer)
    installer.install()
