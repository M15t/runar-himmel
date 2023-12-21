#!/usr/bin/env bash

# Ensure the database container is online and usable
# echo "Waiting for database..."
until docker compose exec db sh -c 'mysql -u root -p$MYSQL_ROOT_PASSWORD -D $MYSQL_DATABASE -e "SELECT 1"' &> /dev/null

do
	# printf "."
	sleep 1
done

#* EnableMySQL: replace the `until ...` line with the one below
# until docker compose exec db sh -c 'mysql -u root -p$MYSQL_ROOT_PASSWORD -D $MYSQL_DATABASE -e "SELECT 1"' &> /dev/null
#* EnablePostgreSQL: replace the `until ...` line with the one below
# until docker compose exec db sh -c 'psql -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1"' &> /dev/null
