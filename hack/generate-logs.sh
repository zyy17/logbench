#! /bin/bash

set -e

YEAR=2025
MONTH=03
START_DAY=2
DAY_COUNT=1
OUTPUT_DIR=.output
#IMAGE=registry.cn-hangzhou.aliyuncs.com/zyyinternal/logbench:20250309-39fde7f
IMAGE=localhost:5001/logbench:latest

for day in $(seq $START_DAY $((START_DAY + DAY_COUNT - 1))); do
  for hour in {0..23}; do
    # Calculate start and end time.
    start_time=$(printf "%04d-%02d-%02dT%02d:00:00Z" $YEAR $MONTH $day $hour)
    
    # Calculate end time.
    if [ $hour -eq 23 ]; then
      # If it's the last hour of the day, add 1 to the day and set the hour to 00.
      next_day=$((day + 1))
      end_time=$(printf "%04d-%02d-%02dT00:00:00Z" $YEAR $MONTH $next_day)
    else
      # Otherwise, the same day, add 1 to the hour.
      end_time=$(printf "%04d-%02d-%02dT%02d:00:00Z" $YEAR $MONTH $day $((hour + 1)))
    fi
    
    echo "Processing day $day, hour $hour: $start_time to $end_time"
    
    docker run -v $(pwd)/$OUTPUT_DIR:/output \
      --rm --entrypoint mocklogs $IMAGE \
      -start-time "$start_time" \
      -end-time "$end_time" \
      -interval-count 80 \
      -cluster cluster1 \
      -app app1 \
      -interval 1s \
      -output-dir /output \
      -format parquet
    
    echo "Completed day $day, hour $hour"
    echo "------------------------"
  done
done
