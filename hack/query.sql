SELECT `greptime_timestamp`, `cluster`, `app`, `message`,
  `region`,
  `cloud-provider`,
  `environment`,
  `product`,
  `sub-product`,
  `service`
FROM
  `$cluster`.`$app`
WHERE
  AND (
    "$search" = ""
    OR (message LIKE '%$search%')
  )
  AND $__timeFilter(`greptime_timestamp`)
ORDER BY
  `greptime_timestamp`
LIMIT $limit;
