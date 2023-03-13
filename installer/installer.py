#!/usr/bin/env python3


###############################################################
###     		    Python package imports      			###
###############################################################
from pathlib import Path
import copy
import datetime
import functools
import json
import logging
import os
import platform
import pwd
import re
import shutil
import signal
import subprocess
import sys
import tarfile
import urllib.request as request

###############################################################
###     				Installer defaults    				###
###############################################################
LAST_N_RELEASES = 5
DEFAULT_CHEQD_HOME_DIR = "/home/cheqd"
DEFAULT_INSTALL_PATH = "/usr/bin"
DEFAULT_CHEQD_USER = "cheqd"
DEFAULT_BINARY_NAME = "cheqd-noded"
DEFAULT_COSMOVISOR_BINARY_NAME = "cosmovisor"
MAINNET_CHAIN_ID = "cheqd-mainnet-1"
TESTNET_CHAIN_ID = "cheqd-testnet-6"
PRINT_PREFIX = "********* "
# Set branch dynamically in CI workflow for testing if Python dev mode is enabled and DEFAULT_DEBUG_BRANCH is set
# Otherwise, use the main branch
DEFAULT_DEBUG_BRANCH = os.getenv("DEFAULT_DEBUG_BRANCH") if os.getenv("DEFAULT_DEBUG_BRANCH") != None else "main"

###############################################################
###     		Cosmovisor configuration      				###
###############################################################
DEFAULT_LATEST_COSMOVISOR_VERSION = "v1.2.0"
COSMOVISOR_BINARY_URL = "https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2F{}/cosmovisor-{}-linux-{}.tar.gz"
DEFAULT_USE_COSMOVISOR = "yes"
DEFAULT_BUMP_COSMOVISOR = "yes"
DEFAULT_DAEMON_ALLOW_DOWNLOAD_BINARIES = "true"
DEFAULT_DAEMON_RESTART_AFTER_UPGRADE = "true"
DEFAULT_DAEMON_POLL_INTERVAL = "300s"
DEFAULT_UNSAFE_SKIP_BACKUP = "true"
DEFAULT_DAEMON_RESTART_DELAY = "120s"


###############################################################
###     			Systemd configuration      				###
###############################################################
STANDALONE_SERVICE_TEMPLATE = f"https://raw.githubusercontent.com/cheqd/cheqd-node/{DEFAULT_DEBUG_BRANCH}/build-tools/cheqd-noded.service"
COSMOVISOR_SERVICE_TEMPLATE = f"https://raw.githubusercontent.com/cheqd/cheqd-node/{DEFAULT_DEBUG_BRANCH}/build-tools/cheqd-cosmovisor.service"
LOGROTATE_TEMPLATE = f"https://raw.githubusercontent.com/cheqd/cheqd-node/{DEFAULT_DEBUG_BRANCH}/build-tools/logrotate.conf"
RSYSLOG_TEMPLATE = f"https://raw.githubusercontent.com/cheqd/cheqd-node/{DEFAULT_DEBUG_BRANCH}/build-tools/rsyslog.conf"
DEFAULT_STANDALONE_SERVICE_NAME = 'cheqd-noded'
DEFAULT_COSMOVISOR_SERVICE_NAME = 'cheqd-cosmovisor'
DEFAULT_STANDALONE_SERVICE_FILE_PATH = f"/lib/systemd/system/{DEFAULT_STANDALONE_SERVICE_NAME}.service"
DEFAULT_COSMOVISOR_SERVICE_FILE_PATH = f"/lib/systemd/system/{DEFAULT_COSMOVISOR_SERVICE_NAME}.service"
DEFAULT_LOGROTATE_FILE = "/etc/logrotate.d/cheqd-node"
DEFAULT_RSYSLOG_FILE = "/etc/rsyslog.d/cheqd-node.conf"


###############################################################
###     		Network configuration files    				###
###############################################################
GENESIS_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/%s/networks/{}/genesis.json" % (
    DEFAULT_DEBUG_BRANCH)
SEEDS_FILE = "https://raw.githubusercontent.com/cheqd/cheqd-node/%s/networks/{}/seeds.txt" % (
    DEFAULT_DEBUG_BRANCH)


###############################################################
###     				Node snapshots      				###
###############################################################
DEFAULT_INIT_FROM_SNAPSHOT = "yes"
TESTNET_SNAPSHOT = "https://snapshots-cdn.cheqd.net/testnet/{}/cheqd-testnet-6_{}.tar.lz4"
MAINNET_SNAPSHOT = "https://snapshots-cdn.cheqd.net/mainnet/{}/cheqd-mainnet-1_{}.tar.lz4"
MAX_SNAPSHOT_DAYS = 7


###############################################################
###     	Default node environment variables      	    ###
###############################################################
DEFAULT_RPC_PORT = "26657"
DEFAULT_P2P_PORT = "26656"
CHEQD_NODED_HOME = "/home/cheqd/.cheqdnode"
CHEQD_NODED_NODE = "tcp://localhost:26657"
CHEQD_NODED_MONIKER = platform.node()
CHEQD_NODED_CHAIN_ID = MAINNET_CHAIN_ID
CHEQD_NODED_MINIMUM_GAS_PRICES = "50ncheq"
CHEQD_NODED_LOG_LEVEL = "error"
CHEQD_NODED_LOG_FORMAT = "json"
CHEQD_NODED_FASTSYNC_VERSION = "v0"
CHEQD_NODED_P2P_MAX_PACKET_MSG_PAYLOAD_SIZE = 10240


###############################################################
###     	    Common, reusable functions    	            ###
###############################################################

# Set logging configuration
if sys.flags.dev_mode:
    logging.basicConfig(format='%(asctime)s %(levelname)s: %(message)s',
                    datefmt='%d-%b-%Y %H:%M:%S', 
                    level=logging.DEBUG)
    logging.raiseExceptions = True
    logging.propagate = True
else:
    logging.basicConfig(format='%(asctime)s %(levelname)s: %(message)s',
                        datefmt='%d-%b-%Y %H:%M:%S', 
                        level=logging.INFO)
    logging.raiseExceptions = True
    logging.propagate = True

# Handle Ctrl+C / SIGINT halts requests
def sigint_handler(signal, frame):
    logging.info(f'Exiting installer')
    sys.exit(0)

signal.signal(signal.SIGINT, sigint_handler)

def is_valid_url(url) -> bool:
    # Helper function to check if the URL is valid
    try:
        status_code = request.urlopen(url).getcode()
        if status_code == 200:
            return True
    except request.HTTPError:
        logging.exception(f"URL is not valid: {url}")
        raise

def search_and_replace(search_text, replace_text, file_path):
    # Common function to search and replace text in a file
    file = open(file_path, "r")
    for line in file:
        line = line.strip()
        if search_text in line:
            with open(file_path, 'r') as file:
                data = file.read()
                data = data.replace(line, replace_text)
            with open(file_path, 'w') as file:
                file.write(data)
    file.close()

def post_process(func):
    # Common function to post-process commands
    @functools.wraps(func)
    def wrapper(*args, **kwds):
        _allow_error = kwds.pop('allow_error', False)
        try:
            value = func(*args, **kwds)
        except subprocess.CalledProcessError as err:
            if err.returncode and _allow_error:
                return err
            logging.exception(err)
        return value
    return wrapper

def default_answer(func):
    # Common function to add default answer to questions
    @functools.wraps(func)
    def wrapper(*args, **kwds):
        _default = kwds.get('default', "")
        if _default:
            args = list(args)
            args[-1] += f" [default: {_default}]:{os.linesep}"
        value = func(*args)
        return value if value != "" else _default
    return wrapper


###############################################################
###         Release class: Get cheqd-node releases   	    ###
###############################################################
class Release:
    def __init__(self, release_map):
        self.version = release_map['tag_name']
        self.url = release_map['html_url']
        self.assets = release_map['assets']
        self.is_prerelease = release_map['prerelease']

    def get_release_url(self):
        # Construct the URL to download selected release from GitHub
        # This fetches the release tagged "latest", plus any other releases or pre-releases.
        # Release version numbers are in format "vX.Y.Z", but the release URL does not include the "v".
        # We also determine the OS and architecture, and construct the URL to download the release.
        try:
            os_arch = str.lower(platform.machine())
            # Python returns "x86_64" for 64-bit OS, but the release URL uses "amd64" since that's the Go convention.
            if os_arch == 'x86_64':
                os_arch = 'amd64'
            else:
                os_arch = 'arm64'
            os_name = str.lower(platform.system())
            for _url_item in self.assets:
                _url = _url_item["browser_download_url"]
                version_without_v_prefix = self.version.replace('v', '', 1)
                if os.path.basename(_url) == f"cheqd-noded-{version_without_v_prefix}-{os_name}-{os_arch}.tar.gz":
                    if is_valid_url(_url):
                        logging.debug(f"Release URL for binary download: {_url}")
                        return _url
                    else:
                        logging.exception(f"Release URL is not valid: {_url}")
            else:
                logging.exception(f"No asset found to download for release: {self.version}")
        except Exception as e:
            logging.exception(f"Failed to get cheqd-node binaries from GitHub. Reason: {e}")

    def __str__(self):
        return f"Name: {self.version}"


###############################################################
###         Installer class: Configure installation  	    ###
###############################################################
class Installer():
    def __init__(self, interviewer):
        self.version = interviewer.release.version
        self.release = interviewer.release
        self.interviewer = interviewer
        self._snapshot_url = ""
        
    @property
    def snapshot_url(self):
        return self._snapshot_url
    
    @snapshot_url.setter
    def snapshot_url(self, value):
        self._snapshot_url = value

    @property
    def binary_path(self):
        # Get the path to the cheqd-node binary on the local system
        return os.path.join(os.path.realpath(os.path.curdir), DEFAULT_BINARY_NAME)

    @property
    def cosmovisor_service_cfg(self):
        # Modify cheqd-cosmovisor.service template file to replace values for environment variables
        # The template file is fetched from the GitHub repo
        # Some of these variables are explicitly asked during the installer process. Others are set to default values.
        try:
            # Set service file path
            fname = os.path.basename(COSMOVISOR_SERVICE_TEMPLATE)

            # Fetch the template file from GitHub
            self.exec(f"wget -c {COSMOVISOR_SERVICE_TEMPLATE}")

            # Replace the values for environment variables in the template file
            with open(fname) as f:
                s = re.sub(
                    r'({CHEQD_ROOT_DIR}|{DEFAULT_BINARY_NAME}|{COSMOVISOR_DAEMON_ALLOW_DOWNLOAD_BINARIES}|{COSMOVISOR_DAEMON_RESTART_AFTER_UPGRADE}|{DEFAULT_DAEMON_POLL_INTERVAL}|{DEFAULT_UNSAFE_SKIP_BACKUP}|{DEFAULT_DAEMON_RESTART_DELAY})',
                    lambda m: {'{CHEQD_ROOT_DIR}': self.cheqd_root_dir,
                            '{DEFAULT_BINARY_NAME}': DEFAULT_BINARY_NAME,
                            '{COSMOVISOR_DAEMON_ALLOW_DOWNLOAD_BINARIES}':  self.interviewer.daemon_allow_download_binaries,
                            '{COSMOVISOR_DAEMON_RESTART_AFTER_UPGRADE}': self.interviewer.daemon_restart_after_upgrade,
                            '{DEFAULT_DAEMON_POLL_INTERVAL}': DEFAULT_DAEMON_POLL_INTERVAL,
                            '{DEFAULT_UNSAFE_SKIP_BACKUP}': DEFAULT_UNSAFE_SKIP_BACKUP,
                            '{DEFAULT_DAEMON_RESTART_DELAY}': DEFAULT_DAEMON_RESTART_DELAY}[m.group()],
                    f.read()
                )
            
            # Remove the template file
            self.remove_safe(fname)
            return s
        except Exception as e:
            logging.exception(f"Failed to set up service file from template. Reason: {e}")

    @property
    def rsyslog_cfg(self):
        # Modify rsyslog template file to replace values for environment variables
        # The template file is fetched from GitHub repo
        # Some of these variables are explicitly asked during the installer process. Others are set to default values.
        try:
            # Determine the binary name for logging based on installation type
            if self.interviewer.is_cosmovisor_needed:
                binary_name = DEFAULT_COSMOVISOR_BINARY_NAME
            else:
                binary_name = DEFAULT_BINARY_NAME

            # Set template file path
            fname = os.path.basename(RSYSLOG_TEMPLATE)

            # Fetch the template file from GitHub
            if is_valid_url(RSYSLOG_TEMPLATE):
                self.exec(f"wget -c {RSYSLOG_TEMPLATE}")
            else:
                logging.exception(f"URL is not valid: {RSYSLOG_TEMPLATE}")

            # Replace the values for environment variables in the template file
            with open(fname) as f:
                s = re.sub(
                    r'({BINARY_FOR_LOGGING}|{CHEQD_LOG_DIR})',
                    lambda m: {'{BINARY_FOR_LOGGING}': binary_name,
                                '{CHEQD_LOG_DIR}': self.cheqd_log_dir}[m.group()],
                    f.read()
                )
            
            # Remove the template file
            self.remove_safe(fname)
            return s
        except Exception as e:
            logging.exception(f"Failed to set up rsyslog from template. Reason: {e}")
        
    @property
    def logrotate_cfg(self):
        # Modify logrotate template file to replace values for environment variables
        # The logrotate template file is fetched from the GitHub repo
        # Logrotate is used to rotate the log files of the cheqd-node every day, and keep a maximum of 7 days of logs.
        try:
            # Set template file path
            fname = os.path.basename(LOGROTATE_TEMPLATE)

            # Fetch the template file from GitHub
            if is_valid_url(LOGROTATE_TEMPLATE):
                self.exec(f"wget -c {LOGROTATE_TEMPLATE}")
            else:
                logging.exception(f"URL is not valid: {LOGROTATE_TEMPLATE}")

            # Replace the values for environment variables in the template file
            with open(fname) as f:
                s = re.sub(
                    r'({CHEQD_LOG_DIR})',
                    lambda m: {'{CHEQD_LOG_DIR}': self.cheqd_log_dir}[m.group()],
                    f.read()
                )
            
            # Remove the template file
            self.remove_safe(fname)
            return s
        except Exception as e:
            logging.exception(f"Failed to set up logrotate from template. Reason: {e}")

    @property
    def cheqd_root_dir(self):
        # CHEQD_NODED_HOME variable can be picked up by cheqd-noded, so this should be set as an environment variable later
        # Default: /home/cheqd/.cheqdnode
        cheqd_noded_home = os.path.join(self.interviewer.home_dir, ".cheqdnode")
        return cheqd_noded_home

    @property
    def cheqd_config_dir(self):
        # cheqd-noded config directory
        # Default: /home/cheqd/.cheqdnode/config
        return os.path.join(self.cheqd_root_dir, "config")

    @property
    def cheqd_data_dir(self):
        # cheqd-noded data directory
        # Default: /home/cheqd/.cheqdnode/data
        return os.path.join(self.cheqd_root_dir, "data")

    @property
    def cheqd_log_dir(self):
        # cheqd-noded log directory
        # Default: /home/cheqd/.cheqdnode/log
        return os.path.join(self.cheqd_root_dir, "log")

    @property
    def cosmovisor_root_dir(self):
        # cosmovisor root directory
        # Default: /home/cheqd/.cheqdnode/cosmovisor
        return os.path.join(self.cheqd_root_dir, "cosmovisor")

    @property
    def cosmovisor_cheqd_bin_path(self):
        # cheqd-noded binary path if installed with cosmovisor
        # Default: /home/cheqd/.cheqdnode/cosmovisor/current/bin/cheqd-noded
        return os.path.join(self.cosmovisor_root_dir, f"current/bin/{DEFAULT_BINARY_NAME}")

    @property
    def cosmovisor_download_url(self):
        # Compute the download URL for cosmovisor binary based on the OS architecture and version number
        try:
            os_arch = platform.machine()
            if os_arch == 'x86_64':
                os_arch = 'amd64'
            else:
                os_arch = 'arm64'
            cosmovisor_download_url = COSMOVISOR_BINARY_URL.format(DEFAULT_LATEST_COSMOVISOR_VERSION, DEFAULT_LATEST_COSMOVISOR_VERSION, os_arch)
            if is_valid_url(cosmovisor_download_url):
                logging.debug(f"Cosmovisor download URL: {cosmovisor_download_url}")
                return cosmovisor_download_url
            else:
                logging.exception(f"Cosmovisor download URL is not valid: {cosmovisor_download_url}")
        except Exception as e:
            logging.exception(f"Failed to compute Cosmovisor download URL. Reason: {e}")

    @post_process
    def exec(self, cmd, use_stdout=True, suppress_err=False):
        # Helper function to safely execute shell commands
        logging.info(f"Executing command: {cmd}")
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

    def remove_safe(self, path, is_dir=False) -> bool:
        # Helper function to remove a file or directory safely
        try:
            if is_dir and os.path.exists(path):
                shutil.rmtree(path)
                logging.warning(f"Removed {path}")
                return True
            elif os.path.exists(path):
                os.remove(path)
                logging.warning(f"Removed {path}")
                return True
            else:
                logging.exception(f"{path} does not exist")
                return False
        except Exception as e:
            logging.exception(f"Failed to remove {path}. Reason: {e}")
            return False

    def install(self) -> bool:
        # Main function that controls calls to installation process functions
        try:
            # Download and extract cheqd-node binary
            if self.get_binary():
                logging.info("Successfully downloaded and extracted cheqd-noded binary")
            else:
                logging.error("Failed to download and extract binary")
                raise
            
            # Carry out pre-installation steps
            # Mostly relevant if installing from scratch or re-installing
            if self.pre_install():
                logging.info("Pre-installation steps completed successfully")
            else:
                logging.error("Failed to complete pre-installation steps")
                raise
            
            # Create cheqd user if it doesn't exist
            if self.prepare_cheqd_user():
                logging.info("User/group cheqd setup successfully")
            else:
                logging.error("Failed to setup user/group cheqd")
                raise
            
            # Setup directories needed for installation
            self.prepare_directory_tree()

            # Setup Cosmovisor binary if needed
            if self.interviewer.is_cosmovisor_needed or self.interviewer.is_cosmovisor_bump_needed:
                self.install_cosmovisor()

            # if self.interviewer.is_cosmovisor_bump_needed:
            #     logging.info("Bumping Cosmovisor")
            #     self.bump_cosmovisor()

            if not self.interviewer.is_cosmovisor_needed and not self.interviewer.is_cosmovisor_bump_needed:
                logging.info(
                    f"Moving binary from {self.binary_path} to {DEFAULT_INSTALL_PATH}")
                shutil.chown(self.binary_path,
                            DEFAULT_CHEQD_USER,
                            DEFAULT_CHEQD_USER)
                shutil.move(self.binary_path, os.path.join(
                    DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))
                self.configure_node_settings()
            
            self.set_cheqd_env_vars()
            
            if self.interviewer.is_setup_needed:
                self.configure_node_settings()

            if self.interviewer.init_from_snapshot:
                self.snapshot_url = self.get_snapshot_url()
                if self.snapshot_url:
                    logging.info(
                        "Downloading snapshot and extracting archive. This can take a *really* long time...")
                    self.download_snapshot()
                    self.untar_from_snapshot()

            self.setup_node_systemd()
            self.setup_logging_systemd()
            logging.info("The cheqd-noded binary has been successfully installed")
            return True
        except Exception as e:
            logging.exception(f"Failed to install cheqd-noded. Reason: {e}")
    
    def get_binary(self) -> bool:
        # Download cheqd-noded binary and extract it
        # Also remove the downloaded archive file, if applicable
        try:
            logging.info("Downloading cheqd-noded binary...")
            binary_url = self.release.get_release_url()
            fname = os.path.basename(binary_url)

            # Download the binary from GitHub
            self.exec(f"wget -c {binary_url}")
            
            # Check tar archive exists before extracting
            if fname.find(".tar.gz") != -1:
                # Extract the binary from the archive file
                # Using tarfile to extract is a safer option than just executing a command
                tar = tarfile.open(fname)
                tar.extractall()

                # Remove the archive file
                self.remove_safe(fname)

                # Make the binary executable
                # 0755 is equivalent to chmod +x
                os.chmod(DEFAULT_BINARY_NAME, 0o755)
                return True
            else:
                logging.error(f"Unable to extract cheqd-noded binary from archive file: {fname}")
                return False
        except Exception as e:
            logging.exception("Failed to download cheqd-noded binary. Reason: {e}")

    def pre_install(self) -> bool:
        # Pre-installation steps
        # Removes the following existing cheqd-noded data and configurations:
        # 1. ~/.cheqdnode directory
        # 2. cheqd-noded / cosmovisor binaries
        # 3. systemd service files
        try:
            if self.interviewer.is_from_scratch:
                logging.warning("Removing user's data and configs")

                # Remove cheqd-node service files safely first
                self.remove_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME, DEFAULT_STANDALONE_SERVICE_FILE_PATH)
                self.remove_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME, DEFAULT_COSMOVISOR_SERVICE_FILE_PATH)
                
                # Remove logging service files safely
                self.remove_safe(DEFAULT_RSYSLOG_FILE)
                self.remove_safe(DEFAULT_LOGROTATE_FILE)
                self.reload_systemd()

                # Remove cheqd-node data and binaries
                self.remove_safe(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_COSMOVISOR_BINARY_NAME))
                self.remove_safe(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))
                self.remove_safe(self.cheqd_root_dir, is_dir=True)
                return True

            # Scenario: User has installed cheqd-noded without cosmovisor, AND now wants to install cheqd-noded with Cosmovisor
            if self.interviewer.is_cosmovisor_needed and os.path.exists(DEFAULT_STANDALONE_SERVICE_FILE_PATH):
                self.remove_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME, DEFAULT_STANDALONE_SERVICE_FILE_PATH)
                self.remove_safe(DEFAULT_RSYSLOG_FILE)
                self.remove_safe(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))
                self.reload_systemd()
                return True
        except Exception as e:
            logging.exception("Could not complete pre-installation steps. Reason: {e}")
            raise
            return False

    def prepare_cheqd_user(self) -> bool:
        # Create "cheqd" user/group if it doesn't exist
        try:
            if not self.does_user_exist(DEFAULT_CHEQD_USER):
                logging.info(f"Creating {DEFAULT_CHEQD_USER} group")
                self.exec(f"addgroup {DEFAULT_CHEQD_USER} --quiet --system")
                logging.info(f"Creating {DEFAULT_CHEQD_USER} user and adding to {DEFAULT_CHEQD_USER} group")
                self.exec(
                    f"adduser --system {DEFAULT_CHEQD_USER} --home {self.interviewer.home_dir} --shell /bin/bash --ingroup {DEFAULT_CHEQD_USER} --quiet")
                return True
            else:
                logging.info(f"User {DEFAULT_CHEQD_USER} already exists. Skipping creation...")
                return True
        except Exception as e:
            logging.exception(f"Failed to create {DEFAULT_CHEQD_USER} user. Reason: {e}")
            return False

    def does_user_exist(self, username) -> bool:
        # Helper function to see if a given user exists on the system
        try:
            pwd.getpwnam(username)
            logging.debug(f"User {username} already exists")
            return True
        except KeyError:
            logging.debug(f"User {username} does not exist")
            return False

    def prepare_directory_tree(self):
        # Needed only in case of clean installation
        # 1. Create ~/.cheqdnode directory
        # 2. Set directory permissions to default cheqd user
        # 3. Create ~/.cheqdnode/log directory
        try:
            # Create root directory for cheqd-noded
            if not os.path.exists(self.cheqd_root_dir):
                logging.info("Creating main directory for cheqd-noded")
                os.makedirs(self.cheqd_root_dir)

                logging.info(f"Setting directory permissions to default cheqd user: {DEFAULT_CHEQD_USER}")
                shutil.chown(self.interviewer.home_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)
            else:
                logging.info(f"Skipping main directory creation because {self.cheqd_root_dir} already exists")

            # Setup logging related directories
            # 1. Create log directory if it doesn't exist
            # 2. Create default stdout.log file in log directory
            # 3. Set ownership of log directory to syslog:cheqd
            if not os.path.exists(self.cheqd_log_dir):
                # Create ~/.cheqdnode/log directory
                logging.info("Creating log directory for cheqd-noded")
                os.makedirs(self.cheqd_log_dir)

                # Create blank ~/.cheqdnode/log/stdout.log file
                Path(os.path.join(self.cheqd_log_dir, "stdout.log")).touch(exist_ok=True)

                logging.info(f"Setting up ownership permissions for {self.cheqd_log_dir} directory")
                shutil.chown(self.cheqd_log_dir, 'syslog', DEFAULT_CHEQD_USER)
            else:
                logging.info(f"Skipping log directory creation because {self.cheqd_log_dir} already exists")

            # Create symlink from cheqd-noded log folder from /var/log/cheqd-node
            # This step is necessary since many logging tools look for logs in /var/log
            if not os.path.exists("/var/log/cheqd-node"):
                logging.info("Creating a symlink from cheqd-noded log folder to /var/log/cheqd-node")
                os.symlink(self.cheqd_log_dir, "/var/log/cheqd-node", target_is_directory=True)
            else:
                logging.info("Skipping linking because /var/log/cheqd-node already exists")
        except Exception as e:
            logging.exception(f"Failed to prepare directory tree for {DEFAULT_CHEQD_USER}. Reason: {e}")

    def install_cosmovisor(self):
        # Install binaries for cheqd-noded and Cosmovisor
        # Cosmovisor is only installed if requested by the user
        # cheqd-noded binary is always installed, but the installation location depends whether user
        # chose to install with Cosmovisor or standalone
        try:
            logging.info("Setting up Cosmovisor")
            
            if self.get_cosmovisor():
                logging.info("Successfully downloaded Cosmovisor")

                # Set environment variables for Cosmovisor
                self.set_cosmovisor_env_vars()

                # Move Cosmovisor binary to installation directory if it doesn't exist or bump needed
                # This is executed is there is no Cosmovisor binary in the installation directory
                # or if the user has requested a bump for Cosmovisor
                logging.info(f"Moving Cosmovisor binary to {DEFAULT_INSTALL_PATH}/{DEFAULT_COSMOVISOR_BINARY_NAME}")
                shutil.move(DEFAULT_COSMOVISOR_BINARY_NAME, DEFAULT_INSTALL_PATH)

                # Check if Cosmovisor was successfully installed
                if self.interviewer.check_cosmovisor_installed():
                    logging.info("Cosmovisor successfully installed")

                    # Initialize Cosmovisor
                    self.exec(f"""su -l -c 'cosmovisor init ./{DEFAULT_BINARY_NAME}' {DEFAULT_CHEQD_USER}""")
                    
                    # Remove cheqd-noded binary from /usr/bin if it's not a symlink
                    if not os.path.islink(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME)):
                        logging.warn(f"Removing {DEFAULT_BINARY_NAME} from {DEFAULT_INSTALL_PATH} because it is not a symlink")
                        os.remove(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))

                        # Move cheqd-noded binary to Cosmovisor bin path
                        logging.info(f"Moving cheqd-noded binary from {self.binary_path} to {self.cosmovisor_cheqd_bin_path}")
                        shutil.move(DEFAULT_BINARY_NAME, self.cosmovisor_cheqd_bin_path)
                    else:
                        logging.debug(f"{DEFAULT_INSTALL_PATH}/{DEFAULT_BINARY_NAME} is already symlink. Skipping removal...")

                    # Create symlink to cheqd-noded binary in Cosmovisor bin path
                    if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME)):
                        logging.info(f"Creating symlink to {self.cosmovisor_cheqd_bin_path}")
                        os.symlink(self.cosmovisor_cheqd_bin_path, os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))

                    # Change owner of Cosmovisor directory to default cheqd user
                    logging.info(f"Changing ownership of {self.cosmovisor_root_dir} to {DEFAULT_CHEQD_USER} user")
                    shutil.chown(self.cosmovisor_root_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)
                else:
                    logging.error("Failed to install Cosmovisor")
                    raise
            else:
                logging.error("Failed to download Cosmovisor")
                raise
            
            # Steps to execute only if this is an upgrade
            # The upgrade-info.json file is required for Cosmovisor to function correctly
            if self.interviewer.is_upgrade and os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")):
                logging.info(f"Copying ~/.cheqdnode/data/upgrade-info.json file to ~/.cheqdnode/cosmovisor/current/")
                shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                    os.path.join(self.cosmovisor_root_dir, "current"))
            else:
                logging.debug("Skipping copying of upgrade-info.json file because it doesn't exist")
        except Exception as e:
            logging.exception(f"Failed to setup Cosmovisor. Reason: {e}")

    def bump_cosmovisor(self):
        try:
            self.stop_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME)
            self.remove_safe(fname)

            # Remove cosmovisor artifacts...
            self.remove_safe("CHANGELOG.md")
            self.remove_safe("README.md")
            self.remove_safe("LICENSE")

            # move the new binary to installation directory
            logging.info(f"Moving Cosmovisor binary to installation directory")
            shutil.move(os.path.join(os.path.realpath(os.path.curdir), DEFAULT_COSMOVISOR_BINARY_NAME), os.path.join(
                DEFAULT_INSTALL_PATH, DEFAULT_COSMOVISOR_BINARY_NAME))

            if not os.path.exists(os.path.join(self.cosmovisor_root_dir, "current")):
                logging.info(
                    f"Creating symlink for current Cosmovisor version")
                os.symlink(os.path.join(self.cosmovisor_root_dir, "genesis"),
                           os.path.join(self.cosmovisor_root_dir, "current"))

            logging.info(
                f"Moving binary from {self.binary_path} to {self.cosmovisor_cheqd_bin_path}")
            shutil.move(os.path.join(os.path.realpath(
                os.path.curdir), self.binary_path), os.path.join(self.cosmovisor_cheqd_bin_path))
            shutil.chown(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_COSMOVISOR_BINARY_NAME),
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)
            self.exec(
                "sudo chmod +x {}".format(f'{DEFAULT_INSTALL_PATH}/{DEFAULT_COSMOVISOR_BINARY_NAME}'))

            if not os.path.exists(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME)):
                logging.info(
                    f"Creating symlink to {self.cosmovisor_cheqd_bin_path}")
                os.symlink(self.cosmovisor_cheqd_bin_path,
                           os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))

            if self.interviewer.is_upgrade and \
                    os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")):
                logging.info(
                    f"Copying upgrade-info.json file to cosmovisor/current/")
                shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                            os.path.join(self.cosmovisor_root_dir, "current"))
                logging.info(f"Changing owner to {DEFAULT_CHEQD_USER} user")
                shutil.chown(self.cosmovisor_root_dir,
                             DEFAULT_CHEQD_USER,
                             DEFAULT_CHEQD_USER)

            logging.info(
                f"Changing directory ownership for Cosmovisor to {DEFAULT_CHEQD_USER} user")
            shutil.chown(self.cosmovisor_root_dir,
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)
            self.reload_systemd()
        except Exception as e:
            logging.exception(f"Failed to bump Cosmovisor. Reason: {e}")

    def get_cosmovisor(self) -> bool:
        # Download Cosmovisor binary and extract it
        # Also remove the downloaded archive file, if applicable
        try:
            logging.info("Downloading Cosmovisor binary...")
            cosmovisor_download_url = self.download_and_unzip(self.cosmovisor_download_url)
            fname = os.path.basename(cosmovisor_download_url)

            # Download Cosmovisor binary from GitHub
            self.exec(f"wget -c {cosmovisor_download_url}")

            # Check tar archive exists before extracting
            if fname.find(".tar.gz") != -1:
                # Extract Cosmovisor binary from the archive file
                # Using tarfile to extract is a safer option than just executing a command
                tar = tarfile.open(fname)
                tar.extractall()

                # Remove Cosmovisor artifacts...
                self.remove_safe("CHANGELOG.md")
                self.remove_safe("README.md")
                self.remove_safe("LICENSE")
                self.remove_safe(fname)

                # Make the binary executable
                # 0755 is equivalent to chmod +x
                os.chmod(DEFAULT_COSMOVISOR_BINARY_NAME, 0o755)
                return True
            else:
                logging.error(f"Unable to extract Cosmovisor binary from archive file: {fname}")
                return False
        except Exception as e:
            logging.exception("Failed to download Cosmovisor binary. Reason: {e}")

    def set_cosmovisor_env_vars(self) -> bool:
        # Set environment variables for Cosmovisor
        try:
            self.set_environment_variable("DAEMON_NAME", DEFAULT_BINARY_NAME, overwrite=True)
            self.set_environment_variable("DAEMON_HOME", self.cheqd_root_dir, overwrite=False)
            self.set_environment_variable("DAEMON_ALLOW_DOWNLOAD_BINARIES", 
                self.interviewer.daemon_allow_download_binaries, overwrite=True)
            self.set_environment_variable("DAEMON_RESTART_AFTER_UPGRADE",
                self.interviewer.daemon_restart_after_upgrade, overwrite=True)
            self.set_environment_variable("DAEMON_POLL_INTERVAL", DEFAULT_DAEMON_POLL_INTERVAL, overwrite=False)
            self.set_environment_variable("UNSAFE_SKIP_BACKUP", DEFAULT_UNSAFE_SKIP_BACKUP, overwrite=False)
        except Exception as e:
            logging.exception(f"Failed to set environment variables for Cosmovisor. Reason: {e}")

    def set_environment_variable(self, env_var_name, env_var_value, overwrite=True):
        # Set an environment variable
        # By default, existing environment variables are overwritten
        # This can be changed by setting the overwrite parameter to False
        # Environment variables are set for the current session as well as for all users
        try:
            logging.debug(f"Checking whether {env_var_name} is set")

            if not os.environ(env_var_name) or overwrite:
                logging.debug(f"Setting {env_var_name} to {env_var_value}")
                
                # Set the environment variable for the current session
                os.environ[env_var_name] = env_var_value

                # Modify the system's environment variables
                # This will set the variable permanently for all users
                with open("/etc/environment", "a") as env_file:
                    env_file.write(f"\n{env_var_name}={env_var_value}")
                
                # Reload the environment variables
                os.system("source /etc/environment")
            else:
                logging.debug(f"Environment variable {env_var_name} already set or overwrite is disabled")
        except Exception as e:
            logging.exception(f"Failed to set environment variable {env_var_name}. Reason: {e}")
        finally:
            env_file.close()

    def check_systemd_service_active(self, service_name) -> bool:
        # Check if a given systemd service is active
        try:
            logging.debug(f"Checking whether {service_name} service is active")
            # Check if the service exists
            cmd_exists = f'systemctl status {service_name}.service'
            exists = os.system(cmd_exists)

            if exists == 0:
                # Check if the service is active
                cmd_active = f'systemctl is-active --quiet {service_name}.service'
                active = os.system(cmd_active)
                if active == 0:
                    logging.debug(f"Service {service_name} is active")
                    return True
                else:
                    logging.debug(f"Service {service_name} is not active")
                    return False
            else:
                logging.debug(f"Service {service_name} is not installed")
                return False
        except Exception as e:
            logging.exception(f"Failed to check whether {service_name} service is active. Reason: {e}")
            return False

    def check_systemd_service_enabled(self, service_name) -> bool:
        # Check if a given systemd service is enabled
        try:
            logging.debug(f"Checking whether {service_name} service is enabled")
            # Check if the service exists
            cmd_exists = f'systemctl status {service_name}.service'
            exists = os.system(cmd_exists)
            if exists == 0:
                # Check if the service is enabled
                cmd_enabled = f'systemctl is-enabled --quiet {service_name}.service'
                enabled = os.system(cmd_enabled)
                if enabled == 0:
                    logging.debug(f"Service {service_name} is enabled")
                    return True
                else:
                    logging.debug(f"Service {service_name} is not enabled")
                    return False
            else:
                logging.debug(f"Service {service_name} is not installed")
                return False
        except Exception as e:
            logging.exception(f"Failed to check whether {service_name} service is enabled. Reason: {e}")
    
    def reload_systemd(self) -> bool:
        # Reload systemd config
        try:
            logging.debug("Reload systemd config and reset failed services")

            # Reload systemd config
            reload = os.system(f'systemctl daemon-reload --quiet')

            # Reset failed services
            reset = os.system(f'systemctl reset-failed --quiet')

            if reload == 0 and reset == 0:
                logging.info("Reloaded systemd config and reset failed services")
                return True
            else:
                logging.error("Failed to reload systemd config and reset failed services")
                return False
                raise
        except Exception as e:
            logging.exception(f"Error disabling {service_name}: Reason: {e}")
    
    def disable_systemd_service(self, service_name) -> bool:
        # Disable a given systemd service
        try:
            if self.check_systemd_service_enabled(service_name):
                disabled = os.system(f"systemctl disable --quiet {service_name}.service")
                if disabled == 0:
                    logging.info(f"{service_name} has been disabled")
                    return True
                else:
                    logging.error(f"{service_name} could not be disabled")
                    return False
                    raise
            else:
                logging.debug(f"{service_name} is already disabled")
                return True
        except Exception as e:
            logging.exception(f"Error disabling {service_name}: Reason: {e}")

    def enable_systemd_service(self, service_name) -> bool:
        # Enable a given systemd service
        try:
            if self.reload_systemd():
                if not self.check_systemd_service_enabled(service_name):
                    enabled = os.system(f"systemctl enable --quiet {service_name}.service")
                    if enabled == 0:
                        logging.info(f"{service_name} has been enabled")
                        return True
                    else:
                        logging.error(f"{service_name} could not be enabled")
                        return False
                        raise
                else:
                    logging.debug(f"{service_name} is already enabled")
                    return True
            else:
                logging.error(f"Failed to reload systemd config and reset failed services")
                return False
                raise
        except Exception as e:
            logging.exception(f"Error disabling {service_name}: Reason: {e}")

    def stop_systemd_service(self, service_name) -> bool:
        # Stop and disable a given systemd service
        try:
            if self.check_systemd_service_active(service_name):
                stopped = os.system(f"systemctl stop --quiet {service_name}.service")
                if stopped == 0:
                    logging.info(f"{service_name} has been stopped")
                    return True
                else:
                    logging.error(f"{service_name} could not be stopped")
                    return False
                    raise
            else:
                logging.debug(f"{service_name} is not active")
                return True
        except Exception as e:
            logging.exception(f"Error stopping {service_name}: Reason: {e}")
            return False

    def restart_systemd_service(self, service_name) -> bool:
        # Restart a given systemd service
        try:
            # If the service is not enabled, enable it before restarting
            if not self.check_systemd_service_enabled(service_name):
                self.enable_systemd_service(service_name)

            # Reload systemd services before restarting
            if self.reload_systemd():
                restarted = os.system(f"systemctl restart --quiet {service_name}.service")
                if restarted == 0:
                    logging.info(f"{service_name} has been restarted")
                    return True
                else:
                    logging.error(f"{service_name} could not be restarted")
                    return False
                    raise
            else:
                logging.error(f"Failed to restart {service_name}")
                return False
                raise
        except Exception as e:
            logging.exception(f"Error restarting {service_name}: Reason: {e}")
            raise

    def remove_systemd_service(self, service_name, service_file) -> bool:
        # Remove a given systemd service
        try:
            if os.path.exists(service_file):
                # Stop the service if it is active
                if self.stop_systemd_service(service_name):
                    # Disable the service
                    if self.disable_systemd_service(service_name):
                        # Remove the service file
                        self.remove_safe(service_file)
                        logging.warning(f"{service_name} has been removed")
                        return True
                    else:
                        logging.error(f"{service_name} could not be removed")
                        return False
        except Exception as e:
            logging.exception(f"Error removing {service_name}: Reason: {e}")

    def setup_node_systemd(self):
        # Setup cheqd-noded related systemd services
        # If user selected Cosmovisor install, then cheqd-cosmovisor.service will be setup
        # If user selected Standalone install, then cheqd-noded.service will be setup
        try:
            logging.info("Setting up systemd config")

            # Check if systemd service files already exist for cheqd-node
            service_file_exists = os.path.exists(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH) or os.path.exists(
                DEFAULT_STANDALONE_SERVICE_FILE_PATH)

            # WARNING: Revisit this logic and check the condition
            if not self.interviewer.is_upgrade or \
                    self.interviewer.rewrite_node_systemd or \
                    not service_file_exists:
                self.remove_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME, DEFAULT_COSMOVISOR_SERVICE_FILE_PATH)
                self.remove_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME, DEFAULT_STANDALONE_SERVICE_FILE_PATH)

                if self.interviewer.is_cosmovisor_needed:
                    # Setup cheqd-cosmovisor.service if requested
                    logging.info("Enabling cheqd-cosmovisor.service in systemd")
                    with open(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH, mode="w") as fd:
                        fd.write(self.cosmovisor_service_cfg)
                    self.enable_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME)
                else:
                    logging.info("Enabling cheqd-noded.service in systemd")
                    self.exec(f"curl -s {STANDALONE_SERVICE_TEMPLATE} > {DEFAULT_STANDALONE_SERVICE_FILE_PATH}")
                    self.enable_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME)
        except Exception as e:
            logging.exception(f"Failed to setup systemd service for cheqd-node. Reason: {e}")

    # Setup logging related systemd services
    def setup_logging_systemd(self):
        # Install cheqd-node configuration for rsyslog if either of the following conditions are met:
        # 1. rsyslog service file is not present
        # 2. user wants to rewrite rsyslog service file
        try:
            if not os.path.exists(DEFAULT_RSYSLOG_FILE) or self.interviewer.rewrite_rsyslog:
                # Warn user if rsyslog service file already exists
                if os.path.exists(DEFAULT_RSYSLOG_FILE):
                    logging.warning(f"Existing rsyslog configuration at {DEFAULT_RSYSLOG_FILE} will be overwritten")
                
                # Determine the binary name for logging based on installation type
                if self.interviewer.is_cosmovisor_needed:
                    binary_name = DEFAULT_COSMOVISOR_BINARY_NAME
                else:
                    binary_name = DEFAULT_BINARY_NAME

                logging.info(f"Configuring rsyslog systemd service for {binary_name} logging")

                # Modify rsyslog template file with values specific to the installation
                with open(DEFAULT_RSYSLOG_FILE, mode="w") as fname:
                    fname.write(self.rsyslog_cfg)

                # Restarting rsyslog can take a lot of time: https://github.com/rsyslog/rsyslog/issues/3133
                if self.restart_systemd_service("rsyslog.service"):
                    logging.info("Successfully configured rsyslog service")
                else:
                    logging.error("Failed to configure rsyslog service")
                    raise
        except Exception as e:
            logging.exception(f"Failed to setup rsyslog service for {binary_name} logging. Reason: {e}")

        # Install cheqd-node configuration for logrotate if either of the following conditions are met:
        # 1. logrotate service file is not present
        # 2. user wants to rewrite logrotate service file
        if not os.path.exists(DEFAULT_LOGROTATE_FILE) or self.interviewer.rewrite_logrotate:
            try:
                # Warn user if logrotate service file already exists
                if os.path.exists(DEFAULT_LOGROTATE_FILE):
                    logging.warning(f"Existing logrotate configuration at {DEFAULT_LOGROTATE_FILE} will be overwritten")

                logging.info(f"Configuring logrotate systemd service for cheqd-node logging")

                # Modify logrotate template file with values specific to the installation
                with open(DEFAULT_LOGROTATE_FILE, mode="w") as fname:
                    fname.write(self.logrotate_cfg)

                if self.restart_systemd_service("logrotate.service"):
                    logging.info("Successfully configured logrotate service")
                    if self.restart_systemd_service("logrotate.timer"):
                        logging.info("Successfully configured logrotate timer")
                    else:
                        logging.exception("Failed to configure logrotate timer")
                else:
                    logging.exception("Failed to configure logrotate service")
            except Exception as e:
                logging.exception(
                    f"Failed to setup logrotate service. Reason: {e}")

    def configure_node_settings(self):
        # Init the node with provided moniker
        if not os.path.exists(os.path.join(self.cheqd_config_dir, 'genesis.json')):
            self.exec(
                f"""sudo su -c 'cheqd-noded init {self.interviewer.moniker}' {DEFAULT_CHEQD_USER}""")

            # Downloading genesis file
            self.exec(
                f"curl {GENESIS_FILE.format(self.interviewer.chain)} > {os.path.join(self.cheqd_config_dir, 'genesis.json')}")
            shutil.chown(os.path.join(self.cheqd_config_dir, 'genesis.json'),
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)

        # Replace the default RCP port to listen to anyone
        rpc_default_value = 'laddr = "tcp://127.0.0.1:{}"'.format(
            DEFAULT_RPC_PORT)
        new_rpc_default_value = 'laddr = "tcp://0.0.0.0:{}"'.format(
            DEFAULT_RPC_PORT)
        search_and_replace(rpc_default_value, new_rpc_default_value, os.path.join(
            self.cheqd_config_dir, "config.toml"))

        # Set create empty blocks to false by default
        create_empty_blocks_search_text = 'create_empty_blocks = true'
        create_empty_blocks_replace_text = 'create_empty_blocks = false'
        search_and_replace(create_empty_blocks_search_text, create_empty_blocks_replace_text, os.path.join(
            self.cheqd_config_dir, "config.toml"))

        # Setting up the external_address
        if self.interviewer.external_address:
            external_address_search_text = 'external_address = ""'
            external_address_replace_text = 'external_address = "{}:{}"'.format(
                self.interviewer.external_address, self.interviewer.p2p_port)
            search_and_replace(external_address_search_text, external_address_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))

        # Setting up the seeds
        seeds = self.exec(
            f"curl {SEEDS_FILE.format(self.interviewer.chain)}").stdout.decode("utf-8").strip()
        seeds_search_text = 'seeds = ""'
        seeds_replace_text = 'seeds = "{}"'.format(seeds)
        search_and_replace(seeds_search_text, seeds_replace_text, os.path.join(
            self.cheqd_config_dir, "config.toml"))

        # Setting up the RPC port
        if self.interviewer.rpc_port:
            rpc_laddr_search_text = 'laddr = "tcp://0.0.0.0:{}"'.format(
                DEFAULT_RPC_PORT)
            rpc_laddr_replace_text = 'laddr = "tcp://0.0.0.0:{}"'.format(
                self.interviewer.rpc_port)
            search_and_replace(rpc_laddr_search_text, rpc_laddr_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))
        # Setting up the P2P port
        if self.interviewer.p2p_port:
            p2p_laddr_search_text = 'laddr = "tcp://0.0.0.0:{}"'.format(
                DEFAULT_P2P_PORT)
            p2p_laddr_replace_text = 'laddr = "tcp://0.0.0.0:{}"'.format(
                self.interviewer.p2p_port)
            search_and_replace(p2p_laddr_search_text, p2p_laddr_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))

        # Setting up min gas-price
        if self.interviewer.gas_price:
            min_gas_price_search_text = 'minimum-gas-prices = '
            min_gas_price_replace_text = 'minimum-gas-prices = "{}"'.format(
                self.interviewer.gas_price)
            search_and_replace(min_gas_price_search_text, min_gas_price_replace_text, os.path.join(
                self.cheqd_config_dir, "app.toml"))

        # Setting up persistent peers
        if self.interviewer.persistent_peers:
            persistent_peers_search_text = 'persistent_peers = ""'
            persistent_peers_replace_text = 'persistent_peers = "{}"'.format(
                self.interviewer.persistent_peers)
            search_and_replace(persistent_peers_search_text, persistent_peers_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))

        # Setting up log level
        if self.interviewer.log_level:
            log_level_search_text = 'log_level'
            log_level_replace_text = 'log_level = "{}"'.format(
                self.interviewer.log_level)
            search_and_replace(log_level_search_text, log_level_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))
        else:
            log_level_search_text = 'log_level'
            log_level_replace_text = 'log_level = "{}"'.format(
                CHEQD_NODED_LOG_LEVEL)
            search_and_replace(log_level_search_text, log_level_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))

        # Setting up log format
        if self.interviewer.log_format:
            log_format_search_text = 'log_format'
            log_format_replace_text = 'log_format = "{}"'.format(
                self.interviewer.log_format)
            search_and_replace(log_format_search_text, log_format_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))
        else:
            log_format_search_text = 'log_format'
            log_format_replace_text = 'log_format = "{}"'.format(
                CHEQD_NODED_LOG_FORMAT)
            search_and_replace(log_format_search_text, log_format_replace_text, os.path.join(
                self.cheqd_config_dir, "config.toml"))

    def mkdir_p(self, dir_name):
        try:
            os.mkdir(dir_name)
        except FileExistsError as err:
            logging.info(f"Directory {dir_name} already exists")

    def set_cheqd_env_vars(self):
        self.set_environment_variable("DEFAULT_CHEQD_HOME_DIR",
                          f"{self.interviewer.cheqd_root_dir}")
        self.set_environment_variable("CHEQD_NODED_CHAIN_ID", f"{self.interviewer.chain}")

    def compare_checksum(self, file_path):
        # Set URL for correct checksum file for snapshot
        checksum_url = os.path.join(os.path.dirname(
            self.snapshot_url), "md5sum.txt")
        # Get checksum file
        published_checksum = self.exec(
            f"curl -s {checksum_url} | tail -1 | cut -d' ' -f 1").stdout.strip()
        logging.info(f"Comparing published checksum with local checksum")
        local_checksum = self.exec(
            f"md5sum {file_path} | tail -1 | cut -d' ' -f 1").stdout.strip()
        if published_checksum == local_checksum:
            logging.info(f"Checksums match. Download is OK.")
            return True
        elif published_checksum != local_checksum:
            logging.info(f"Checksums do not match. Download got corrupted.")
            return False
        else:
            logging.exception(f"Error encountered when comparing checksums.")

    def install_dependencies(self):
        try:
            logging.info("Installing dependencies")
            self.exec("sudo apt-get update")
            logging.info(f"Install pv to show progress of extraction")
            self.exec("sudo apt-get install -y pv")
        except Exception as e:
            logging.exception(f"Failed to install dependencies. Reason: {e}")
            
    def get_snapshot_url(self) -> str:
        template = TESTNET_SNAPSHOT if self.interviewer.chain in TESTNET_CHAIN_ID else MAINNET_SNAPSHOT
        _date = datetime.date.today()
        _days_counter = 0
        _is_url_valid = False

        while not _is_url_valid and _days_counter <= MAX_SNAPSHOT_DAYS:
            _url = template.format(_date.strftime(
                "%Y-%m-%d"), _date.strftime("%Y-%m-%d"))
            _is_url_valid = is_valid_url(_url)
            _days_counter += 1
            _date -= datetime.timedelta(days=1)

        if not _is_url_valid:
            logging.exception("Could not find the valid snapshot for the last {} days".format(
                MAX_SNAPSHOT_DAYS))
        return _url

    def download_snapshot(self):
        try:
            archive_name = os.path.basename(self.snapshot_url)
            self.mkdir_p(self.cheqd_data_dir)
            shutil.chown(self.cheqd_data_dir,
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)
            # Fetch size of snapshot archive. Uses curl to fetch headers and looks for Content-Length.
            archive_size = self.exec(
                f"curl -s --head {self.snapshot_url} | awk '/content-length/ {{print $2}}'").stdout.strip()
            # Check how much free disk space is available wherever the cheqd root directory is mounted
            free_disk_space = self.exec(
                f"df -P -B1 {self.cheqd_root_dir} | tail -1 | awk '{{print $4}}'").stdout.strip()
            if int(archive_size) < int(free_disk_space):
                logging.info(
                    f"Downloading snapshot archive. This may take a while...")
                self.exec(
                    f"wget -c {self.snapshot_url} -P {self.cheqd_root_dir}")
                archive_path = os.path.join(self.cheqd_root_dir, archive_name)
                if self.compare_checksum(archive_path) is True:
                    logging.info(
                        f"Snapshot download was successful and checksums match.")
                else:
                    logging.info(
                        f"Snapshot download was successful but checksums do not match.")
                    logging.exception(
                        f"Snapshot download was successful but checksums do not match.")
            elif int(archive_size) > int(free_disk_space):
                logging.exception(
                    f"Snapshot archive is too large to fit in free disk space. Please free up some space and try again.")
            else:
                logging.exception(
                    f"Error encountered when downloading snapshot archive.")
        except Exception as e:
            logging.exception(f"Failed to download snapshot. Reason: {e}")

    def untar_from_snapshot(self):
        try:
            archive_path = os.path.join(
                self.cheqd_root_dir, os.path.basename(self.snapshot_url))
            # Check if there is enough space to extract snapshot archive
            self.install_dependencies()
            logging.info(
                f"Extracting snapshot archive. This may take a while...")

            # Extract to cheqd node data directory EXCEPT for validator state
            self.exec(
                f"sudo su -c 'pv {archive_path} | tar --use-compress-program=lz4 -xf - -C {self.cheqd_root_dir} --exclude priv_validator_state.json' {DEFAULT_CHEQD_USER}")

            # Delete snapshot archive file
            logging.info(
                f"Snapshot extraction was successful. Deleting snapshot archive.")
            self.remove_safe(archive_path)
            # Workaround to make this work with Cosmovisor since it expects upgrade-info.json file in cosmovisor/current directory
            if self.interviewer.is_cosmovisor_needed:
                if os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")):
                    logging.info(
                        f"Copying upgrade-info.json file to cosmovisor/current/")
                    shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                                os.path.join(self.cosmovisor_root_dir, "current"))
                logging.info(f"Changing owner to {DEFAULT_CHEQD_USER} user")
                shutil.chown(self.cosmovisor_root_dir,
                             DEFAULT_CHEQD_USER,
                             DEFAULT_CHEQD_USER)
            shutil.chown(self.cheqd_data_dir,
                         DEFAULT_CHEQD_USER,
                         DEFAULT_CHEQD_USER)
        except Exception as e:
            logging.exception(f"Failed to extract snapshot. Reason: {e}")


###############################################################
###         Interviewer class: Ask user for settings  	    ###
###############################################################
class Interviewer:
    def __init__(self, home_dir=DEFAULT_CHEQD_HOME_DIR, chain=CHEQD_NODED_CHAIN_ID):
        self._home_dir = home_dir
        self._is_upgrade = False
        self._is_cosmovisor_needed = True
        self._is_cosmovisor_bump_needed = True
        self._is_cosmovisor_installed = False
        self._systemd_service_file = ""
        self._init_from_snapshot = False
        self._release = None
        self._chain = chain
        self._is_setup_needed = False
        self._moniker = CHEQD_NODED_MONIKER
        self._external_address = ""
        self._rpc_port = DEFAULT_RPC_PORT
        self._p2p_port = DEFAULT_P2P_PORT
        self._gas_price = CHEQD_NODED_MINIMUM_GAS_PRICES
        self._persistent_peers = ""
        self._log_level = CHEQD_NODED_LOG_LEVEL
        self._log_format = CHEQD_NODED_LOG_FORMAT
        self._daemon_allow_download_binaries = DEFAULT_DAEMON_ALLOW_DOWNLOAD_BINARIES
        self._daemon_restart_after_upgrade = DEFAULT_DAEMON_RESTART_AFTER_UPGRADE
        self._is_from_scratch = False
        self._rewrite_node_systemd = False
        self._rewrite_rsyslog = False
        self._rewrite_logrotate = False

    ### This section sets @property variables ###
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
    def systemd_service_file(self) -> str:
        return self._systemd_service_file

    @property
    def rewrite_node_systemd(self) -> bool:
        return self._rewrite_node_systemd

    @property
    def rewrite_rsyslog(self) -> bool:
        return self._rewrite_rsyslog

    @property
    def rewrite_logrotate(self) -> bool:
        return self._rewrite_logrotate

    @property
    def is_cosmovisor_needed(self) -> bool:
        return self._is_cosmovisor_needed

    @property
    def is_cosmovisor_bump_needed(self) -> bool:
        return self._is_cosmovisor_bump_needed

    @property
    def is_cosmovisor_installed(self) -> bool:
        return self._is_cosmovisor_installed

    @property
    def init_from_snapshot(self) -> bool:
        return self._init_from_snapshot

    @property
    def chain(self) -> str:
        return self._chain

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

    @property
    def persistent_peers(self) -> str:
        return self._persistent_peers

    @property
    def log_level(self) -> str:
        return self._log_level

    @property
    def log_format(self) -> str:
        return self._log_format

    @property
    def daemon_allow_download_binaries(self) -> str:
        return self._daemon_allow_download_binaries

    @property
    def daemon_restart_after_upgrade(self) -> str:
        return self._daemon_restart_after_upgrade

    ### This section sets @property variables ###
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

    @systemd_service_file.setter
    def systemd_service_file(self, ssf):
        self._systemd_service_file = ssf

    @rewrite_node_systemd.setter
    def rewrite_node_systemd(self, rns):
        self._rewrite_node_systemd = rns

    @rewrite_rsyslog.setter
    def rewrite_rsyslog(self, rr):
        self._rewrite_rsyslog = rr

    @rewrite_logrotate.setter
    def rewrite_logrotate(self, rl):
        self._rewrite_logrotate = rl

    @is_cosmovisor_needed.setter
    def is_cosmovisor_needed(self, icn):
        self._is_cosmovisor_needed = icn

    @is_cosmovisor_bump_needed.setter
    def is_cosmovisor_bump_needed(self, icbn):
        self._is_cosmovisor_bump_needed = icbn

    @is_cosmovisor_installed.setter
    def is_cosmovisor_installed(self, icbn):
        self._is_cosmovisor_installed = ici

    @init_from_snapshot.setter
    def init_from_snapshot(self, ifs):
        self._init_from_snapshot = ifs

    @chain.setter
    def chain(self, chain):
        self._chain = chain

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

    @persistent_peers.setter
    def persistent_peers(self, persistent_peers):
        self._persistent_peers = persistent_peers

    @log_level.setter
    def log_level(self, log_level):
        self._log_level = log_level

    @log_format.setter
    def log_format(self, log_format):
        self._log_format = log_format

    @daemon_allow_download_binaries.setter
    def daemon_allow_download_binaries(self, daemon_allow_download_binaries):
        self._daemon_allow_download_binaries = daemon_allow_download_binaries

    @daemon_restart_after_upgrade.setter
    def daemon_restart_after_upgrade(self, daemon_restart_after_upgrade):
        self._daemon_restart_after_upgrade = daemon_restart_after_upgrade

    ### This section contains helper functions for the interviewer ###

    # Set value to default answer for a question
    @default_answer
    def ask(self, question, **kwargs):
        return str(input(question)).strip()

    @post_process
    def exec(self, cmd, use_stdout=True, suppress_err=False, check=True):
        logging.info(f"Executing command: {cmd}")
        kwargs = {
            "shell": True,
            "check": check,
        }
        if use_stdout:
            kwargs["stdout"] = subprocess.PIPE
        else:
            kwargs["capture_output"] = True

        if suppress_err:
            kwargs["stderr"] = subprocess.DEVNULL
        return subprocess.run(cmd, **kwargs)

    # Check if cheqd-noded is installed
    def is_node_installed(self) -> bool:
        try:
            if shutil.which("cheqd-noded") is not None:
                return True
            else:
                return False
        except Exception as e:
            logging.exception(
                f"Could not check if cheqd-noded is already installed. Reason: {e}")

    # Check if Cosmovisor is installed
    def check_cosmovisor_installed(self) -> bool:
        try:
            if shutil.which("cosmovisor") is not None:
                self.is_cosmovisor_installed = True
                return True
            else:
                self.is_cosmovisor_installed = False
                return False
        except Exception as e:
            logging.exception(f"Could not check if Cosmovisor is already installed. Reason: {e}")

    # Check if a systemd config is installed for a given service file
    def is_systemd_config_installed(self, systemd_service_file) -> bool:
        try:
            if os.path.exists(systemd_service_file):
                return True
            else:
                return False
        except Exception as e:
            logging.exception(
                f"Could not check if {systemd_service_file} already exists. Reason: {e}")

    # Get list of last N releases for cheqd-node from GitHub
    def get_releases(self):
        try:
            req = request.Request(
                "https://api.github.com/repos/cheqd/cheqd-node/releases")
            req.add_header("Accept", "application/vnd.github.v3+json")
            with request.urlopen(req) as response:
                r_list = json.loads(response.read().decode("utf-8").strip())
                return [Release(r) for r in r_list]
        except Exception as e:
            logging.exception(
                f"Could not get releases from GitHub. Reason: {e}")

    # The "latest" stable release may not be in last N releases, so we need to get it separately
    def get_latest_release(self):
        try:
            req = request.Request(
                "https://api.github.com/repos/cheqd/cheqd-node/releases/latest")
            req.add_header("Accept", "application/vnd.github.v3+json")
            with request.urlopen(req) as response:
                return Release(json.loads(response.read().decode("utf-8")))
        except Exception as e:
            logging.exception(
                f"Could not get latest release from GitHub. Reason: {e}")

    # Compile a list of releases to be displayed to the user
    # The "latest" stable release is always displayed first
    def remove_release_from_list(self, r_list, elem):
        try:
            copy_r_list = copy.deepcopy(r_list)
            for i, release in enumerate(r_list):
                if release.version == elem.version:
                    copy_r_list.pop(i)
                    return copy_r_list
        except Exception as e:
            logging.exception(
                f"Could not assemble list of releases to show to the user. Reason: {e}")

    # Ask user to select a version of cheqd-node to install
    def ask_for_version(self):
        try:
            default = self.get_latest_release()
            all_releases = self.get_releases()
            all_releases = self.remove_release_from_list(all_releases, default)
            all_releases.insert(0, default)

            print(f"Latest stable cheqd-noded release version is {default}")
            print(f"List of cheqd-noded releases: ")

            # Print list of releases
            for i, release in enumerate(all_releases[0: LAST_N_RELEASES]):
                print(f"{i + 1}. {release.version}")

            release_num = int(self.ask(
                "Choose list option number above to select version of cheqd-node to install", default=1))

            # Check if user input is valid
            if release_num >= 1 and release_num <= LAST_N_RELEASES:
                self.release = all_releases[release_num - 1]
            else:
                raise ValueError(
                    f"Invalid release number picked from list of releases: {release_num}")

        except Exception as e:
            logging.exception(
                f"Failed to selected version of cheqd-noded. Reason: {e}")

    # Set cheqd user's home directory
    def ask_for_home_directory(self) -> str:
        try:
            self.home_dir = self.ask(
                f"Set path for cheqd user's home directory", default=DEFAULT_CHEQD_HOME_DIR)
        except Exception as e:
            logging.exception(
                f"Failed to set cheqd user's home directory. Reason: {e}")

    # Ask whether user wants to do a install from scratch
    def ask_for_setup(self):
        try:
            answer = self.ask(
                f"Do you want to setup a new cheqd-node installation? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.is_setup_needed = True
            elif answer.lower().startswith("n"):
                self.is_setup_needed = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'")
                self.ask_for_setup()
        except Exception as e:
            logging.exception(
                f"Failed to set fresh installation parameters. Reason: {e}")

    # Ask user which network to join
    def ask_for_chain(self):
        try:
            answer = int(self.ask(
                "Select cheqd network to join:\n"
                f"1. Mainnet ({MAINNET_CHAIN_ID})\n"
                f"2. Testnet ({TESTNET_CHAIN_ID})\n", default=1))
            if answer == 1:
                self.chain = "mainnet"
            elif answer == 2:
                self.chain = "testnet"
            else:
                logging.error(
                    f"Invalid network selected during installation. Please choose either 1 or 2.")
                self.ask_for_chain()
        except Exception as e:
            logging.exception(
                f"Failed to set network/chain to join. Reason: {e}")

    # Ask user whether to install with Cosmovisor
    def ask_for_cosmovisor(self):
        try:
            logging.info(f"Installing cheqd-node with Cosmovisor allows for automatic unattended upgrades for valid software upgrade proposals. See https://docs.cosmos.network/main/tooling/cosmovisor for more information.\n")
            answer = self.ask(
                f"Install cheqd-noded using Cosmovisor? (yes/no)", default=DEFAULT_USE_COSMOVISOR)
            if answer.lower().startswith("y"):
                self.is_cosmovisor_needed = True
            elif answer.lower().startswith("n"):
                self.is_cosmovisor_needed = False
            else:
                logging.error(
                    f"Invalid input provided during installation. Please choose either 'yes' or 'no'.")
                self.ask_for_cosmovisor()
        except Exception as e:
            logging.exception(
                f"Failed to set whether installation should be done with Cosmovisor. Reason: {e}")

    # Ask user whether to bump Cosmovisor to latest version
    def ask_for_cosmovisor_bump(self):
        try:
            answer = self.ask(
                f"Do you want to bump your Cosmovisor to {DEFAULT_LATEST_COSMOVISOR_VERSION} ? (yes/no)", default=DEFAULT_BUMP_COSMOVISOR)
            if answer.lower().startswith("y"):
                self.is_cosmovisor_bump_needed = True
            elif answer.lower().startswith("n"):
                self.is_cosmovisor_bump_needed = False
            else:
                logging.error(
                    f"Invalid input provided during installation. Please choose either 'yes' or 'no'.")
                self.ask_for_cosmovisor_bump()
        except Exception as e:
            logging.exception(
                f"Failed to set whether Cosmovisor should be bumped to latest version. Reason: {e}")

    # Ask user whether to allow Cosmovisor to automatically download binaries for scheduled upgrades
    def ask_for_daemon_allow_download_binaries(self):
        try:
            answer = self.ask(
                f"Do you want Cosmovisor to automatically download binaries for scheduled upgrades? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.daemon_allow_download_binaries = "true"
            elif answer.lower().startswith("n"):
                self.daemon_allow_download_binaries = "false"
            else:
                logging.error(
                    f"Invalid input provided during installation. Please choose either 'yes' or 'no'.")
                self.ask_for_daemon_allow_download_binaries()
        except Exception as e:
            logging.exception(
                f"Failed to set whether Cosmovisor should automatically download binaries for scheduled upgrades. Reason: {e}")

    # Ask whether Cosmovisor should restart daemon after upgrade
    def ask_for_daemon_restart_after_upgrade(self):
        try:
            answer = self.ask(
                f"Do you want Cosmovisor to automatically restart cheqd-noded service after an upgrade has been applied? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.daemon_restart_after_upgrade = "true"
            elif answer.lower().startswith("n"):
                self.daemon_restart_after_upgrade = "false"
            else:
                logging.error(
                    f"Invalid input provided during installation. Please choose either 'yes' or 'no'.")
                self.ask_for_daemon_restart_after_upgrade()
        except Exception as e:
            logging.exception(
                f"Failed to set whether Cosmovisor should automatically restart cheqd-noded service after an upgrade has been applied. Reason: {e}")

    # Ask user for node moniker
    def ask_for_moniker(self):
        try:
            logging.info(f"Moniker is a human-readable name for your cheqd-node. This is NOT the same as your validator name, and is only used to uniquely identify your node for Tendermint P2P address book. It can be edited later in your ~.cheqdnode/config/config.toml file.\n")
            answer = self.ask(
                f"Provide a moniker for your cheqd-node", default=CHEQD_NODED_MONIKER)
            if answer is not None:
                self.moniker = answer
            else:
                logging.error(
                    f"Invalid moniker provided during cheqd-noded setup.")
        except Exception as e:
            logging.exception(f"Failed to set moniker. Reason: {e}")

    # Ask for node's external IP address or DNS name
    def ask_for_external_address(self):
        try:
            logging.info(f"External address is the publicly accessible IP address or DNS name of your cheqd-node. This is used to advertise your node's P2P address to other nodes in the network. If you are running your node behind a NAT, you should set this to your public IP address or DNS name. If you are running your node on a public IP address, you can leave this blank to automatically fetch your IP address via DNS resolver lookup. This sends a `dig` request to whoami.cloudflare.com\n\n")
            answer = self.ask(
                f"What is the externally-reachable IP address or DNS name for your cheqd-node? [default: Fetch automatically via DNS resolver lookup]: {os.linesep}")
            if answer:
                self.external_address = answer
            else:
                self.external_address = str(self.exec(
                    "dig +short txt ch whoami.cloudflare @1.1.1.1").stdout).strip("""b'"\\n""")
        except Exception as e:
            logging.exception(f"Failed to set external address. Reason: {e}")

    # Ask for node's P2P port
    def ask_for_p2p_port(self):
        try:
            answer = self.ask(f"Specify your node's P2P port",
                              default=DEFAULT_P2P_PORT)
            if answer is not None:
                self.p2p_port = answer
            else:
                self.p2p_port = DEFAULT_P2P_PORT
        except Exception as e:
            logging.exception(f"Failed to set P2P port. Reason: {e}")

    # Ask for node's RPC port
    def ask_for_rpc_port(self):
        try:
            answer = self.ask(f"Specify your node's RPC port",
                              default=DEFAULT_RPC_PORT)
            if answer is not None:
                self.rpc_port = answer
            else:
                self.rpc_port = DEFAULT_RPC_PORT
        except Exception as e:
            logging.exception(f"Failed to set RPC port. Reason: {e}")

    # (Optional) Ask for node's persistent peers
    def ask_for_persistent_peers(self):
        try:
            logging.info(f"Persistent peers are nodes that you want to always keep connected to. Values for persistent peers should be specified in format: <nodeID>@<IP>:<port>,<nodeID>@<IP>:<port>... \n")
            answer = self.ask(
                f"Specify persistent peers [default: none]: {os.linesep}")
            if answer is not None:
                self.persistent_peers = answer
            else:
                self.persistent_peers = ""
        except Exception as e:
            logging.exception(f"Failed to set persistent peers. Reason: {e}")

    # (Optional) Ask for minimum gas prices
    def ask_for_gas_price(self):
        try:
            logging.info(
                f"Minimum gas prices are the minimum amount of CHEQ tokens you are willing to accept as a validator to process a transaction.\n")
            answer = self.ask(f"Specify minimum gas price",
                              default=CHEQD_NODED_MINIMUM_GAS_PRICES)
            if answer is not None:
                self.gas_price = answer
            else:
                self.gas_price = default = CHEQD_NODED_MINIMUM_GAS_PRICES
        except Exception as e:
            logging.exception(f"Failed to set minimum gas prices. Reason: {e}")

    # (Optional) Ask for node's log level
    def ask_for_log_level(self):
        try:
            self.log_level = self.ask(
                f"Specify log level (trace|debug|info|warn|error|fatal|panic)", default=CHEQD_NODED_LOG_LEVEL)
        except Exception as e:
            logging.exception(f"Failed to set log level. Reason: {e}")

    # (Optional) Ask for node's log format
    def ask_for_log_format(self):
        try:
            self.log_format = self.ask(
                f"Specify log format (json|plain)", default=CHEQD_NODED_LOG_FORMAT)
        except Exception as e:
            logging.exception(f"Failed to set log format. Reason: {e}")
    
    # If an existing installation is detected, ask user if they want to upgrade
    def ask_for_upgrade(self):
        try:
            logging.warning(
                f"Existing cheqd-node configuration folder detected.\n")
            answer = self.ask(
                f"Do you want to upgrade an existing cheqd-node installation? (yes/no)", default="no")
            if answer.lower().startswith("y"):
                self.is_upgrade = True
            elif answer.lower().startswith("n"):
                self.is_upgrade = False
            else:
                logging.exception(
                    f"Invalid input provided during installation.")
        except Exception as e:
            logging.exception(f"Failed to upgrade cheqd-node. Reason: {e}")

    # If an install from scratch is requested, warn the user and check if they want to proceed
    def ask_for_install_from_scratch(self):
        try:
            logging.warning(
                f"Doing a fresh installation of cheqd-node will remove ALL existing configuration and data.\nCAUTION: Please ensure you have a backup of your existing configuration and data before proceeding!\n")
            answer = self.ask(
                f"Do you want to do fresh installation of cheqd-node? (yes/no)", default="no")
            if answer.lower().startswith("y"):
                self.is_from_scratch = True
            elif answer.lower().startswith("n"):
                self.is_from_scratch = False
            else:
                logging.exception(
                    f"Invalid input provided during installation.")
        except Exception as e:
            logging.exception(
                f"Failed to set whether cheqd-node should install from scratch. Reason: {e}")

    # If an existing installation is detected, ask user if they want to overwrite existing systemd configuration
    def ask_for_rewrite_node_systemd(self):
        try:
            answer = self.ask(
                f"Overwrite existing systemd configuration for cheqd-node? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.rewrite_node_systemd = True
            elif answer.lower().startswith("n"):
                self.rewrite_node_systemd = False
            else:
                logging.exception(
                    f"Invalid input provided during installation.")
        except Exception as e:
            logging.exception(
                f"Failed to set whether overwrite existing systemd configuration for cheqd-node. Reason: {e}")

    # If an existing installation is detected, ask user if they want to overwrite existing logrotate configuration
    def ask_for_rewrite_logrotate(self):
        try:
            answer = self.ask(
                f"Overwrite existing configuration for logrotate? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.rewrite_logrotate = True
            elif answer.lower().startswith("n"):
                self.rewrite_logrotate = False
            else:
                logging.exception(
                    f"Invalid input provided during installation.")
        except Exception as e:
            logging.exception(
                f"Failed to set whether overwrite existing configuration for logrotate. Reason: {e}")

    # If an existing installation is detected, ask user if they want to overwrite existing rsyslog configuration
    def ask_for_rewrite_rsyslog(self):
        try:
            answer = self.ask(
                f"Overwrite existing configuration for cheqd-node logging? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.rewrite_rsyslog = True
            elif answer.lower().startswith("n"):
                self.rewrite_rsyslog = False
            else:
                logging.exception(
                    f"Invalid input provided during installation.")
        except Exception as e:
            logging.exception(
                f"Failed to set whether overwrite existing rsyslog configuration for cheqd-node. Reason: {e}")

    # Ask user if they want to download a snapshot of the existing chain to speed up node synchronization.
    # This is only applicable if installing from scratch.
    # This question is asked last because it is the most time consuming.
    def ask_for_init_from_snapshot(self):
        try:
            logging.warning(
                f"CAUTION: Downloading a snapshot replaces your existing copy of chain data! Usually safe to use this option when doing a fresh installation.\n")
            answer = self.ask(
                f"Do you want to download a snapshot of the existing chain to speed up node synchronization? (yes/no)", default=DEFAULT_INIT_FROM_SNAPSHOT)
            if answer.lower().startswith("y"):
                self.init_from_snapshot = True
            elif answer.lower().startswith("n"):
                self.init_from_snapshot = False
            else:
                logging.exception(
                    f"Invalid input provided during installation.")
        except Exception as e:
            logging.exception(
                f"Failed to set whether init snapshot. Reason: {e}")


if __name__ == '__main__':
    # Order of questions to ask the user if installing:
    # 1. Version of cheqd-noded to install
    # 2. Home directory for cheqd user
    # 3. Install new version of cheqd-noded
    # 4. Chain ID to join
    # 5. Install Cosmovisor if not installed, or bump Cosmovisor version
    # 6. (if applicable) Cosmovisor settings
    # 7. Node configuration settings
    # 8. Download snapshot to bootsrap node
    def install_steps():
        try:
            interviewer.ask_for_version()
            interviewer.ask_for_home_directory()
            interviewer.ask_for_setup()
            interviewer.ask_for_chain()

            if interviewer.is_cosmovisor_installed is False:
                interviewer.ask_for_cosmovisor()
            else:
                interviewer.ask_for_cosmovisor_bump()

            if interviewer.is_cosmovisor_needed is True:
                interviewer.ask_for_daemon_allow_download_binaries()
                interviewer.ask_for_daemon_restart_after_upgrade()

            if interviewer.is_setup_needed is True:
                interviewer.ask_for_moniker()
                interviewer.ask_for_external_address()
                interviewer.ask_for_p2p_port()
                interviewer.ask_for_rpc_port()
                interviewer.ask_for_persistent_peers()
                interviewer.ask_for_gas_price()
                interviewer.ask_for_log_level()
                interviewer.ask_for_log_format()

            interviewer.ask_for_init_from_snapshot()

        except Exception as e:
            logging.exception(
                f"Unable to complete user interview process for installation. Reason for exiting: {e}")

    # Order of questions to ask the user if installing:
    # 1. Version of cheqd-noded to install
    # 2. Home directory for cheqd user
    # 3. Install Cosmovisor if not installed, or bump Cosmovisor version
    # 4. (if applicable) Cosmovisor settings
    # 6. Rewrite node systemd config
    # 7. Rewrite rsyslog config
    # 8. Rewrite logrotate config
    def upgrade_steps():
        try:
            interviewer.ask_for_version()
            interviewer.ask_for_home_directory()

            if interviewer.is_cosmovisor_installed is False:
                interviewer.ask_for_cosmovisor()
            else:
                interviewer.ask_for_cosmovisor_bump()

            if interviewer.is_cosmovisor_needed is True:
                interviewer.ask_for_daemon_allow_download_binaries()
                interviewer.ask_for_daemon_restart_after_upgrade()

            if interviewer.is_systemd_config_installed(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH) is True or interviewer.is_systemd_config_installed(DEFAULT_STANDALONE_SERVICE_FILE_PATH) is True:
                interviewer.ask_for_rewrite_node_systemd()

            if interviewer.is_systemd_config_installed(DEFAULT_RSYSLOG_FILE) is True:
                interviewer.ask_for_rewrite_rsyslog()

            if interviewer.is_systemd_config_installed(DEFAULT_LOGROTATE_FILE) is True:
                interviewer.ask_for_rewrite_logrotate()

        except Exception as e:
            logging.exception(
                f"Unable to complete user interview process for upgrade. Reason for exiting: {e}")

    ### This section is where the Interviewer class is invoked ###
    try:
        interviewer = Interviewer()

        # Check if cheqd-noded is already installed
        installed = interviewer.is_node_installed()

        # Check if Cosmovisor is already installed
        cosmovisor_installed = interviewer.check_cosmovisor_installed()

        # If no cheqd-noded binary is found, install from scratch
        if installed is False:
            install_steps()

        else:
            # If cheqd-noded binary is found, ask user if they want to upgrade or install from scratch
            interviewer.ask_for_upgrade()

            # If user wants to upgrade, execute upgrade steps
            if interviewer.is_upgrade is True:
                upgrade_steps()

            else:
                # If user declines upgrade, ask if they want to install from scratch
                interviewer.ask_for_install_from_scratch()

                if interviewer.is_from_scratch is True:
                    install_steps()
                else:
                    logging.error("Aborting installation to prevent overwriting existing node installation. Exiting...")
                    raise

    except Exception as e:
        logging.exception(f"Unable to complete user interview process. Reason for exiting: {e}")
        raise

    ### This section where the Installer class is invoked ###
    try:
        installer = Installer(interviewer)
        if installer.install():
            logging.info(f"Installation of cheqd-noded {interviewer.version} completed successfully!")
            sys.exit(0)
        else:
            logging.error(f"Installation of cheqd-noded {interviewer.version} failed. Exiting...")
            raise
    except Exception as e:
        logging.exception(f"Unable to execute installation process. Reason for exiting: {e}")
        sys.exit(1)
