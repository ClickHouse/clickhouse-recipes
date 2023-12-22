# SQL Dynamic Column Selection with ClickHouse

In this recipe, we'll learn how to do dynamic column selection with ClickHouse.

## Download NYC Taxis Dataset

We're going to use the NYC Taxis Dataset. 
Download at least one of the Yellow Taxi Trip Records Parquet files from https://www.nyc.gov/site/tlc/about/tlc-trip-record-data.page.

For example, January 2023

```bash
curl https://d37ci6vzurychx.cloudfront.net/trip-data/yellow_tripdata_2023-01.parquet
```

## Download Clickhouse

```bash
curl https://clickhouse.com/ | sh
```

## Setup ClickHouse

Launch ClickHouse Local

```bash
./clickhouse local -m --path trips.chdb
```

Create database

```sql
CREATE DATABASE taxis;
```

```sql
USE taxis;
```

## Ingest data into ClickHouse

Create `trips` table

```sql
CREATE TABLE trips ENGINE MergeTree ORDER BY (tpep_pickup_datetime) AS 
from file('yellow tripdata Jan 2023.parquet', Parquet)
select *
SETTINGS schema_inference_make_columns_nullable = 0;
```

## Querying ClickHouse


Get all the amount columns

```sql
FROM trips 
SELECT COLUMNS('.*_amount')
LIMIT 10;
```

Get only the amount/fee/tax columns

```sql
FROM trips 
SELECT 
  COLUMNS('.*_amount|fee|tax')
LIMIT 10;
```

Compute the max of all the amount columns

```sql
FROM trips 
SELECT 
  COLUMNS('.*_amount|fee|tax')
  APPLY(max);
```


Compute the average of all the amount columns

```sql
FROM trips 
SELECT 
  COLUMNS('.*_amount|fee|tax')
  APPLY(avg);
```


Rounding (chaining functions)

```sql
FROM trips 
SELECT 
  COLUMNS('.*_amount|fee|tax')
  APPLY(avg)
  APPLY(round)
FORMAT Vertical;
```

Rounding to 2 dp (lambda)

```sql
FROM trips 
SELECT 
  COLUMNS('.*_amount|fee|tax')
  APPLY(avg)
  APPLY(col -> round(col, 2))
FORMAT Vertical;
```

Replace a field value

```sql
FROM trips 
SELECT 
  COLUMNS('.*_amount|fee|tax')
  REPLACE(total_amount*2 AS total_amount) 
  APPLY(avg)
  APPLY(col -> round(col, 2))
FORMAT Vertical;
```

Exclude a field

```sql
FROM trips 
SELECT 
  COLUMNS('.*_amount|fee|tax') EXCEPT(tolls_amount)
  REPLACE(total_amount*2 AS total_amount) 
  APPLY(avg)
  APPLY(col -> round(col, 2))
FORMAT Vertical;
```