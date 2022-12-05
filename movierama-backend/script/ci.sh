#!/bin/sh
set -e

migrate_db() {
  # Attempt to run migrations, retry if the database container is not yet ready.
  echo "Running migrations on ${MYSQL_DB}"
  done=0
  i=0
  while :; do
    ./migrate -database mysql://${MYSQL_USERNAME}:${MYSQL_PASSWORD}@tcp\(${MYSQL_WRITE}:${MYSQL_PORT}\)/${MYSQL_DB} -verbose -source file://./migration up
    if [ $? -eq 0 ]; then
      break
    fi
    if [ ${i} -eq 20 ]; then
      echo "Error: Database connection failed."
      exit 1
    fi
    sleep 2
    echo "Database not ready (attempt #$i), retrying.."
    i=$(expr $i + 1)
  done
}

set +e
# Run default migration.
migrate_db
set -e
echo "Migration succeeded!"

# lint
echo "Checking lint"
golint -set_exit_status=1 $(go list -mod=vendor ./...)
echo "Lint success!"

# test
echo "Running tests"
go test $(go list ./... | grep -v 'docs') -mod=vendor -race -cover -tags=integration -coverprofile=coverage.txt -covermode=atomic
echo "Test execution completed!"
