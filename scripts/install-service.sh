#!/bin/bash

set -e

apt update && apt install sudo vim make gcc git curl rsyslog systemd -y

# create user
sudo useradd ghoul -m
cd /home/ghoul

# Install GO
curl -O https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz
tar -xzf go1.11.4.linux-amd64.tar.gz
sudo mv go /usr/local

# Setting GO Paths
echo 'export GOPATH=$HOME/go' | sudo tee -a /home/ghoul/.profile
echo 'export PATH="$PATH:/usr/local/go/bin:$GOPATH/bin"' | sudo tee -a /home/ghoul/.profile
echo 'export GO111MODULE=on' | sudo tee -a /home/ghoul/.profile

# Clone & build server from source
git clone https://ghoul.git ghoul
source /home/ghoul/.profile
cd ghoul
make build
chown -R ghoul:ghoul .

# Install systemd service
sudo mv ghoul.service /lib/systemd/system/.
sudo chmod 755 /lib/systemd/system/ghoul.service

# Modify "/etc/rsyslog.conf" and uncomment the lines below
#   module(load="imtcp")
#   input(type="imtcp" port="514")
# Then, create “/etc/rsyslog.d/30-ghoul.conf” with the following content:
#   if $programname == 'ghoul' or $syslogtag == 'ghoul' then /var/log/ghoul/ghoul.log
#   & stop
# Now restart the rsyslog service and ghoul service. View logs by:
#   tail -f /var/log/ghoul/ghoul.log

# Install syncer
# */5 * * * * cd /home/ghoul/ghoul-api && ./syncer >> /var/log/ghoul/syncer.log 2>&1
