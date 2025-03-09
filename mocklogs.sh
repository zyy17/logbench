./bin/mocklogs \
  -start-time 2025-03-08T00:00:00Z \
  -end-time 2025-03-08T00:01:00Z \
  -interval-count 80 \
  -cluster cluster1 \
  -app app1 \
  -interval 1s \
  -format parquet
