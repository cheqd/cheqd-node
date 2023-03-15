#!/usr/bin/env python3


###############################################################
###     		    Python package imports      			###
###############################################################
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
import socket
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
DEFAULT_LATEST_COSMOVISOR_VERSION = "v1.3.0"
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
    # If PYTHONDEVMODE = 1, show more detailed logging messages
    logging.basicConfig(format='[%(levelname)s]: %(message)s', level=logging.DEBUG)
    logging.raiseExceptions = True
    logging.propagate = True
else:
    # Else show logging messages INFO level and above
    logging.basicConfig(format='[%(levelname)s]: %(message)s', level=logging.INFO)
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
    try:
        file = open(file_path, "r")
        for line in file:
            line = line.strip()
            if search_text in line:
                with open(file_path, "r") as file:
                    data = file.read()
                    data = data.replace(line, replace_text)
                with open(file_path, "w") as file:
                    file.write(data)
    except Exception as e:
        logging.exception(f"Failed to search and replace text in {file_path}. Reason: {e}")
        raise

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
    def cheqd_home_dir(self):
        # Root directory for cheqd-noded
        # Default: /home/cheqd
        return self.interviewer.home_dir

    @property
    def cheqd_backup_dir(self):
        # Root directory for cheqd-noded
        # Default: /home/cheqd/backup
        return os.path.join(self.cheqd_home_dir, "backup")

    @property
    def cheqd_root_dir(self):
        # Root directory for cheqd-noded
        # Default: /home/cheqd/.cheqdnode
        return os.path.join(self.cheqd_home_dir, ".cheqdnode")

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
    def cosmovisor_binary_path(self):
        # Path where Cosmovisor binary will be installed
        # Default: /usr/bin/cosmovisor
        return os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_COSMOVISOR_BINARY_NAME)

    @property
    def temporary_cosmovisor_binary_path(self):
        # Temporary path for Cosmovisor binary just after it's downloaded
        # This is NOT the final install path
        return os.path.join(os.path.realpath(os.path.curdir), DEFAULT_COSMOVISOR_BINARY_NAME)

    @property
    def standalone_node_binary_path(self):
        # Path where cheqd-noded binary will be installed
        # Default: /usr/bin/cheqd-noded
        # When installing with Cosmovisor, this will be a symlink to Cosmovisor directory
        return os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME)

    @property
    def temporary_node_binary_path(self):
        # Temporary path for cheqd-node binary just after it's downloaded
        # This is NOT the final install path
        return os.path.join(os.path.realpath(os.path.curdir), DEFAULT_BINARY_NAME)

    @property
    def cosmovisor_current_bin_path(self):
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
            _url = COSMOVISOR_BINARY_URL.format(DEFAULT_LATEST_COSMOVISOR_VERSION, DEFAULT_LATEST_COSMOVISOR_VERSION, os_arch)
            if is_valid_url(_url):
                logging.debug(f"Cosmovisor download URL: {_url}")
                return _url
            else:
                logging.exception(f"Cosmovisor download URL is not valid: {_url}")
        except Exception as e:
            logging.exception(f"Failed to compute Cosmovisor download URL. Reason: {e}")

    @property
    def cosmovisor_service_cfg(self):
        # Modify cheqd-cosmovisor.service template file to replace values for environment variables
        # The template file is fetched from the GitHub repo
        # Some of these variables are explicitly asked during the installer process. Others are set to default values.
        try:
            # Set service file path
            fname = os.path.basename(COSMOVISOR_SERVICE_TEMPLATE)

            # Fetch the template file from GitHub
            if is_valid_url(COSMOVISOR_SERVICE_TEMPLATE):
                with request.urlopen(COSMOVISOR_SERVICE_TEMPLATE) as response, open(fname, "w") as file:
                    file.write(response.read())

                    # Replace the values for environment variables in the template file
                    s = re.sub(
                        r'({CHEQD_ROOT_DIR}|{DEFAULT_BINARY_NAME}|{COSMOVISOR_DAEMON_ALLOW_DOWNLOAD_BINARIES}|{COSMOVISOR_DAEMON_RESTART_AFTER_UPGRADE}|{DEFAULT_DAEMON_POLL_INTERVAL}|{DEFAULT_UNSAFE_SKIP_BACKUP}|{DEFAULT_DAEMON_RESTART_DELAY})',
                        lambda m: {'{CHEQD_ROOT_DIR}': self.cheqd_root_dir,
                                '{DEFAULT_BINARY_NAME}': DEFAULT_BINARY_NAME,
                                '{COSMOVISOR_DAEMON_ALLOW_DOWNLOAD_BINARIES}':  self.interviewer.daemon_allow_download_binaries,
                                '{COSMOVISOR_DAEMON_RESTART_AFTER_UPGRADE}': self.interviewer.daemon_restart_after_upgrade,
                                '{DEFAULT_DAEMON_POLL_INTERVAL}': DEFAULT_DAEMON_POLL_INTERVAL,
                                '{DEFAULT_UNSAFE_SKIP_BACKUP}': DEFAULT_UNSAFE_SKIP_BACKUP,
                                '{DEFAULT_DAEMON_RESTART_DELAY}': DEFAULT_DAEMON_RESTART_DELAY}[m.group()],
                        file.read()
                    )
                
                # Remove the template file
                self.remove_safe(fname)
                return s
            else:
                logging.exception(f"URL is not valid: {RSYSLOG_TEMPLATE}")
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
                with request.urlopen(RSYSLOG_TEMPLATE) as response, open(fname, "w") as file:
                    file.write(response.read())
                
                    # Replace the values for environment variables in the template file
                    s = re.sub(
                        r'({BINARY_FOR_LOGGING}|{CHEQD_LOG_DIR})',
                        lambda m: {'{BINARY_FOR_LOGGING}': binary_name,
                                    '{CHEQD_LOG_DIR}': self.cheqd_log_dir}[m.group()],
                        file.read()
                    )

                    # Remove the template file
                    self.remove_safe(fname)
                    return s
            else:
                logging.exception(f"URL is not valid: {RSYSLOG_TEMPLATE}")
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
                with request.urlopen(LOGROTATE_TEMPLATE) as response, open(fname, "w") as file:
                    file.write(response.read())

                    # Replace the values for environment variables in the template file
                    s = re.sub(
                        r'({CHEQD_LOG_DIR})',
                        lambda m: {'{CHEQD_LOG_DIR}': self.cheqd_log_dir}[m.group()],
                        file.read()
                    )

                # Remove the template file
                self.remove_safe(fname)
                return s
            else:
                logging.exception(f"URL is not valid: {LOGROTATE_TEMPLATE}")
        except Exception as e:
            logging.exception(f"Failed to set up logrotate from template. Reason: {e}")

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
                return False

            # Create cheqd user if it doesn't exist
            if self.prepare_cheqd_user():
                logging.info("User/group cheqd setup successfully")
            else:
                logging.error("Failed to setup user/group cheqd")
                return False
            
            # Setup directories needed for installation
            if self.prepare_directory_tree():
                logging.info("Directory tree setup successfully")
            else:
                logging.error("Failed to setup directory tree")
                return False

            # Carry out pre-installation steps
            # Mostly relevant if installing from scratch or re-installing
            if self.pre_install():
                logging.info("Pre-installation steps completed successfully")
            else:
                logging.error("Failed to complete pre-installation steps")
                return False

            # Setup Cosmovisor binary if needed
            if self.interviewer.is_cosmovisor_needed or self.interviewer.is_cosmovisor_bump_needed:
                if self.install_cosmovisor():
                    logging.info("Successfully installed Cosmovisor")
                else:
                    logging.error("Failed to setup Cosmovisor")
                    return False
            # If Cosmovisor is not needed, treat it as a standalone installation
            else:
                if self.install_standalone():
                    logging.info("Successfully installed cheqd-noded as a standalone binary")
                else:
                    logging.error("Failed to setup cheqd-noded as a standalone binary")
                    return False
            
            # Setup cheqd-noded environment variables
            # These are independent of Cosmovisor environment variables
            # Set them regardless of whether Cosmovisor is used or not
            self.set_cheqd_env_vars()
            
            # Configure cheqd-noded settings
            # This edits the config.toml and app.toml files
            if self.configure_node_settings():
                logging.info("Successfully configured cheqd-noded settings")
            else:
                logging.error("Failed to configure cheqd-noded settings")
                return False

            # Configure systemd service for cheqd-noded
            # Sets up either a standalone service or a Cosmovisor service
            # ONLY enables it without activating it
            if self.setup_node_systemd():
                logging.info("Successfully configured systemd service for node operations")
            else:
                logging.error("Failed to configure systemd service for node operations")
                return False

            # Configure systemd services for rsyslog and logrotate
            if self.setup_logging_systemd():
                logging.info("Successfully configured systemd service for logging")
            else:
                logging.error("Failed to configure systemd service for logging")
                return False

            # Download and extract snapshot if needed
            if self.interviewer.init_from_snapshot:
                # Check if snapshot download was successful
                if self.download_snapshot():
                    logging.info("Successfully downloaded snapshot")
                else:
                    logging.error("Failed to download snapshot")
                    return False
                                
                if self.extract_snapshot():
                    logging.info("Successfully extracted snapshot")
                else:
                    logging.error("Failed to extract snapshot")
                    return False
            else:
                logging.debug("Skipping snapshot download and extraction as it was not requested")

            # Return True if all steps were successful
            logging.info("Installation steps completed successfully")
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
            with request.urlopen(binary_url) as response, open(fname, "wb") as file:
                file.write(response.read())
            
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
        except Exception as e:
            logging.exception("Failed to download cheqd-noded binary. Reason: {e}")
            return False

    def pre_install(self) -> bool:
        # Pre-installation steps
        # Removes the following existing cheqd-noded data and configurations:
        # 1. ~/.cheqdnode directory
        # 2. cheqd-noded / cosmovisor binaries
        # 3. systemd service files
        try:
            # Stop existing systemd services first if running
            # Check if the service is running before stopping it
            self.stop_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME)
            self.stop_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME)

            # Create backup directory if it doesn't exist
            os.makedirs(self.cheqd_backup_dir, exist_ok=True)

            # Make a copy of validator key and state before removing user data
            # Use shutil.copytree() when copying directories
            # Use shutil.copy() when copying files instead of shutil.copyfile() since it preserves file metadata
            logging.info("Backing up user's config folder and selected validator secrets from data folder")

            if os.path.exists(self.cheqd_config_dir):
                # Backup ~/.cheqdnode/config/ folder
                shutil.copytree(self.cheqd_config_dir, self.cheqd_backup_dir)
            else:
                logging.debug("No config folder found to backup. Skipping...")

            # Backup ~/.cheqdnode/data/priv_validator_key.json
            # Without this file, a validator node will get jailed!
            if os.path.exists(os.path.join(self.cheqd_data_dir, "priv_validator_key.json")):
                shutil.copy(os.path.join(self.cheqd_data_dir, "priv_validator_state.json"), 
                    os.path.join(self.cheqd_backup_dir, "priv_validator_state.json"))
            else:
                logging.debug("No validator state file found to backup. Skipping...")
            
            # Backup ~/.cheqdnode/data/upgrade-info.json
            # This file is required for Cosmovisor to track and understand where upgrade is needed
            if os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")):
                shutil.copyfile(os.path.join(self.cheqd_data_dir, "upgrade-info.json"), 
                    os.path.join(self.cheqd_backup_dir, "upgrade-info.json"))
            else:
                logging.debug("No upgrade-info.json file found to backup. Skipping...")

            # Change ownership of backup directory to cheqd user
            shutil.chown(self.cheqd_backup_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)

            if self.interviewer.is_from_scratch or self.interviewer.is_setup_needed:
                # Remove cheqd-node data and binaries
                logging.warning("Removing user's data and configs")
                self.remove_safe(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_COSMOVISOR_BINARY_NAME))
                self.remove_safe(os.path.join(DEFAULT_INSTALL_PATH, DEFAULT_BINARY_NAME))
                self.remove_safe(self.cheqd_root_dir, is_dir=True)
                return True
            else:
                logging.debug("No pre-installation steps needed")
                return True
        except Exception as e:
            logging.exception(f"Could not complete pre-installation steps. Reason: {e}")
            return False

    def prepare_cheqd_user(self) -> bool:
        # Create "cheqd" user/group if it doesn't exist
        try:
            if not self.does_user_exist(DEFAULT_CHEQD_USER):
                logging.info(f"Creating {DEFAULT_CHEQD_USER} group")
                self.exec(f"addgroup {DEFAULT_CHEQD_USER} --quiet --system")

                logging.info(f"Creating {DEFAULT_CHEQD_USER} user and adding to {DEFAULT_CHEQD_USER} group")
                self.exec(
                    f"adduser --system {DEFAULT_CHEQD_USER} --home {self.cheqd_home_dir} --shell /bin/bash --ingroup {DEFAULT_CHEQD_USER} --quiet")
                
                # Set permissions for cheqd home directory to cheqd:cheqd
                logging.info(f"Setting permissions for {self.cheqd_home_dir} to {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER}")
                shutil.chown(self.cheqd_home_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)
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

    def prepare_directory_tree(self) -> bool:
        # Needed only in case of clean installation
        # 1. Create ~/.cheqdnode directory
        # 2. Set directory permissions to default cheqd user
        # 3. Create ~/.cheqdnode/log directory
        try:
            # Create root directory for cheqd-noded
            if not os.path.exists(self.cheqd_root_dir):
                logging.info("Creating main directory for cheqd-noded")
                os.makedirs(self.cheqd_root_dir, exist_ok=True)
            else:
                logging.info(f"Skipping main directory creation because {self.cheqd_root_dir} already exists")

            # Setup logging related directories
            # 1. Create log directory if it doesn't exist
            # 2. Create default stdout.log file in log directory
            # 3. Set ownership of log directory to syslog:cheqd
            if not os.path.exists(self.cheqd_log_dir):
                # Create ~/.cheqdnode/log directory
                logging.info("Creating log directory for cheqd-noded")
                os.makedirs(self.cheqd_log_dir, exist_ok=True)

                # Create blank ~/.cheqdnode/log/stdout.log file. Overwrite if it already exists.
                # Using the .open() method without doing anything in it will create the file
                # "w" mode is used to overwrite the file if it already exists
                with open(os.path.join(self.cheqd_log_dir, "stdout.log"), "w") as file:
                    logging.debug("Created blank stdout.log file")

                logging.info(f"Setting up ownership permissions for {self.cheqd_log_dir} directory")
                shutil.chown(self.cheqd_log_dir, "syslog", DEFAULT_CHEQD_USER)
            else:
                logging.info(f"Skipping log directory creation because {self.cheqd_log_dir} already exists")

            # Create symlink from cheqd-noded log folder from /var/log/cheqd-node
            # This step is necessary since many logging tools look for logs in /var/log
            if not os.path.exists("/var/log/cheqd-node"):
                logging.info("Creating a symlink from cheqd-noded log folder to /var/log/cheqd-node")
                os.symlink(self.cheqd_log_dir, "/var/log/cheqd-node", target_is_directory=True)
            else:
                logging.info("Skipping linking because /var/log/cheqd-node already exists")
            
            # Return True if all steps were successful
            return True
        except Exception as e:
            logging.exception(f"Failed to prepare directory tree for {DEFAULT_CHEQD_USER}. Reason: {e}")
            return False

    def install_cosmovisor(self) -> bool:
        # Install binaries for cheqd-noded and Cosmovisor
        # Cosmovisor is only installed if requested by the user
        # cheqd-noded binary is installed in Cosmovisor bin path under this scenario
        try:
            logging.info("Setting up Cosmovisor...")
            
            # Download Cosmovisor binary and set environment variables
            if self.get_cosmovisor():
                logging.info("Successfully downloaded Cosmovisor")

                # Set environment variables for Cosmovisor
                self.set_cosmovisor_env_vars()
                logging.info("Successfully set Cosmovisor environment variables")
            else:
                logging.error("Failed to download Cosmovisor")
                return False

            # Move Cosmovisor binary to installation directory if it doesn't exist or bump needed
            # This is executed is there is no Cosmovisor binary in the installation directory
            # or if the user has requested a bump for Cosmovisor
            # shutil.move() will overwrite the file if it already exists
            logging.info(f"Moving Cosmovisor {self.temporary_cosmovisor_binary_path} to {self.cosmovisor_binary_path}")
            shutil.move(self.temporary_cosmovisor_binary_path, self.cosmovisor_binary_path)

            # Set ownership of Cosmovisor binary to root:root
            shutil.chown(self.cosmovisor_binary_path, "root", "root")

            # Move cheqd-noded binary to /usr/bin
            logging.info(f"Copying cheqd-noded binary from {self.temporary_node_binary_path} to {self.standalone_node_binary_path}")
            shutil.copy(self.temporary_node_binary_path, self.standalone_node_binary_path)

            # Set ownership of cheqd-noded binary to root:root
            shutil.chown(self.standalone_node_binary_path, "root", "root")

            # Initialize Cosmovisor if it's not already initialized
            # This is done by checking whether the Cosmovisor root directory exists
            if not os.path.exists(self.cosmovisor_root_dir):
                self.exec(f"sudo -u {DEFAULT_CHEQD_HOME_DIR} bash -c 'cosmovisor init {self.standalone_node_binary_path}'")

                # Set ownership of cheqd root directory to cheqd:cheqd
                # This is necessary because the command above may site ownership of the directory to root:root
                shutil.chown(self.cheqd_root_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)
            else:
                logging.info("Cosmovisor directory already exists. Skipping initialisation...")
            
            # Remove cheqd-noded binary from /usr/bin if it's not a symlink
            if not os.path.islink(self.standalone_node_binary_path):
                logging.warn(f"Removing {DEFAULT_BINARY_NAME} from {DEFAULT_INSTALL_PATH} because it is not a symlink")
                os.remove(self.standalone_node_binary_path)

                # Move cheqd-noded binary to Cosmovisor bin path
                # shutil.move() will overwrite the file if it already exists
                logging.info(f"Moving cheqd-noded binary from {self.temporary_node_binary_path} to {self.cosmovisor_current_bin_path}")
                shutil.move(self.temporary_node_binary_path, self.cosmovisor_current_bin_path)

                # Set ownership of cheqd-noded binary to cheqd:cheqd
                # This is ONLY done when the binary is moved to Cosmovisor bin path
                shutil.chown(self.cosmovisor_current_bin_path, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)

                # Create symlink to cheqd-noded binary in Cosmovisor bin path
                # Target comes first, then the location of the symlink
                logging.info(f"Creating symlink to {self.cosmovisor_current_bin_path}")
                os.symlink(self.cosmovisor_current_bin_path, self.standalone_node_binary_path)
            else:
                logging.info(f"{self.cosmovisor_current_bin_path} is already symlink. Skipping removal...")

            # Steps to execute only if this is an upgrade
            # The upgrade-info.json file is required for Cosmovisor to track upgrades
            if self.interviewer.is_upgrade \
                and os.path.exists(os.path.join(self.cheqd_data_dir, "upgrade-info.json")) \
                and not os.path.exists(os.path.join(self.cosmovisor_root_dir, "current/upgrade-info.json")):
                logging.info(f"Copying ~/.cheqdnode/data/upgrade-info.json file to ~/.cheqdnode/cosmovisor/current/")

                # shutil.copy() preserves the file metadata
                shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                    os.path.join(self.cosmovisor_root_dir, "current/upgrade-info.json"), follow_symlinks=True)
            else:
                logging.debug("Skipped copying upgrade-info.json file because it doesn't exist")
            
            # Change owner of Cosmovisor directory to cheqd:cheqd
            logging.info(f"Changing ownership of {self.cosmovisor_root_dir} to {DEFAULT_CHEQD_USER} user")
            shutil.chown(self.cosmovisor_root_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)

            # Return True if all steps were successful
            return True
        except Exception as e:
            logging.exception(f"Failed to setup Cosmovisor. Reason: {e}")
            return False

    def get_cosmovisor(self) -> bool:
        # Download Cosmovisor binary and extract it
        # Also remove the downloaded archive file, if applicable
        try:
            logging.info("Downloading Cosmovisor binary...")
            binary_url = self.cosmovisor_download_url
            fname = os.path.basename(binary_url)

            # Download Cosmovisor binary from GitHub
            with request.urlopen(binary_url) as response, open(fname, "wb") as file:
                file.write(response.read())

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

                # Return True if all steps were successful
                return True
            else:
                logging.error(f"Unable to extract Cosmovisor binary from archive file: {fname}")
                return False
        except Exception as e:
            logging.exception("Failed to download Cosmovisor binary. Reason: {e}")

    def install_standalone(self) -> bool:
        # Install cheqd-noded as a standalone binary
        # cheqd-noded binary is installed in /usr/bin under this scenario
        try:
            logging.info("Setting up standalone cheqd-noded binary...")
            
            # Remove symlink for cheqd-noded if it exists
            if os.path.islink(self.standalone_node_binary_path):
                logging.warn(f"Removing symlink {self.standalone_node_binary_path}")
                os.remove(self.standalone_node_binary_path)
            else:
                logging.info(f"{self.standalone_node_binary_path} is not a symlink. Skipping removal...")

            # Move cheqd-noded binary to /usr/bin
            # shutil.move() will overwrite the file if it already exists
            logging.info(f"Moving cheqd-noded binary from {self.temporary_node_binary_path} to {self.standalone_node_binary_path}")
            shutil.move(self.temporary_node_binary_path, self.standalone_node_binary_path)
            
            # Set ownership of cheqd-noded binary to root:root
            logging.info(f"Changing ownership of {self.standalone_node_binary_path} to root:root")
            shutil.chown(self.standalone_node_binary_path, "root", "root")

            # Remove Cosmovisor directory if it exists
            if os.path.exists(self.cosmovisor_root_dir):
                logging.warn(f"Removing Cosmovisor directory from {self.cosmovisor_root_dir} because it is not required for a standalone installation")
                self.remove_safe(self.cosmovisor_root_dir, is_dir=True)
            else:
                logging.debug(f"{self.cosmovisor_root_dir} doesn't exist. Skipping removal...")

            # Return True if all steps were successful
            return True
        except Exception as e:
            logging.exception(f"Failed to setup Cosmovisor. Reason: {e}")
            return False
    
    def set_cosmovisor_env_vars(self):
        # Set environment variables for Cosmovisor
        try:
            self.set_environment_variable("DAEMON_NAME", DEFAULT_BINARY_NAME)
            self.set_environment_variable("DAEMON_HOME", self.cheqd_root_dir, overwrite=False)
            self.set_environment_variable("DAEMON_ALLOW_DOWNLOAD_BINARIES", 
                self.interviewer.daemon_allow_download_binaries)
            self.set_environment_variable("DAEMON_RESTART_AFTER_UPGRADE",
                self.interviewer.daemon_restart_after_upgrade)
            self.set_environment_variable("DAEMON_POLL_INTERVAL", 
                DEFAULT_DAEMON_POLL_INTERVAL, overwrite=False)
            self.set_environment_variable("UNSAFE_SKIP_BACKUP", 
                DEFAULT_UNSAFE_SKIP_BACKUP, overwrite=False)
        except Exception as e:
            logging.exception(f"Failed to set environment variables for Cosmovisor. Reason: {e}")
            raise
    
    def set_cheqd_env_vars(self):
        # Set environment variables for cheqd-noded binary
        # Applicable for both standalone and Cosmovisor installations
        # Only environment variables that are required required for transactions are set here
        try:
            self.set_environment_variable("CHEQD_NODED_NODE", 
                f"tcp://localhost:{self.interviewer.rpc_port}", overwrite=False)
            if self.interviewer.chain == "testnet":
                self.set_environment_variable("CHEQD_NODED_CHAIN_ID", TESTNET_CHAIN_ID)
            else:
                self.set_environment_variable("CHEQD_NODED_CHAIN_ID", MAINNET_CHAIN_ID)
        except Exception as e:
            logging.exception(f"Failed to set environment variables for cheqd-noded. Reason: {e}")
            raise

    def set_environment_variable(self, env_var_name, env_var_value, overwrite=True):
        # Set an environment variable
        # By default, existing environment variables are overwritten
        # This can be changed by setting the overwrite parameter to False
        # Environment variables are set for the current session as well as for all users
        try:
            logging.debug(f"Checking whether {env_var_name} is set")

            if os.getenv(env_var_name) is None or overwrite:
                logging.debug(f"Setting {env_var_name} to {env_var_value}")
                
                # Set the environment variable for the current session
                os.environ[env_var_name] = env_var_value

                # Modify the system's environment variables
                # This will set the variable permanently for all users
                with open("/etc/environment", "a") as env_file:
                    env_file.write(f"export {env_var_name}={env_var_value}")
                
                # Reload the environment variables
                os.system("source /etc/environment")
            else:
                logging.debug(f"Environment variable {env_var_name} already set or overwrite is disabled")
        except Exception as e:
            logging.exception(f"Failed to set environment variable {env_var_name}. Reason: {e}")
            raise

    def configure_node_settings(self) -> bool:
        # Configure cheqd-noded settings in app.toml and config.toml
        # Some of these need to be set based on user input for setup needed from scratch only
        # Others are needed regardless of whether the node is being setup from scratch or an upgrade path
        try:
            # Set file paths for common configuration files
            app_toml_path = os.path.join(self.cheqd_config_dir, "app.toml")
            config_toml_path = os.path.join(self.cheqd_config_dir, "config.toml")
            genesis_file_path = os.path.join(self.cheqd_config_dir, 'genesis.json')

            # Set URLs for files to be downloaded
            genesis_url = GENESIS_FILE.format(self.interviewer.chain)
            seeds_url = SEEDS_FILE.format(self.interviewer.chain)

            # These changes are required only when NEW node setup is needed
            if self.interviewer.is_setup_needed:
                # Don't execute an init in case a validator key already exists
                if not os.path.exists(os.path.join(self.cheqd_config_dir, 'priv_validator_key.json')):
                    # Initialize the node
                    logging.info(f"Initialising {self.cheqd_root_dir} directory")
                    self.exec(f"""sudo -u {DEFAULT_CHEQD_USER} -c 'cheqd-noded init {self.interviewer.moniker}'""")
                else:
                    logging.debug(f"Validator key already exists in {self.cheqd_config_dir}. Skipping cheqd-noded init...")
                
                # Check if genesis file exists
                # If not, download it from the GitHub repo
                if is_valid_url(genesis_url) and not os.path.exists(genesis_file_path):
                    logging.debug(f"Downloading genesis file for {self.interviewer.chain}")
                    
                    with request.urlopen(genesis_url) as response, open(genesis_file_path, "w") as file:
                        file.write(response.read())
                else:
                    logging.debug(f"Genesis file already exists in {genesis_file_path}")

                # Set seeds from the seeds file on GitHub
                if is_valid_url(seeds_url):
                    logging.debug(f"Setting seeds from {seeds_url}")
                    
                    with request.urlopen(seeds_url) as response:
                        seeds = response.read().decode("utf-8").strip()
                    
                    seeds_search_text = 'seeds = ""'
                    seeds_replace_text = 'seeds = "{}"'.format(seeds)
                    search_and_replace(seeds_search_text, seeds_replace_text, config_toml_path)
                else:
                    logging.exception(f"Invalid URL for seeds file: {seeds_url}")
                    return False

                # Set RPC port to listen to for all origins by default
                rpc_default_value = 'laddr = "tcp://127.0.0.1:{}"'.format(DEFAULT_RPC_PORT)
                new_rpc_default_value = 'laddr = "tcp://0.0.0.0:{}"'.format(DEFAULT_RPC_PORT)
                search_and_replace(rpc_default_value, new_rpc_default_value, config_toml_path)
            else:
                logging.debug("Skipping cheqd-noded init as setup is not needed")

            ### This next section changes values in configuration files only if the user has provided input ###

            # Set external address
            if self.interviewer.external_address:
                external_address_search_text = 'external_address'
                external_address_replace_text = 'external_address = "{}:{}"'.format(
                    self.interviewer.external_address, self.interviewer.p2p_port)
                logging.debug(f"Setting external address to {external_address_replace_text}")
                search_and_replace(external_address_search_text, external_address_replace_text, config_toml_path)
            else:
                logging.debug("External address not set by user. Skipping...")

            # Set P2P port
            if self.interviewer.p2p_port:
                p2p_laddr_search_text = 'laddr = "tcp://0.0.0.0:{}"'.format(DEFAULT_P2P_PORT)
                p2p_laddr_replace_text = 'laddr = "tcp://0.0.0.0:{}"'.format(self.interviewer.p2p_port)
                search_and_replace(p2p_laddr_search_text, p2p_laddr_replace_text, config_toml_path)
            else:
                logging.debug("P2P port not set by user. Skipping...")

            # Setting up the RPC port
            if self.interviewer.rpc_port:
                rpc_laddr_search_text = 'laddr = "tcp://0.0.0.0:{}"'.format(DEFAULT_RPC_PORT)
                rpc_laddr_replace_text = 'laddr = "tcp://0.0.0.0:{}"'.format(self.interviewer.rpc_port)
                search_and_replace(rpc_laddr_search_text, rpc_laddr_replace_text, config_toml_path)
            else:
                logging.debug("RPC port not set by user. Skipping...")

            # Setting up min gas-price
            if self.interviewer.gas_price:
                min_gas_price_search_text = 'minimum-gas-prices'
                min_gas_price_replace_text = 'minimum-gas-prices = "{}"'.format(self.interviewer.gas_price)
                search_and_replace(min_gas_price_search_text, min_gas_price_replace_text, app_toml_path)
            else:
                logging.debug("Minimum gas price not set by user. Skipping...")

            # Setting up persistent peers
            if self.interviewer.persistent_peers:
                persistent_peers_search_text = 'persistent_peers'
                persistent_peers_replace_text = 'persistent_peers = "{}"'.format(self.interviewer.persistent_peers)
                search_and_replace(persistent_peers_search_text, persistent_peers_replace_text, config_toml_path)
            else:
                logging.debug("Persistent peers not set by user. Skipping...")

            # Setting up log level
            if self.interviewer.log_level:
                log_level_search_text = 'log_level'
                log_level_replace_text = 'log_level = "{}"'.format(self.interviewer.log_level)
                search_and_replace(log_level_search_text, log_level_replace_text, config_toml_path)
            else:
                logging.debug("Log level not set by user. Skipping...")

            # Setting up log format
            if self.interviewer.log_format:
                log_format_search_text = 'log_format'
                log_format_replace_text = 'log_format = "{}"'.format(self.interviewer.log_format)
                search_and_replace(log_format_search_text, log_format_replace_text, config_toml_path)
            else:
                logging.debug("Log format not set by user. Skipping...")
            
            # Set ownership of configuration directory to cheqd:cheqd
            logging.info(f"Setting ownership of {self.cheqd_config_dir} to {DEFAULT_CHEQD_USER}:{DEFAULT_CHEQD_USER}")
            shutil.chown(self.cheqd_config_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)

            # Return True if all the above steps were successful
            return True
        except Exception as e:
            logging.exception(f"Failed to configure cheqd-noded settings. Reason: {e}")
            return False

    def setup_node_systemd(self) -> bool:
        # Setup cheqd-noded related systemd services
        # If user selected Cosmovisor install, then cheqd-cosmovisor.service will be setup
        # If user selected Standalone install, then cheqd-noded.service will be setup
        # WARNING: Services should already have been stopped in pre_install() but if it's removed from there,
        # then it should be added here
        try:
            # Remove cheqd-noded.service and cheqd-cosmovisor.service if they exist
            # Also run if setup is from scratch/first-time install
            if self.interviewer.rewrite_node_systemd:
                logging.warning("Removing existing node-related systemd configuration as requested")
                self.remove_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME, DEFAULT_COSMOVISOR_SERVICE_FILE_PATH)
                self.remove_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME, DEFAULT_STANDALONE_SERVICE_FILE_PATH)
            else:
                logging.debug("Node-related systemd configurations don't need to be removed. Skipping...")
            
            # Setup cheqd-cosmovisor.service if requested
            if self.interviewer.is_cosmovisor_needed:
                # Write cheqd-cosmovisor.service file
                # Replace placeholder values with actuals
                with open(DEFAULT_COSMOVISOR_SERVICE_FILE_PATH, "w") as fname:
                    fname.write(self.cosmovisor_service_cfg)

                # Enable cheqd-cosmovisor.service
                self.enable_systemd_service(DEFAULT_COSMOVISOR_SERVICE_NAME)
                return True
            
            # Otherwise, setup cheqd-noded.service for standalone install
            else:
                # Fetch the template file from GitHub
                if is_valid_url(STANDALONE_SERVICE_TEMPLATE):
                    with request.urlopen(STANDALONE_SERVICE_TEMPLATE) as response, open(DEFAULT_STANDALONE_SERVICE_FILE_PATH, "w") as file:
                        file.write(response.read())
                    
                    # Enable cheqd-noded.service
                    self.enable_systemd_service(DEFAULT_STANDALONE_SERVICE_NAME)
                    return True
                else:
                    logging.error(f"Invalid URL provided for standalone service template: {STANDALONE_SERVICE_TEMPLATE}")
                    return False
        except Exception as e:
            logging.exception(f"Failed to setup systemd service for cheqd-node. Reason: {e}")
            return False

    # Setup logging related systemd services
    def setup_logging_systemd(self) -> bool:
        # Install cheqd-node configuration for rsyslog if user wants to rewrite rsyslog service file
        # Also run if setup is from scratch/first-time install
        try:
            if self.interviewer.rewrite_rsyslog:
                # Remove existing rsyslog service file if it exists
                if os.path.exists(DEFAULT_RSYSLOG_FILE):
                    logging.warning("Removing existing rsyslog configuration as requested")
                    self.remove_safe(DEFAULT_RSYSLOG_FILE)
                else:
                    logging.debug("Rsyslog configuration doesn't need to be removed. Skipping...")
                
                # Determine the binary name for logging based on installation type
                if self.interviewer.is_cosmovisor_needed:
                    binary_name = DEFAULT_COSMOVISOR_BINARY_NAME
                else:
                    binary_name = DEFAULT_BINARY_NAME

                logging.info(f"Configuring rsyslog systemd service for {binary_name} logging")

                # Modify rsyslog template file with values specific to the installation
                with open(DEFAULT_RSYSLOG_FILE, "w") as fname:
                    fname.write(self.rsyslog_cfg)

                # Restarting rsyslog can take a lot of time: https://github.com/rsyslog/rsyslog/issues/3133
                if self.restart_systemd_service("rsyslog.service"):
                    logging.info("Successfully configured rsyslog service")
                else:
                    logging.exception("Failed to configure rsyslog service")
                    return False

            # Install cheqd-node configuration for logrotate if user wants to rewrite logrotate service file
            # Also run if setup is from scratch/first-time install
            if self.interviewer.rewrite_logrotate:
                # Remove existing logrotate service file if it exists
                if os.path.exists(DEFAULT_LOGROTATE_FILE):
                    logging.warning("Removing existing logrotate configuration as requested")
                    self.remove_safe(DEFAULT_LOGROTATE_FILE)
                else:
                    logging.debug("Logrotate configuration doesn't need to be removed. Skipping...")

                logging.info(f"Configuring logrotate systemd service for cheqd-node logging")

                # Modify logrotate template file with values specific to the installation
                with open(DEFAULT_LOGROTATE_FILE, "w") as fname:
                    fname.write(self.logrotate_cfg)

                # Restart logrotate.service
                if self.restart_systemd_service("logrotate.service"):
                    logging.info("Successfully configured logrotate service")
                else:
                    logging.exception("Failed to configure logrotate service")
                    return False
                
                # Restart logrotate.timer
                if self.restart_systemd_service("logrotate.timer"):
                    logging.info("Successfully configured logrotate timer")
                else:
                    logging.exception("Failed to configure logrotate timer")
                    return False
            
            # Return True if both rsyslog and logrotate services are configured
            return True
        except Exception as e:
            logging.exception(f"Failed to setup logging systemd services. Reason: {e}")
            return False

    def download_snapshot(self) -> bool:
        # Download snapshot archive if requested by the user
        # This is a blocking operation that will take a while
        try:
            # Only proceed if a valid snapshot URL has been set
            if self.set_snapshot_url():
                logging.info(f"Valid snapshot URL found: {self.snapshot_url}")
                fname = os.path.basename(self.snapshot_url)
                file_path = os.path.join(self.cheqd_root_dir, fname)
            else:
                logging.error(f"No valid snapshot URL found in last {MAX_SNAPSHOT_DAYS} days!")
                return False
            
            # Install dependencies needed to show progress bar
            if self.install_dependencies():
                logging.info("Dependencies required for snapshot restore installed successfully")
            else:
                logging.error("Failed to install dependencies required for snapshot restore")
                return False

            # Fetch size of snapshot archive WITHOUT downloading it
            req = request.Request(self.snapshot_url, method='HEAD')
            response = request.urlopen(req)
            content_length = response.getheader("Content-Length")
            if content_length is not None:
                archive_size = content_length
                logging.debug(f"Snapshot archive size: {content_length} bytes")
            else:
                logging.error(f"Could not determine snapshot archive size")
                return False

            # Free up some disk space by deleting contents of the data folder
            # Otherwise, there may not be enough space to download AND extract the snapshot
            # WARNING: Backup the priv_validator_state.json and upgrade-info.json before doing this!
            if os.path.exists(self.cheqd_backup_dir):
                # Check that backup of validator keys, state, and upgrade info exists before proceeding
                logging.info(f"Backup directory exists: {self.cheqd_backup_dir}")

                # Remove contents of data directory
                logging.warning(f"Contents of {self.cheqd_data_dir} will be deleted to make room for snapshot")
                self.remove_safe(self.cheqd_data_dir, is_dir=True)

                # Recreate data directory
                os.makedirs(self.cheqd_data_dir, exist_ok=True)
                shutil.chown(self.cheqd_data_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)
            else:
                logging.warning(f"Backup directory does not exist. Will not delete data directory.\n")
                logging.warning(f"Free disk space will be calculated without freeing up space.\n")

            # Check how much free disk space is available wherever the cheqd home directory is located
            # First, determine where the home directory is mounted
            fs_stats = os.statvfs(self.cheqd_home_dir)

            # Calculate the free space in bytes
            free_space = fs_stats.f_frsize * fs_stats.f_bavail

            # ONLY download the snapshot if there is enough free disk space
            if int(archive_size) < int(free_space):
                logging.info(f"Downloading snapshot and extracting archive. This can take a *really* long time...")
                
                # Use wget to download since it can show a progress bar while downloading natively
                # This is a blocking operation that will take a while
                # "wget -c" will resume a download if it gets interrupted
                self.exec(f"wget -c {self.snapshot_url} -P {self.cheqd_home_dir}")

                if self.compare_checksum(file_path, self.snapshot_url):
                    logging.info(f"Snapshot download was successful AND checksums match.")
                    return True
                else:
                    logging.error(f"Snapshot download was successful BUT checksums do not match.")
                    logging.warning(f"Removing corrupted snapshot archive: {file_path}")
                    self.remove_safe(file_path)
                    return False
            else:
                logging.error(f"Snapshot is larger than free disk space. Please free up disk space and try again.")
                return False
        except Exception as e:
            logging.exception(f"Failed to download snapshot. Reason: {e}")
            return False

    def set_snapshot_url(self) -> bool:
        # Get latest available snapshot URL from snapshots.cheqd.net for the given chain
        # This checks whether there are any snapshots in past MAX_SNAPSHOT_DAYS (default: 7 days)
        try:
            template = TESTNET_SNAPSHOT if self.interviewer.chain in TESTNET_CHAIN_ID else MAINNET_SNAPSHOT
            snapshot_date = datetime.date.today()
            counter = 0
            valid_url_found = False

            # Iterate over past MAX_SNAPSHOT_DAYS days to find the latest snapshot
            while not valid_url_found and counter <= MAX_SNAPSHOT_DAYS:
                _url = template.format(snapshot_date.strftime(
                    "%Y-%m-%d"), snapshot_date.strftime("%Y-%m-%d"))
                valid_url_found = is_valid_url(_url)
                counter += 1
                snapshot_date -= datetime.timedelta(days=1)

            # Set snapshot URL if found
            if valid_url_found:
                self.snapshot_url = _url
                logging.debug(f"Snapshot URL: {self.snapshot_url}")
                return True
            else:
                logging.debug("Could not find a valid snapshot in last {} days".format(MAX_SNAPSHOT_DAYS))
                return False
        except Exception as e:
            logging.exception(f"Failed to get snapshot URL. Reason: {e}")
            return False

    def install_dependencies(self) -> bool:
        # Install dependencies required for snapshot extraction
        try:
            # Update apt lists before installing dependencies
            logging.info("Updating apt lists")
            self.exec("sudo apt-get update")

            # Use apt-get to install dependencies
            logging.info(f"Install pv to show progress of extraction")
            self.exec("sudo apt-get install -y pv")
            return True
        except Exception as e:
            logging.exception(f"Failed to install dependencies. Reason: {e}")
            return False

    def compare_checksum(self, file_path, snapshot_url) -> bool:
        # Compare checksum of downloaded snapshot with published checksum
        # This is to ensure that the snapshot was downloaded correctly
        try:
            # Split snapshot URL into its components
            url_parts = list(os.path.split(snapshot_url))

            # Replace archive filename with checksum filename
            url_parts[-1] = "md5sum.txt"

            # Construct checksum file URL
            checksum_url = os.path.join("", *url_parts)

            # Fetch published checksum from checksum URL if it exists
            if is_valid_url(checksum_url):
                with request.urlopen(checksum_url) as response:
                    published_checksum = response.read().decode("utf-8").strip()
                
                # Calculate checksum of downloaded snapshot
                # This is a blocking operation that will take a while
                # Python's hashlib.md5() is not used because making it work with large files is a pain
                local_checksum = self.exec(f"md5sum {file_path} | tail -1 | cut -d' ' -f 1").stdout.strip()
                
                # Print checksums for debugging
                logging.debug(f"Published checksum: {published_checksum}")
                logging.debug(f"Local checksum: {local_checksum}")

                # Compare checksums
                if published_checksum == local_checksum:
                    logging.debug(f"Checksums match. Download is OK.")
                    return True
                else:
                    logging.debug(f"Checksums do not match. Download got corrupted.")
                    return False
            else:
                logging.error(f"Checksum URL is invalid. File integrity couldn't be tested.")
                return False
        except Exception as e:
            logging.exception(f"Failed to compare checksums. Reason: {e}")
            return False

    def extract_snapshot(self):
        # Extract snapshot archive to cheqd node data directory
        # This is a blocking operation that will take a while
        # Once extracted, restore files from backup folder
        try:
            # Set file path of snapshot archive
            file_path = os.path.join(self.cheqd_root_dir, os.path.basename(self.snapshot_url))

            # Extract to cheqd node data directory EXCEPT for validator state
            # Snapshot archives are created using lz4 compression since it's more efficient than gzip
            if os.path.exists(file_path):
                logging.info(f"Extracting snapshot archive. This may take a while...")

                # Bash command is used since the Python libraries for lz4 are not installed out-of-the-box
                # Showing a progress bar or an estimate of time remaining is also not easy-to-achieve
                # "pv" is used to show a progress bar while extracting
                self.exec(f"sudo -u {DEFAULT_CHEQD_USER} -c 'pv {file_path} \
                    | tar --use-compress-program=lz4 -xf - -C {self.cheqd_root_dir} \
                    --exclude priv_validator_state.json'")

                # Delete snapshot archive file
                logging.info(f"Snapshot extraction was successful. Deleting snapshot archive.")
                self.remove_safe(file_path)
            else:
                logging.error(f"Snapshot archive file not found. Could not extract snapshot.")
                return False
            
            # Restore files from backup folder
            # Use shutil.copy() instead of shutil.copyfile() to preserve file metadata
            if os.path.exists(self.cheqd_backup_dir):
                logging.info(f"Restoring files from backup folder.")
                
                # Restore priv_validator_state.json
                logging.info(f"Restoring priv_validator_state.json to {self.cheqd_data_dir}")
                shutil.copy(os.path.join(self.cheqd_backup_dir, "priv_validator_state.json"),
                            os.path.join(self.cheqd_data_dir, "priv_validator_state.json"))
                
                # Restore upgrade-info.json
                logging.info(f"Restoring upgrade-info.json to {self.cheqd_data_dir}")
                shutil.copy(os.path.join(self.cheqd_backup_dir, "upgrade-info.json"),
                            os.path.join(self.cheqd_data_dir, "upgrade-info.json"))
                
                # If Cosmovisor is needed, copy upgrade-info.json to ~/.cheqdnode/cosmovisor/current/ directory
                # Otherwise, Cosmovisor will throw an error
                if self.interviewer.is_cosmovisor_needed:
                    logging.info(f"Restoring upgrade-info.json to {self.cosmovisor_root_dir}/current/")
                    shutil.copy(os.path.join(self.cheqd_data_dir, "upgrade-info.json"),
                                os.path.join(self.cosmovisor_root_dir, "current/upgrade-info.json"))

                    # Change ownership of Cosmovisor directory to cheqd user
                    shutil.chown(self.cosmovisor_root_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)
            else:
                logging.warning(f"Backup folder not found. Please check and restore the following folders/files from your own backup:\n~/.cheqdnode/data/priv_validator_state.json, ~/.cheqdnode/data/upgrade-info.json, ~/.cheqdnode/config/")

            # Change ownership of cheqd node data directory to cheqd user
            shutil.chown(self.cheqd_data_dir, DEFAULT_CHEQD_USER, DEFAULT_CHEQD_USER)

            # Return True if snapshot extraction was successful
            return True
        except Exception as e:
            logging.exception(f"Failed to extract snapshot. Reason: {e}")
            return False

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
        except Exception as e:
            logging.exception(f"Error daemon reloading: Reason: {e}")
    
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
                else:
                    logging.debug(f"{service_name} is already enabled")
                    return True
            else:
                logging.error(f"Failed to reload systemd config and reset failed services")
                return False
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
            else:
                logging.error(f"Failed to restart {service_name}")
                return False
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
        self._rpc_port = ""
        self._p2p_port = ""
        self._gas_price = ""
        self._persistent_peers = ""
        self._log_level = ""
        self._log_format = ""
        self._daemon_allow_download_binaries = DEFAULT_DAEMON_ALLOW_DOWNLOAD_BINARIES
        self._daemon_restart_after_upgrade = DEFAULT_DAEMON_RESTART_AFTER_UPGRADE
        self._is_from_scratch = False
        self._rewrite_node_systemd = True
        self._rewrite_rsyslog = True
        self._rewrite_logrotate = True

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
    def is_cosmovisor_installed(self, ici):
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
            logging.exception(f"Could not check if cheqd-noded is already installed. Reason: {e}")

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

    # Check whether external address provided is valid IP address
    def check_ip_address(self, ip_address) -> bool:
        try:
            socket.inet_aton(ip_address)
            logging.debug(f"IP address {ip_address} is valid")
            return True
        except socket.error:
            logging.debug(f"IP address {ip_address} is invalid")
            return False

    # Check whether external address provided is valid DNS name
    def check_dns_name(self, dns_name) -> bool:
        try:
            socket.gethostbyname(dns_name)
            logging.debug(f"DNS name {dns_name} is valid")
            return True
        except socket.error:
            logging.debug(f"DNS name {dns_name} is invalid")
            return False

    # Check if a systemd config is installed for a given service file
    def is_systemd_config_installed(self, systemd_service_file) -> bool:
        try:
            if os.path.exists(systemd_service_file):
                return True
            else:
                return False
        except Exception as e:
            logging.exception(f"Could not check if {systemd_service_file} already exists. Reason: {e}")

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
            logging.exception(f"Could not get releases from GitHub. Reason: {e}")

    # The "latest" stable release may not be in last N releases, so we need to get it separately
    def get_latest_release(self):
        try:
            req = request.Request(
                "https://api.github.com/repos/cheqd/cheqd-node/releases/latest")
            req.add_header("Accept", "application/vnd.github.v3+json")
            with request.urlopen(req) as response:
                return Release(json.loads(response.read().decode("utf-8")))
        except Exception as e:
            logging.exception(f"Could not get latest release from GitHub. Reason: {e}")

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
            logging.exception(f"Could not assemble list of releases to show to the user. Reason: {e}")

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

            # Check that user selected a valid release number
            if release_num >= 1 and release_num <= LAST_N_RELEASES and isinstance(release_num, int):
                self.release = all_releases[release_num - 1]
                logging.debug(f"Release version selection: {self.release.version}")
            else:
                logging.error(f"Invalid release number picked from list of releases: {release_num}")
                logging.error(f"Please choose a number between 1 and {LAST_N_RELEASES}\n")
                self.ask_for_version()
        except Exception as e:
            logging.exception(f"Failed to selected version of cheqd-noded. Reason: {e}")

    # Set cheqd user's home directory
    def ask_for_home_directory(self) -> str:
        try:
            self.home_dir = self.ask(
                f"Set path for cheqd user's home directory", default=DEFAULT_CHEQD_HOME_DIR)
            logging.debug(f"Setting home directory to {self.home_dir}")
        except Exception as e:
            logging.exception(f"Failed to set cheqd user's home directory. Reason: {e}")

    # Ask whether user wants to do a install from scratch
    def ask_for_setup(self):
        try:
            answer = self.ask(
                f"Do you want to setup a new cheqd-node installation? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.is_setup_needed = True
            elif answer.lower().startswith("n") and self.is_node_installed():
                self.is_setup_needed = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'\n")
                self.ask_for_setup()
        except Exception as e:
            logging.exception(f"Failed to set fresh installation parameters. Reason: {e}")

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
                logging.error(f"Invalid network selected during installation. Please choose either 1 or 2.\n")
                self.ask_for_chain()
            
            # Set debug message
            logging.debug(f"Setting network to join as {self.chain}")
        except Exception as e:
            logging.exception(f"Failed to set network/chain to join. Reason: {e}")

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
                logging.error(f"Invalid input provided during installation. Please choose either 'yes' or 'no'.\n")
                self.ask_for_cosmovisor()
        except Exception as e:
            logging.exception(f"Failed to set whether installation should be done with Cosmovisor. Reason: {e}")

    # Ask user whether to bump Cosmovisor to latest version
    def ask_for_cosmovisor_bump(self):
        try:
            answer = self.ask(
                f"Do you want to bump your Cosmovisor to {DEFAULT_LATEST_COSMOVISOR_VERSION}? (yes/no)", default=DEFAULT_BUMP_COSMOVISOR)
            if answer.lower().startswith("y"):
                self.is_cosmovisor_bump_needed = True
            elif answer.lower().startswith("n"):
                self.is_cosmovisor_bump_needed = False
            else:
                logging.error(f"Invalid input provided during installation. Please choose either 'yes' or 'no'.\n")
                self.ask_for_cosmovisor_bump()
        except Exception as e:
            logging.exception(f"Failed to set whether Cosmovisor should be bumped to latest version. Reason: {e}")

    # Ask user whether to allow Cosmovisor to automatically download binaries for scheduled upgrades
    def ask_for_daemon_allow_download_binaries(self):
        try:
            answer = self.ask(
                f"Do you want Cosmovisor to automatically download binaries for scheduled upgrades? (yes/no)", 
                default="yes")
            if answer.lower().startswith("y"):
                self.daemon_allow_download_binaries = "true"
            elif answer.lower().startswith("n"):
                self.daemon_allow_download_binaries = "false"
            else:
                logging.error(f"Invalid input provided during installation. Please choose either 'yes' or 'no'.\n")
                self.ask_for_daemon_allow_download_binaries()
        except Exception as e:
            logging.exception(
                f"Failed to set whether Cosmovisor should automatically download binaries. Reason: {e}")

    # Ask whether Cosmovisor should restart daemon after upgrade
    def ask_for_daemon_restart_after_upgrade(self):
        try:
            answer = self.ask(
                f"Do you want Cosmovisor to automatically restart after an upgrade? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.daemon_restart_after_upgrade = "true"
            elif answer.lower().startswith("n"):
                self.daemon_restart_after_upgrade = "false"
            else:
                logging.error(f"Invalid input provided during installation. Please choose either 'yes' or 'no'.\n")
                self.ask_for_daemon_restart_after_upgrade()
        except Exception as e:
            logging.exception(f"Failed to set whether Cosmovisor should automatically restart after an upgrade. Reason: {e}")

    # Ask user for node moniker
    def ask_for_moniker(self):
        try:
            logging.info(f"Moniker is a human-readable name for your cheqd-node.\nThis is NOT the same as your validator name, and is only used to uniquely identify your node for Tendermint P2P address book.\nIt can be edited later in your ~/.cheqdnode/config/config.toml file.\n")
            self.moniker = self.ask(
                f"Provide a moniker for your cheqd-node", default=CHEQD_NODED_MONIKER)
            if self.moniker is not None and isinstance(self.moniker, str):
                logging.debug(f"Moniker set to {self.moniker}")
            else:
                logging.error(f"Invalid moniker provided during cheqd-noded setup.\n")
                self.ask_for_moniker()
        except Exception as e:
            logging.exception(f"Failed to set moniker. Reason: {e}")

    # Ask for node's external IP address or DNS name
    def ask_for_external_address(self):
        try:
            logging.info(f"External address is the publicly accessible IP address or DNS name of your cheqd-node.\nThis is used to advertise your node's P2P address to other nodes in the network.\n- If you are running your node behind a NAT, you should set this to your public IP address or DNS name\n- If you are running your node on a public IP address, you can leave this blank to automatically fetch your IP address via DNS resolver lookup.\n- Automatic fetching sends a `dig` request to whoami.cloudflare.com\n")
            
            answer = self.ask(
                f"What is the externally-reachable IP address or DNS name for your cheqd-node? [default: Fetch automatically via DNS resolver lookup]: {os.linesep}")
            
            # If user provided an answer, check if it's a valid IP address or DNS name
            if answer:
                if self.check_ip_address(answer) or self.check_dns_name(answer):
                    self.external_address = answer
                else:
                    logging.error(f"Invalid IP address or DNS name provided. Please enter a valid IP address or DNS name.\n")
                    self.ask_for_external_address()
            # If user didn't provide an answer, fetch IP address via DNS resolver lookup
            else:
                self.external_address = str(self.exec(
                    "dig +short txt ch whoami.cloudflare @1.1.1.1").stdout).strip("""b'"\\n""")

            logging.debug(f"External address set to {self.external_address}")
        except Exception as e:
            logging.exception(f"Failed to set external address. Reason: {e}")

    # Ask for node's P2P port
    def ask_for_p2p_port(self):
        try:
            self.p2p_port = int(self.ask(f"Specify your node's P2P port", default=DEFAULT_P2P_PORT))
            if isinstance(self.p2p_port, int):
                logging.debug(f"P2P port set to {self.p2p_port}")
            else:
                logging.error(f"Invalid P2P port provided. Please enter a valid port number.\n")
                self.ask_for_p2p_port()
        except Exception as e:
            logging.exception(f"Failed to set P2P port. Reason: {e}")

    # Ask for node's RPC port
    def ask_for_rpc_port(self):
        try:
            self.rpc_port = int(self.ask(f"Specify your node's RPC port", default=DEFAULT_RPC_PORT))
            if isinstance(self.rpc_port, int):
                logging.debug(f"RPC port set to {self.rpc_port}")
            else:
                logging.error(f"Invalid RPC port provided. Please enter a valid port number.\n")
                self.ask_for_rpc_port()
        except Exception as e:
            logging.exception(f"Failed to set RPC port. Reason: {e}")

    # (Optional) Ask for node's persistent peers
    def ask_for_persistent_peers(self):
        try:
            logging.info(f"Persistent peers are nodes that you want to always keep connected to. Values for persistent peers should be specified in format: <nodeID>@<IP>:<port>,<nodeID>@<IP>:<port>...\n")
            answer = self.ask(
                f"Specify persistent peers [default: none]: {os.linesep}")
            if answer is not None:
                self.persistent_peers = answer
                logging.debug(f"Persistent peers set to {self.persistent_peers}")
            else:
                self.persistent_peers = ""
                logging.debug(f"No persistent peers set.")
        except Exception as e:
            logging.exception(f"Failed to set persistent peers. Reason: {e}")

    # (Optional) Ask for minimum gas prices
    def ask_for_gas_price(self):
        try:
            logging.info(
                f"Minimum gas prices is the price you are willing to accept as a validator to process a transaction.\nValues should be entered in format <number>ncheq (e.g., 50ncheq)\n")
            self.gas_price = self.ask(f"Specify minimum gas price", default=CHEQD_NODED_MINIMUM_GAS_PRICES)
            if self.gas_price.endswith("ncheq"):
                logging.debug(f"Minimum gas price set to {self.gas_price}")
            else:
                logging.error(f"Invalid minimum gas price provided. Valid format is <number>ncheq.\n")
                self.ask_for_gas_price()
        except Exception as e:
            logging.exception(f"Failed to set minimum gas prices. Reason: {e}")

    # (Optional) Ask for node's log level
    def ask_for_log_level(self):
        try:
            self.log_level = self.ask(
                f"Specify log level (trace|debug|info|warn|error|fatal|panic)", default=CHEQD_NODED_LOG_LEVEL)
            if self.log_level in ["trace", "debug", "info", "warn", "error", "fatal", "panic"]:
                logging.debug(f"Log level set to {self.log_level}")
            else:
                logging.error(f"Invalid log level provided. Please enter a valid log level.\n")
                self.ask_for_log_level()
        except Exception as e:
            logging.exception(f"Failed to set log level. Reason: {e}")

    # (Optional) Ask for node's log format
    def ask_for_log_format(self):
        try:
            self.log_format = self.ask(f"Specify log format (json|plain)", default=CHEQD_NODED_LOG_FORMAT)
            if self.log_format in ["json", "plain"]:
                logging.debug(f"Log format set to {self.log_format}")
            else:
                logging.error(f"Invalid log format provided. Please enter a valid log format.\n")
                self.ask_for_log_format()
        except Exception as e:
            logging.exception(f"Failed to set log format. Reason: {e}")
    
    # If an existing installation is detected, ask user if they want to upgrade
    def ask_for_upgrade(self):
        try:
            logging.warning(f"Existing cheqd-node configuration folder detected.\n")
            answer = self.ask(f"Do you want to upgrade an existing cheqd-node installation? (yes/no)", default="no")
            if answer.lower().startswith("y"):
                self.is_upgrade = True
            elif answer.lower().startswith("n"):
                self.is_upgrade = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'\n")
                self.ask_for_upgrade()
        except Exception as e:
            logging.exception(f"Failed to set whether installation should be upgraded. Reason: {e}")

    # If an install from scratch is requested, warn the user and check if they want to proceed
    def ask_for_install_from_scratch(self):
        try:
            logging.warning(f"Doing a fresh installation of cheqd-node will remove ALL existing configuration and data.\nPlease ensure you have a backup of your existing configuration and data before proceeding!\n")
            answer = self.ask(
                f"Do you want to do fresh installation of cheqd-node? (yes/no)", default="no")
            if answer.lower().startswith("y"):
                self.is_from_scratch = True
            elif answer.lower().startswith("n"):
                self.is_from_scratch = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'\n")
                self.ask_for_install_from_scratch()
        except Exception as e:
            logging.exception(f"Failed to set whether to install from scratch. Reason: {e}")

    # If an existing installation is detected, ask user if they want to overwrite existing systemd configuration
    def ask_for_rewrite_node_systemd(self):
        try:
            answer = self.ask(
                f"Overwrite existing systemd configuration for node-related services? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.rewrite_node_systemd = True
            elif answer.lower().startswith("n"):
                self.rewrite_node_systemd = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'\n")
                self.ask_for_rewrite_node_systemd()
        except Exception as e:
            logging.exception(f"Failed to set whether overwrite existing systemd configuration. Reason: {e}")

    # If an existing installation is detected, ask user if they want to overwrite existing logrotate configuration
    def ask_for_rewrite_logrotate(self):
        try:
            answer = self.ask(f"Overwrite existing configuration for logrotate? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.rewrite_logrotate = True
            elif answer.lower().startswith("n"):
                self.rewrite_logrotate = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'\n")
                self.ask_for_rewrite_logrotate()
        except Exception as e:
            logging.exception(f"Failed to set whether overwrite existing configuration for logrotate. Reason: {e}")

    # If an existing installation is detected, ask user if they want to overwrite existing rsyslog configuration
    def ask_for_rewrite_rsyslog(self):
        try:
            answer = self.ask(f"Overwrite existing configuration for cheqd-node logging? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.rewrite_rsyslog = True
            elif answer.lower().startswith("n"):
                self.rewrite_rsyslog = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'\n")
                self.ask_for_rewrite_rsyslog()
        except Exception as e:
            logging.exception(f"Failed to set whether overwrite existing rsyslog configuration. Reason: {e}")

    # Ask user if they want to download a snapshot of the existing chain to speed up node synchronization.
    # This is only applicable if installing from scratch.
    # This question is asked last because it is the most time consuming.
    def ask_for_init_from_snapshot(self):
        try:
            logging.info(f"Downloading a snapshot allows you to get a copy of the blockchain data to speed up node bootstrapping\nSnapshots can be 100 GBs so downloading can take a really long time!\nExisting chain data folder will be replaced! Usually safe to use this option when doing a fresh installation.\n")
            answer = self.ask(
                f"Do you want to download a snapshot of the existing chain to speed up node synchronization? (yes/no)", default="yes")
            if answer.lower().startswith("y"):
                self.init_from_snapshot = True
            elif answer.lower().startswith("n"):
                self.init_from_snapshot = False
            else:
                logging.error(f"Please choose either 'yes' or 'no'\n")
                self.ask_for_init_from_snapshot()
        except Exception as e:
            logging.exception(f"Failed to set whether init snapshot. Reason: {e}")


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
            logging.exception(f"Unable to complete user interview process for installation. Reason for exiting: {e}")

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
            logging.exception(f"Unable to complete user interview process for upgrade. Reason for exiting: {e}")

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
                    sys.exit(1)

    except Exception as e:
        logging.exception(f"Unable to complete user interview process. Reason for exiting: {e}")
        raise

    ### This section where the Installer class is invoked ###
    try:
        installer = Installer(interviewer)
        if installer.install():
            logging.info(f"Installation of cheqd-noded {installer.version} completed successfully!\n")
            logging.info(f"Please review the configuration files manually and use systemctl to start the node.\n")
            logging.info(f"Documentation: https://docs.cheqd.io/node\n")
            sys.exit(0)
        else:
            logging.error(f"Installation of cheqd-noded {installer.version} failed. Exiting...")
            logging.info(f"Documentation: https://docs.cheqd.io/node\n")
            sys.exit(1)

    except Exception as e:
        logging.exception(f"Unable to execute installation process. Reason for exiting: {e}")
        sys.exit(1)
