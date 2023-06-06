#!/bin/bash
# Script to install bettercron as a systemd daemon

# Set variables
ARCH=$(uname -m)
REPO="marcsanchezg/bettercron" # Replace with the actual owner and repo name
ARCHIVE_NAME="bettercron_Linux_$ARCH.tar.gz" # Replace with the actual archive name
INSTALL_DIR="/usr/bin" # Directory to install the binary
CONFIG_DIR="/etc/bettercron" # Directory to install the config file
SERVICE_DIR="/etc/systemd/system" # Directory to install the systemd service

# Create user
sudo useradd -m -s /usr/sbin/nologin -G sudo -p "$(openssl passwd -1 '')" bettercron


if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  # In a git repo

  TEMPORARY_DIR="/tmp/bettercron" # Directory to store temporary exporter tar.gz
  # Download and extract the latest release
  curl -L "https://github.com/$REPO/releases/latest/download/$ARCHIVE_NAME" -o "$ARCHIVE_NAME"
  sudo mkdir -p "$TEMPORARY_DIR"
  sudo tar -xvzf $ARCHIVE_NAME -C /tmp/bettercron
else
  # Not in a git repo
  
  TEMPORARY_DIR=$(dirname "$0") # Directory to store temporary exporter tar.gz
fi

# Copy binary
sudo cp "$TEMPORARY_DIR/bettercron" "$INSTALL_DIR/bettercron"
sudo chown bettercron:bettercron "$INSTALL_DIR/bettercron"
sudo chmod +x "$INSTALL_DIR/bettercron"

# Create config file
sudo mkdir -p "$CONFIG_DIR"
sudo cp  "$TEMPORARY_DIR/config.yaml" "$CONFIG_DIR/config.yaml"
sudo chown bettercron:bettercron "$CONFIG_DIR/config.yaml"

# Create systemd service
sudo cp -r "$TEMPORARY_DIR/bettercron.service" "$SERVICE_DIR/bettercron.service"
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable bettercron.service
sudo systemctl start bettercron.service

# Cleanup remaining files
#sudo rm -rf "$TEMPORARY_DIR"
#sudo rm "$ARCHIVE_NAME"
