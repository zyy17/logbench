# logbench

logbench is a tool to benchmark the performance of log ingestion.

## Usage

```console
./bin/logbench \
  -batch-size 500 \
  -ops 10000 \
  -endpoint http://localhost:4000 \
  -db public \
  -pipeline greptime_identity
```

## Build

- Build the binary:

  ```console
  make
  ```

- Build the local image:

  ```console
  make build-local-test-image
  ```
