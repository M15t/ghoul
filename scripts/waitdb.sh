#!/usr/bin/env bash

# Ensure the database container is online and usable
# echo "Waiting for database..."
until docker exec -i ghoul-api.db mysql -u ghoul -pghoul123 -D ghoul -e "SELECT 1" &> /dev/null
# EnablePostgreSQL: remove the line above, uncomment the following
# until docker exec -i ghoul-api.db psql -h localhost -U ghoul -d ghoul -c "SELECT 1" &> /dev/null
do
  # printf "."
  sleep 1
done
