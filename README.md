# logbench

logbench is a tool to benchmark the performance of log ingestion.

## Usage

```console
./bin/logbench \
  -batch-size 500 \
  -total-size 5000100 \
  -ops 10000 \
  -endpoint http://localhost:4000 \
  -db public \
  -table testlogs \
  -pipeline greptime_identity
```

## Build

```console
make
```
