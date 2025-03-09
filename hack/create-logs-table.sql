CREATE TABLE IF NOT EXISTS `public`.`logsbench` (
  `greptime_timestamp` TimestampNanosecond NOT NULL TIME INDEX,
  `app` STRING NULL INVERTED INDEX,
  `cluster` STRING NULL INVERTED INDEX,
  `message` STRING NULL,
  `region` STRING NULL,
  `cloud-provider` STRING NULL,
  `environment` STRING NULL,
  `product` STRING NULL,
  `sub-product` STRING NULL,
  `service` STRING NULL
) WITH (
  append_mode = 'true'
);
