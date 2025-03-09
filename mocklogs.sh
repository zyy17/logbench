./bin/mocklogs \
  -start-time 2025-03-08T00:00:00Z \
  -end-time 2025-03-08T00:00:20Z \
  -interval-count 20 \
  -cluster cluster1 \
  -app app1 \
  -interval 1s \
  -format parquet
