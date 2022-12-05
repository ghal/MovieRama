#!/bin/sh
set -e

migrate_db() {
  # Attempt to run migrations, retry if the database container is not yet ready.
  echo "Running migrations on ${MYSQL_DB}"
  i=0
  ./migrate -database mysql://${MYSQL_USERNAME}:${MYSQL_PASSWORD}@tcp\(${MYSQL_WRITE}:${MYSQL_PORT}\)/${MYSQL_DB} -verbose -source file://./migration up
  while [ $? -ne 0 -a $i -lt 20 ]; do
    echo "Database not ready (attempt #$i), retrying.."
    sleep 3
    i=$(expr $i + 1)
    ./migrate -database mysql://${MYSQL_USERNAME}:${MYSQL_PASSWORD}@tcp\(${MYSQL_WRITE}:${MYSQL_PORT}\)/${MYSQL_DB} -verbose -source file://./migration up
  done
}

check_success() {
  # Exit if the last command failed.
  if [ $? -ne 0 ]; then
    echo "Last command failed, exiting.."
    exit 1
  fi
}

set +e
# Run default migration.
migrate_db
check_success

echo "Migration succeeded, starting movierama.."

set -e
../movierama
