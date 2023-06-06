# bettercron
A better cron alternative, written in go to learn and practise

## Installation
### Download and decompress
Just run install.sh script located on the root of the repository (still not on point)

For manual installation, download the latest release and decompress the tar.gz file (lines done for amd64 arch, adapt it for your processor architecture)

```bash
curl -l https://github.com/marcsanchezg/bettercron/releases/download/v0.0.3/bettercron_0.0.3_linux_amd64.tar.gz /tmp/bettercron.tar.gz
tar -xvzf /tmp/bettercron.tar.gz
```

### create bettercron user
We're gonna create a user for running bettercron as, it will automatically create a username, with no login and with sudo permissions

```bash
sudo useradd -m -s /usr/sbin/nologin -G sudo -p "$(openssl passwd -1 '')" bettercron
```

### Set up log and config
The next step is to create a config directory and a log file plus adding permissions to it

```bash
# Copy binary
sudo cp "/tmp/bettercron/bettercron" "/usr/bin/bettercron"
sudo chown bettercron:bettercron "/usr/bin/bettercron"
sudo chmod +x "/usr/bin/bettercron"

# Create config file and directory
sudo mkdir -p "/etc/bettercron"
sudo cp "/tmp/bettercron/example/config.yaml" "/etc/bettercron/config.yaml"
sudo chown -R bettercron:bettercron "/etc/bettercron"
```

### Set-up systemd daemon
We're gonna copy bettercron.service and enable it

```bash
# Create systemd service
sudo cp "/tmp/bettercron/example/bettercron.service" "/etc/systemd/system/bettercron.service"

# Enable and start the service
sudo systemctl enable bettercron.service
sudo systemctl start bettercron.service
```

## bettercron flags

```bash
> bettercron --help
Usage of bettercron:
  -config string
    	Define yaml file location (default "/etc/bettercron/config.yaml")
  -help
    	Show help information
  -log string
    	Define log file location (default "/var/log/bettercron.log")

```