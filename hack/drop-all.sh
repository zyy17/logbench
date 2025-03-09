#! /usr/bin/env bash

set -e

DB_HOST=127.0.0.1
DB_PORT=4002
ClusterNum=40
AppNum=50

# Drop all databases.
for ((i=1; i<=ClusterNum; i++)); do
    echo "Drop database cluster$i"
    mysql -h $DB_HOST -P $DB_PORT -e "DROP DATABASE cluster$i;"
done
