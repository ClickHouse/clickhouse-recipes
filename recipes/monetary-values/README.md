# Working with monetary values in ClickHouse

In this recipe, we'll learn how to work with monetary values in ClickHouse

The data used is a subset of `https://datasets-documentation.s3.eu-west-3.amazonaws.com/forex/csv/year_month/*.csv.zst`.

## Download Clickhouse

```bash
curl https://clickhouse.com/ | sh
```

## Setup ClickHouse

Launch ClickHouse Local

```bash
./clickhouse local -m
```

## Running monetary computations

Return the first 10 values

```sql
FROM 'data/forex.csv'
SELECT *
LIMIT 10;
```

Compute the (imprecise) diff between bid/ask prices

```sql
FROM 'data/forex.csv'
SELECT bid, ask, bid - ask AS diff
LIMIT 10;
```

**Don't do this!**

Instead, use the [`toDecimal32`](https://clickhouse.com/docs/en/sql-reference/functions/type-conversion-functions#todecimal3264128256) or [`toDecimal64`](https://clickhouse.com/docs/en/sql-reference/functions/type-conversion-functions#todecimal3264128256) functions.
The second parameter of this function is the number of decimal places.

```sql
FROM 'data/forex.csv'
SELECT toDecimal32(bid, 4) AS bid, toDecimal32(ask, 4) AS ask, bid - ask AS diff
LIMIT 10;
```

## Storing monetary values

When storing monetary values, use the `Decimal` data type.

```sql
CREATE TABLE forex (
	`datetime` DateTime64(3),
	`bid` Decimal(11, 5) CODEC(ZSTD(1)),
	`ask` Decimal(11, 5) CODEC(ZSTD(1)),
	`base` LowCardinality(String),
	`quote` LowCardinality(String)
)
ENGINE = MergeTree
ORDER BY (base, quote, datetime);
```

```sql
INSERT INTO forex
SELECT * 
FROM 'data/forex.csv';
```

You can then query the table

```sql
FROM forex 
select *, bid - ask  AS diff
LIMIT 10;
```
