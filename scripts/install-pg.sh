#!/bin/bash

set -e

sudo add-apt-repository "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -sc)-pgdg main"
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

sudo apt-get update
sudo apt-get install postgresql-9.6 postgresql-contrib-9.6 -y

## Set password for default user `postgres`
# Enter this cmd, then enter:
#   \password postgres
# and follow instructions for setting postgres admin password. Press Ctrl+D or type \\q to quit psql terminal
sudo -u postgres psql postgres

sudo -u postgres createuser -D -A -P ghoul
sudo -u postgres createdb -O ghoul ghoul

# TODO: Update pg_hba.conf and postgresql.conf
