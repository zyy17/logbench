#! /usr/bin/env bash

set -e

DB_HOST=127.0.0.1
DB_PORT=4002
ClusterNum=4
AppNum=5

# Create databases.
for ((i=1; i<=ClusterNum; i++)); do
    echo "Create database cluster$i"
    mysql -h $DB_HOST -P $DB_PORT -e "CREATE DATABASE cluster$i;"
done

# Create tables.
for ((i=1; i<=ClusterNum; i++)); do
    for ((j=1; j<=AppNum; j++)); do
        echo "Create table cluster$i.app$j"
        mysql -h $DB_HOST -P $DB_PORT -e "CREATE TABLE IF NOT EXISTS cluster$i.app$j (
          \`greptime_timestamp\` TIMESTAMP(9) NOT NULL TIME INDEX,
          \`app\` STRING NULL INVERTED INDEX,
          \`cluster\` STRING NULL INVERTED INDEX,
          \`message\` STRING NULL,
          \`region\` STRING NULL,
          \`cloud-provider\` STRING NULL,
          \`environment\` STRING NULL,
          \`product\` STRING NULL,
          \`sub-product\` STRING NULL,
          \`service\` STRING NULL
          ) WITH (
            append_mode = 'true'
          );"
    done
done
