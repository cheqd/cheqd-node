#!/usr/bin/env python3


###############################################################
###     		    Python package imports      			###
###############################################################
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
import platform
import copy

###############################################################
###     				Installer defaults    				###
###############################################################
LAST_N_RELEASES = 5
DEFAULT_HOME = "/home/cheqd"
DEFAULT_INSTALL_PATH = "/usr/bin"
DEFAULT_CHEQD_USER = "cheqd"
DEFAULT_BINARY_NAME = "cheqd-noded"
DEFAULT_COSMOVISOR_BINARY_NAME = "cosmovisor"
DEFAULT_CHAINS = ['testnet', 'mainnet']
DEFAULT_CHAIN = "mainnet"
PRINT_PREFIX = "********* "

###############################################################
###     				Cosmovisor Config      				###
###############################################################
COSMOVISOR_BINARY_URL = "https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz"
DEFAULT_USE_COSMOVISOR = "yes"

###############################################################
###     				Systemd Config      				###
###############################################################
STANDALONE_SERVICE_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/build-tools/node-standalone.service"
COSMOVISOR_SERVICE_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/build-tools/node-cosmovisor.service"
LOGROTATE_TEMPLATE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/build-tools/logrotate.conf"
RSYSLOG_TEMPLATE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/build-tools/rsyslog.conf"
DEFAULT_STANDALONE_SERVICE_NAME = 'cheqd-noded'
DEFAULT_COSMOVISOR_SERVICE_NAME = 'cheqd-cosmovisor'
DEFAULT_STANDALONE_SERVICE_FILE_PATH = f"/lib/systemd/system/{DEFAULT_STANDALONE_SERVICE_NAME}.service"
DEFAULT_COSMOVISOR_SERVICE_FILE_PATH = f"/lib/systemd/system/{DEFAULT_COSMOVISOR_SERVICE_NAME}.service"
DEFAULT_LOGROTATE_FILE = "/etc/logrotate.d/cheqd-node"
DEFAULT_RSYSLOG_FILE = "/etc/rsyslog.d/cheqd-node.conf"

###############################################################
###     		Network configuration files    				###
###############################################################
GENESIS_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/networks/{}/genesis.json"
SEEDS_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/main/networks/{}/seeds.txt"

###############################################################
###     				Node snapshots      				###
###############################################################
DEFAULT_SNAPSHOT_SERVER = "https://snapshots.cheqd.net"
DEFAULT_INIT_FROM_SNAPSHOT = "yes"
TESTNET_SNAPSHOT = "https://cheqd-node-backups.ams3.cdn.digitaloceanspaces.com/testnet/latest/cheqd-testnet-4_{}.tar.gz"
MAINNET_SNAPSHOT = "https://cheqd-node-backups.ams3.cdn.digitaloceanspaces.com/mainnet/latest/cheqd-mainnet-1_{}.tar.gz"
CHECKSUM_URL_BASE = "https://cheqd-node-backups.ams3.cdn.digitaloceanspaces.com/"

###############################################################
###     	    Default node configuration      			###
###############################################################
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

    def get_release_url(self):
        try:
            release_urls = [ a['browser_download_url'] for a in self.assets if (a['browser_download_url'].find("cheqd-node") > 0 and a['browser_download_url'].find(".deb")) == -1]
            if len(release_urls) > 0:
                return release_urls[0]
            else:
                failure_exit(f"No asset found to download for release: {self.version}")
        except:
            failure_exit(f"Failed to get cheqd-node binaries from Github")

    def __str__(self):
        return f"Name: {self.version}"


def failure_exit(reason):
    print(f"Reason for failure: {reason}")
    print("Exiting...")
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
            args[-1] += f" [default: {_default}]:{os.linesep}"
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
        binary_name = DEFAULT_COSMOVISOR_BINARY_NAME if self.interviewer.is_cosmo_needed else DEFAULT_BINARY_NAME
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
        self.log("Downloading cheqd-noded binary...")
        binary_url = self.release.get_release_url()
        fname = os.path.basename(binary_url)
        try:
            self.exec(f"wget -c {binary_url}")
            if fname.find(".tar.gz") != -1:
                self.exec(f"tar -xzf {fname}")
                self.remove_safe(fname)
            self.exec(f"chmod +x {DEFAULT_BINARY_NAME}")
        except:
            failure_exit("Failed to download binary")

    def is_user_exists(self, username) -> bool:
        try:
            pwd.getpwnam(username)
            self.log(f"User {username} already exists")
            return True
        except KeyError:
            self.log(f"User {username} does not exist")
            return False
        

    def remove_safe(self, path, is_dir=False):
        if is_dir and os.path.exists(path):
            shutil.rmtree(path)
        if os.path.exists(path):
            os.remove(path)


    def pre_install(self):
        if self.interviewer.is_from_scratch:
            self.log("Removing user's data and configs")
            self.remove_safe(self.cheqd_root_dir, is_dir=True)
            self.remove_safe(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))
            self.remove_safe(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_COSMOVISOR_BINARY_NAME))
            self.remove_safe(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH)
            self.remove_safe(DEFAULT_STANDALONE_SERVICE_FILE_PATH)
            self.remove_safe(DEFAULT_RSYSLOG_FILE)
            self.remove_safe(DEFAULT_LOGROTATE_FILE)

    def prepare_directory_tree(self):
        """
        Needed only in case of clean installation

        """
        try:
            if not os.path.exists(self.cheqd_root_dir):
                self.log("Creating main directory for cheqd-noded")
                self.mkdir_p(self.cheqd_root_dir)

                self.log(f"Setting directory permissions to default cheqd user: {DEFAULT_CHEQD_USER}")
                self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.interviewer.home_dir}")
            else:
                self.log(f"Skipping main directory creation because {self.cheqd_root_dir} already exists")

            if not os.path.exists(self.cheqd_log_dir):
                self.log("Creating log directory for cheqd-noded")
                self.setup_log_dir()
            else:
                self.log(f"Skipping log directory creation because {self.cheqd_log_dir} already exists")

            if not os.path.exists("/var/log/cheqd-node"):
                self.log("Creating a symlink from cheqd-noded log folder to /var/log/cheqd-node")
                os.symlink(self.cheqd_log_dir, "/var/log/cheqd-node", target_is_directory=True)
            else:
                self.log("Skipping linking because /var/log/cheqd-node already exists")
        except:
            failure_exit("Failed to prepare directory tree for {DEFAULT_CHEQD_USER}")

    def is_service_file_exists(self) -> bool:
        return os.path.exists(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH) or \
            os.path.exists(DEFAULT_STANDALONE_SERVICE_FILE_PATH)

    def setup_systemctl_services(self):
        self.log("Setting up systemd config")
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

            self.log("Enabling systemd service for cheqd-noded")
            self.exec(f"systemctl enable {DEFAULT_COSMOVISOR_SERVICE_NAME if self.interviewer.is_cosmo_needed else DEFAULT_STANDALONE_SERVICE_NAME}")

    def check_systemd_service_on(self, service_name) -> bool:
        # pylint: disable=E1123
        stat = self.exec(f'systemctl is-active {service_name}', suppress_err=True, allow_error=True).returncode
        if stat != 0:
            # pylint: disable=E1123
            stat = self.exec(f'systemctl is-enabled {service_name}', suppress_err=True, allow_error=True).returncode
        return stat == 0

    def remove_systemd_service(self, service_name):
        if self.check_systemd_service_on(service_name):
            self.log(f"Stopping systemd service: {service_name}")
            self.exec(f"systemctl stop {service_name}")

            self.log(f"Disable systemd service: {service_name}")
            self.exec(f"systemctl disable {service_name}")

            self.log("Reset failed systemd services (if any)")
            self.exec("systemctl reset-failed")

            self.log("Reload systemd config")
            self.exec('systemctl daemon-reload')

    def setup_system_configs(self):
        if os.path.exists("/etc/rsyslog.d/"):
            if not os.path.exists(DEFAULT_RSYSLOG_FILE) or self.interviewer.rewrite_rsyslog:
                self.log("Configuring syslog systemd service for cheqd-noded logging")
                with open(DEFAULT_RSYSLOG_FILE, mode="w") as fd:
                    fd.write(self.rsyslog_cfg)
                # Sometimes it can take a lot of time: https://github.com/rsyslog/rsyslog/issues/3133
                self.log("Restarting rsyslog service")
                self.exec("systemctl restart rsyslog")

        if os.path.exists("/etc/logrotate.d"):
            if not os.path.exists(DEFAULT_LOGROTATE_FILE) or self.interviewer.rewrite_logrotate:
                self.log("Configuring log rotation systemd service for cheqd-noded logging")
                with open(DEFAULT_LOGROTATE_FILE, mode="w") as fd:
                    fd.write(self.logrotate_cfg)

        self.log("Restarting logrotate services")
        self.exec("systemctl restart logrotate.service")
        self.exec("systemctl restart logrotate.timer")

        self.setup_systemctl_services()


    def install(self):
        """
        Steps:
        - Remove all data and configurations (if needed)
        - Download cheqd-noded binary
        - Prepare cheqd user
        - Prepare directory tree
        - Setup systemctl configs
        - Setup Cosmovisor (if selected by user)
        - Install cheqd-noded binary at system bin or Cosmovisor bin path
        - Carry out post-install actions
        - Restore and download snapshot (if selected by user)
        """

        self.pre_install()
        self.get_binary()
        self.prepare_cheqd_user()
        self.prepare_directory_tree()
        self.setup_system_configs()

        if self.interviewer.is_cosmo_needed:
            self.log("Setting up Cosmovisor")
            self.setup_cosmovisor()

        if not self.interviewer.is_cosmo_needed:
            self.log(f"Moving binary from {self.binary_path} to {DEFAULT_INSTALL_PATH}")
            self.exec("sudo mv {} {}".format(self.binary_path, DEFAULT_INSTALL_PATH))

        if self.interviewer.is_setup_needed:
            self.post_install()

        if self.interviewer.init_from_snapshot:
            self.log("Downloading snapshot and extracting archive. This can take a *really* long time...")
            self.download_snapshot()
            self.untar_from_snapshot()

    def post_install(self):
        # Init the node with provided moniker
        if not os.path.exists(os.path.join(self.cheqd_config_dir, 'genesis.json')):
            self.exec(f"sudo su -c 'cheqd-noded init {self.interviewer.moniker}' {DEFAULT_CHEQD_USER}")

            # Downloading genesis file
            self.exec(f"curl {GENESIS_FILE.format(self.interviewer.chain)} > {os.path.join(self.cheqd_config_dir, 'genesis.json')}")
            shutil.chown(os.path.join(self.cheqd_config_dir, 'genesis.json'),
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)

        # Setting up the external_address
        if self.interviewer.external_address:
            self.exec(f"sudo su -c 'cheqd-noded configure p2p external-address {self.interviewer.external_address}:{self.interviewer.p2p_port}' {DEFAULT_CHEQD_USER}")

        # Setting up the seeds
        seeds = self.exec(f"curl {SEEDS_FILE.format(self.interviewer.chain)}").stdout.decode("utf-8").strip()
        self.exec(f"sudo su -c 'cheqd-noded configure p2p seeds {seeds}' {DEFAULT_CHEQD_USER}")

        # Setting up the RPC port
        self.exec(f"sudo su -c 'cheqd-noded configure rpc-laddr \"tcp://0.0.0.0:{self.interviewer.rpc_port}\"' {DEFAULT_CHEQD_USER}")

        # Setting up the P2P port
        self.exec(f"sudo su -c 'cheqd-noded configure p2p laddr \"tcp://0.0.0.0:{self.interviewer.p2p_port}\"' {DEFAULT_CHEQD_USER}")

        # Setting up min gas-price
        self.exec(f"sudo su -c 'cheqd-noded configure min-gas-prices {self.interviewer.gas_price}' {DEFAULT_CHEQD_USER}")

    def prepare_cheqd_user(self):
        try:
            if not self.is_user_exists(DEFAULT_CHEQD_USER):
                self.log(f"Creating {DEFAULT_CHEQD_USER} group")
                self.exec(f"addgroup {DEFAULT_CHEQD_USER} --quiet --system")
                self.log(f"Creating {DEFAULT_CHEQD_USER} user and adding to {DEFAULT_CHEQD_USER} group")
                self.exec(
                    f"adduser --system {DEFAULT_CHEQD_USER} --home {self.interviewer.home_dir} --shell /bin/bash --ingroup {DEFAULT_CHEQD_USER} --quiet")
        except:
            failure_exit(f"Failed to create {DEFAULT_CHEQD_USER} user")

    def mkdir_p(self, dir_name):
        try:
            os.mkdir(dir_name)
        except FileExistsError as err:
            self.log(f"Directory {dir_name} already exists")

    def setup_log_dir(self):
        try:
            self.mkdir_p(self.cheqd_log_dir)
            Path(os.path.join(self.cheqd_log_dir, "stdout.log")).touch()
            self.log(f"Setting up ownership permissions for {self.cheqd_log_dir} directory")
            self.exec(f"chown -R syslog:cheqd {self.cheqd_log_dir}")
        except:
            failure_exit(f"Failed to setup {self.cheqd_log_dir} directory")

    def setup_cosmovisor(self):
        try:
            fname= os.path.basename(COSMOVISOR_BINARY_URL)
            self.exec(f"wget -c {COSMOVISOR_BINARY_URL}")
            self.exec(f"tar -xzf {fname}")
            self.remove_safe(fname)
            
            # Remove cosmovisor artifacts...
            self.remove_safe("CHANGELOG.md")
            self.remove_safe("README.md")
            self.mkdir_p(self.cosmovisor_root_dir)
            self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "genesis"))
            self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "genesis/bin"))
            self.mkdir_p(os.path.join(self.cosmovisor_root_dir, "upgrades"))
            if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_COSMOVISOR_BINARY_NAME)):
                self.log(f"Moving Cosmovisor binary to installation directory")
                shutil.move("./cosmovisor", DEFAULT_INSTALL_PATH)

            if not os.path.exists(os.path.join(self.cosmovisor_root_dir, "current")):
                self.log(f"Creating symlink for current Cosmovisor version")
                os.symlink(os.path.join(self.cosmovisor_root_dir, "genesis"),
                        os.path.join(self.cosmovisor_root_dir, "current"))

            self.log(f"Moving binary from {self.binary_path} to {self.cosmovisor_cheqd_bin_path}")
            self.exec("sudo mv {} {}".format(self.binary_path, self.cosmovisor_cheqd_bin_path))

            if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME)):
                self.log(f"Creating symlink to {self.cosmovisor_cheqd_bin_path}")
                os.symlink(self.cosmovisor_cheqd_bin_path,
                        os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))

            if self.interviewer.is_upgrade and \
                os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")):

                self.log(f"Copying upgrade-info.json file to cosmovisor/current/")
                shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                            os.path.join(self.cosmovisor_root_dir, "current"))
                self.log(f"Changing owner to {DEFAULT_CHEQD_USER} user")
                self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cosmovisor_root_dir}")
        
            self.log(f"Changing directory ownership for Cosmovisor to {DEFAULT_CHEQD_USER} user")
            self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cosmovisor_root_dir}")
        except:
            failure_exit(f"Failed to setup Cosmovisor")

    def compare_checksum(self, file_path):
        # Set URL for correct checksum file for snapshot
        checksum_url = os.path.join(CHECKSUM_URL_BASE, self.interviewer.chain, "latest/md5sum.txt")
        # Get checksum file
        published_checksum = self.exec(f"curl -s {checksum_url} | tail -1 | cut -d' ' -f 1").stdout.strip()
        self.log(f"Comparing published checksum with local checksum")
        local_checksum = self.exec(f"md5sum {file_path} | tail -1 | cut -d' ' -f 1").stdout.strip()
        if published_checksum == local_checksum:
            self.log(f"Checksums match. Download is OK.")
            return True
        elif published_checksum != local_checksum:
            self.log(f"Checksums do not match. Download got corrupted.")
            return False
        else:
            failure_exit(f"Error encountered when comparing checksums.")

    def install_dependencies(self):
        try:
            self.log("Installing dependencies")
            self.exec("sudo apt-get update")
            self.log(f"Install pv to show progress of extraction")
            self.exec("sudo apt-get install -y pv")
        except:
            failure_exit(f"Failed to install dependencies")

    def download_snapshot(self):
        try:
            archive_name = os.path.basename(self.interviewer.snapshot_url)
            self.mkdir_p(self.cheqd_data_dir)
            self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cheqd_data_dir}")
            # Fetch size of snapshot archive. Uses curl to fetch headers and looks for Content-Length.
            archive_size = self.exec(f"curl -s --head {self.interviewer.snapshot_url} | awk '/Length/ {{print $2}}'").stdout.strip()
            # Check how much free disk space is available wherever the cheqd root directory is mounted
            free_disk_space = self.exec(f"df -P -B1 {self.cheqd_root_dir} | tail -1 | awk '{{print $4}}'").stdout.strip()
            if int(archive_size) < int(free_disk_space):
                self.log(f"Downloading snapshot archive. This may take a while...")
                self.exec(f"wget -c {self.interviewer.snapshot_url} -P {self.cheqd_root_dir}")
                archive_path = os.path.join(self.cheqd_root_dir, archive_name)
                if self.compare_checksum(archive_path) is True:
                    self.log(f"Snapshot download was successful and checksums match.")
                else:
                    self.log(f"Snapshot download was successful but checksums do not match.")
                    failure_exit(f"Snapshot download was successful but checksums do not match.")
            elif int(archive_size) > int(free_disk_space):
                failure_exit (f"Snapshot archive is too large to fit in free disk space. Please free up some space and try again.")
            else:
                failure_exit (f"Error encountered when downloading snapshot archive.")
        except:
            failure_exit(f"Failed to download snapshot")
        
    def untar_from_snapshot(self):
        try:
            archive_path = os.path.join(self.cheqd_root_dir, os.path.basename(self.interviewer.snapshot_url))
            # Check if there is enough space to extract snapshot archive
            self.install_dependencies()
            self.log(f"Extracting snapshot archive. This may take a while...")

            # Extract to cheqd node data directory EXCEPT for validator state
            self.exec(f"sudo su -c 'pv {archive_path} | tar xzf - -C {self.cheqd_data_dir} --exclude priv_validator_state.json' {DEFAULT_CHEQD_USER}")
            
            # Delete snapshot archive file
            self.log(f"Snapshot extraction was successful. Deleting snapshot archive.")
            self.remove_safe(archive_path)
            # Workaround to make this work with Cosmovisor since it expects upgrade-info.json file in cosmovisor/current directory
            if self.interviewer.is_cosmo_needed:
                if os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")):
                    self.log(f"Copying upgrade-info.json file to cosmovisor/current/")
                    shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                                os.path.join(self.cosmovisor_root_dir, "current"))
                self.log(f"Changing owner to {DEFAULT_CHEQD_USER} user")
                self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cosmovisor_root_dir}")
            self.exec(f"chown -R {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER} {self.cheqd_data_dir}")
        except:
            failure_exit(f"Failed to extract snapshot")
        

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
        if os.path.exists(self.home_dir) and os.path.exists(self.cheqd_root_dir):
            return True
        elif not os.path.exists(self.home_dir) and not os.path.exists(self.cheqd_root_dir):
            return False
        else:
            failure_exit(f"Could not check if cheqd-node is already installed.")

    def is_systemd_config_exists(self) -> bool:
        return os.path.exists(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH) or \
            os.path.exists(DEFAULT_STANDALONE_SERVICE_FILE_PATH)

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

    def get_releases(self):
        req = request.Request("https://api.github.com/repos/cheqd/cheqd-node/releases")
        req.add_header("Accept", "application/vnd.github.v3+json")
        with request.urlopen(req) as response:
            r_list = json.loads(response.read().decode("utf-8").strip())
            return [Release(r) for r in r_list]

    def get_latest_release(self):
        req = request.Request("https://api.github.com/repos/cheqd/cheqd-node/releases/latest")
        req.add_header("Accept", "application/vnd.github.v3+json")
        with request.urlopen(req) as response:
            return Release(json.loads(response.read().decode("utf-8")))

    def remove_release_from_list(self, r_list, elem):
        copy_r_list = copy.deepcopy(r_list)
        for i, release in enumerate(r_list):
            if release.version == elem.version:
                copy_r_list.pop(i)
                return copy_r_list

    def ask_for_version(self):
        default = self.get_latest_release()
        all_releases = self.get_releases()
        self.log(f"Latest stable cheqd-noded release version is {default}")
        self.log(f"List of cheqd-noded releases: ")
        all_releases = self.remove_release_from_list(all_releases, default)
        all_releases.insert(0, default)
        for i, release in enumerate(all_releases[0: LAST_N_RELEASES]):
            print(f"{i + 1}) {release.version}")
        release_num = int(self.ask("Choose list option number above to select version of cheqd-node to install", 
            default=1))
        if release_num >= 1 and release_num <= len(all_releases):
            self.release = all_releases[release_num - 1]
        else:
            failure_exit(f"Invalid release number picked from list of releases: {release_num}")

    @default_answer
    def ask(self, question, **kwargs):
        return str(input(question)).strip()

    def ask_for_home_directory(self, default) -> str:
        self.home_dir = self.ask(
            f"Set path for cheqd user's home directory", default=default)

    def ask_for_upgrade(self):
        answer = self.ask(
            f"Existing cheqd-node configuration folder detected. Do you want to upgrade an existing cheqd-node installation? (yes/no)", default="no")
        if answer.lower().startswith("y"):
            self.is_upgrade = True
        elif answer.lower().startswith("n"):
            self.is_upgrade = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_install_from_scratch(self):
        answer = self.ask(
            f"WARNING: Doing a fresh installation of cheqd-node will remove ALL existing configuration and data. "
            f"CAUTION: Please ensure you have a backup of your existing configuration and data before proceeding. "
            f"Do you want to do fresh installation of cheqd-node? (yes/no)", default="no")
        if answer.lower().startswith("y"):
            self.is_from_scratch = True
        elif answer.lower().startswith("n"):
            self.is_from_scratch = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_rewrite_systemd(self):
        answer = self.ask(
            f"Overwrite existing systemd configuration for cheqd-node? (yes/no)", default="yes")
        if answer.lower().startswith("y"):
            self.rewrite_systemd = True
        elif answer.lower().startswith("n"):
            self.rewrite_systemd = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_rewrite_logrotate(self):
        answer = self.ask(
            f"Overwrite existing configuration for logrotate? (yes/no)", default="yes")
        if answer.lower().startswith("y"):
            self.rewrite_logrotate = True
        elif answer.lower().startswith("n"):
            self.rewrite_logrotate = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_rewrite_rsyslog(self):
        answer = self.ask(
            f"Overwrite existing configuration for cheqd-node logging? (yes/no)", default="yes")
        if answer.lower().startswith("y"):
            self.rewrite_rsyslog = True
        elif answer.lower().startswith("n"):
            self.rewrite_rsyslog = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_cosmovisor(self):
        self.log(f"INFO: Installing cheqd-node with Cosmovisor allows for automatic unattended upgrades for valid software upgrade proposals.")
        answer = self.ask(f"Install cheqd-noded using Cosmovisor? (yes/no)", default=DEFAULT_USE_COSMOVISOR)
        if answer.lower().startswith("y"):
            self.is_cosmo_needed = True
        elif answer.lower().startswith("n"):
            self.is_cosmo_needed = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_init_from_snapshot(self):
        answer = self.ask(
            f"CAUTION: Downloading a snapshot replaces your existing copy of chain data. Usually safe to use this option when doing a fresh installation. "
            f"Do you want to download a snapshot of the existing chain to speed up node synchronisation? (yes/no)", default="yes")
        if answer.lower().startswith("y"):
            self.snapshot_url = self.prepare_url_for_latest()
            self.init_from_snapshot = True
        elif answer.lower().startswith("n"):
            self.init_from_snapshot = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_chain(self):
        answer = self.ask(
            f"Select cheqd network to join ({'/'.join(DEFAULT_CHAINS)})", default="mainnet")
        if answer in DEFAULT_CHAINS:
            self.chain = answer
        else:
            failure_exit(f"Invalid network selected during installation.")

    def ask_for_setup(self):
        answer = self.ask(
            f"Do you want to setup a new cheqd-node? (yes/no)", default="yes")
        if answer.lower().startswith("y"):
            self.is_setup_needed = True
        elif answer.lower().startswith("n"):
            self.is_setup_needed = False
        else:
            failure_exit(f"Invalid input provided during installation.")

    def ask_for_moniker(self):
        answer = self.ask(
            f"Provide a moniker for your cheqd-node", default=platform.node())
        if answer is not None:
            self.moniker = answer
        else:
            failure_exit(f"Invalid moniker provided during cheqd-noded setup.")

    def ask_for_external_address(self):
        answer = self.ask(
            f"What is the externally-reachable IP address or DNS name for your cheqd-node? [default: Fetch automatically via DNS resolver lookup]: {os.linesep}")
        if answer is not None:
            self.external_address = answer
        else:
            try:
                self.external_address = self.exec("dig +short txt ch whoami.cloudflare @1.1.1.1").stdout.replace('"', '').strip()
            except:
                failure_exit(f"Unable to fetch external IP address for your node.")

    def ask_for_rpc_port(self):
        self.rpc_port = self.ask(
            f"Specify port for Tendermint RPC", default=DEFAULT_RPC_PORT)

    def ask_for_p2p_port(self):
        self.p2p_port = self.ask(
            f"Specify port for Tendermint P2P", default=DEFAULT_P2P_PORT)

    def ask_for_gas_price(self):
        self.gas_price = self.ask(
            f"Specify minimum gas price for transactions", default=DEFAULT_GAS_PRICE)

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
    
    # Steps to execute if installing from scratch
    def install_steps():
        interviewer.ask_for_setup()
        interviewer.ask_for_chain()
        interviewer.ask_for_cosmovisor()
        interviewer.ask_for_init_from_snapshot()
        if interviewer.is_setup_needed:
            interviewer.ask_for_moniker()
            interviewer.ask_for_external_address()
            interviewer.ask_for_rpc_port()
            interviewer.ask_for_p2p_port()
            interviewer.ask_for_gas_price()

    # Steps to execute if upgrading existing node
    def upgrade_steps():
        interviewer.ask_for_cosmovisor()
        if interviewer.is_systemd_config_exists():
            interviewer.ask_for_rewrite_systemd()
        if os.path.exists(DEFAULT_RSYSLOG_FILE):
            interviewer.ask_for_rewrite_rsyslog()
        if os.path.exists(DEFAULT_LOGROTATE_FILE):
            interviewer.ask_for_rewrite_logrotate()
    
    # Ask user for information
    interviewer = Interviewer()
    interviewer.ask_for_version()
    interviewer.ask_for_home_directory(default=DEFAULT_HOME)

    # Check if cheqd configuration directory exists
    is_installed = interviewer.is_already_installed()
    
    # First-time new node setup
    if is_installed is False:
        install_steps()
    elif is_installed is True:
        # Check if user wants to upgrade existing cheqd-node installation
        interviewer.ask_for_upgrade()
        if interviewer.is_upgrade is True:
            upgrade_steps()
        elif interviewer.is_upgrade is False:
            interviewer.ask_for_install_from_scratch()
            if interviewer.is_from_scratch is True:
                install_steps()
            else:
                failure_exit("Aborting installation to prevent overwriting existing cheqd-node.")
        else:
            failure_exit("Unable to determine upgrade/installation mode.")
    else:
        failure_exit("Could not execute either install or upgrade steps.")

    # Install
    installer = Installer(interviewer)
    try:
        installer.install()
    except:
        failure_exit("Unable to install cheqd-node.")
