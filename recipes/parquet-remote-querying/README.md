# Remote querying Parquet files with ClickHouse

In this recipe, we'll learn how to query remote Parquet files with ClickHouse

## Download Clickhouse

```bash
curl https://clickhouse.com/ | sh
```

## Setup ClickHouse

Launch ClickHouse Local

```bash
./clickhouse local -m --path trips.chdb
```

## Querying remote Parquet files

The dataset that we’re going to use is available at `vivym/midjourney-messages`.
It contains just over 55 million records spread over 56 Parquet files. 
They’re a little over 150 MB each, for a total of around 8 GB.

The files have the following structure:


```bash
https://huggingface.co/datasets/vivym/midjourney-messages/resolve/main/data/000000.parquet
https://huggingface.co/datasets/vivym/midjourney-messages/resolve/main/data/000001.parquet
...
https://huggingface.co/datasets/vivym/midjourney-messages/resolve/main/data/000055.parquet
```

Get one row from one file

```sql
FROM url('https://huggingface.co/datasets/vivym/midjourney-messages/resolve/main/data/000000.parquet')
SELECT *
LIMIT 1
Format JSONEachRow
SETTINGS max_http_get_redirects=1;
```

Summing the `size` column in one Parquet file

```sql
FROM url('https://huggingface.co/datasets/vivym/midjourney-messages/resolve/main/data/000000.parquet')
SELECT sum(size) as totalSize, formatReadableSize(totalSize)
SETTINGS max_http_get_redirects=1;
```

Summing the `size` column of all the files

```sql
SELECT
    sum(size) AS totalSize,
    formatReadableSize(totalSize)
FROM url('https://huggingface.co/datasets/vivym/midjourney-messages/resolve/main/data/0000{00..55}.parquet')
SETTINGS max_http_get_redirects = 1
```