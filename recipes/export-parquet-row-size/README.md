# Exporting to Parquet based on the number of rows

In this recipe, we'll learn how to export data into Parquet with a limit on the number of rows per Parquet file.

## Download Clickhouse

```bash
curl https://clickhouse.com/ | sh
```

## Setup ClickHouse

Launch ClickHouse Local

```bash
./clickhouse local -m --path trips.chdb
```

## Exporting to Parquet

Using 10 million generated numbers, we can partition those values by dividing `rowNumberInAllBlocks()` by the number of rows we want per file.

```sql
WITH numbers AS (
  FROM system.numbers 
  select * 
  LIMIT 10_000_000
), 1_000_000 AS rowsPerFile
FROM numbers 
SELECT intDiv(rowNumberInAllBlocks(), rowsPerFile) AS partitionId, count() 
GROUP BY partitionId;
```

We can export that to Parquet files using the following query

```sql
insert into function 
file('data/{_partition_id}.parquet') 
partition by intDiv(rowNumberInAllBlocks(), 1_000_000)
select * 
FROM system.numbers 
LIMIT 10_000_000;
```