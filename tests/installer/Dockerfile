###############################################################
###             STAGE 1: Installer test runner              ###
###############################################################

FROM ubuntu:focal AS runner

# Install minimum necessary dependencies
ENV PACKAGES ca-certificates systemd rsyslog
RUN apt-get update && apt-get install $PACKAGES -y

# Add necessary files
COPY . /installer-test
COPY ./../../installer/installer.py /installer-test/installer.py

# Set bash defaults
WORKDIR /installer-test
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Set entrypoint script
ENTRYPOINT [ "/bin/bash", "-c" ]
