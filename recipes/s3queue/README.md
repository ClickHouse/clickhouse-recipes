# Continuously ingesting data from S3 into ClickHouse

In this recipe, we'll learn how to ingest data continuously from an S3 bucket into ClickHouse.

## Download Clickhouse

```bash
curl https://clickhouse.com/ | sh
```

## ClickHouse Server

```bash
cp clickhouse clickhouse-server
cd clickhouse-server
```

Start the server

```bash
./clickhouse server
```


## ClickHouse Client

Create a table with `S3Queue` table engine

```sql
CREATE TABLE ordersQueue (
    orderDate DateTime, 
    gender String,
    customerId UUID,
    cost Float32,
    name String,
    creditCardNumber String,
    address String,
    orderId UUID
)
ENGINE = S3Queue(
    'https://s3queue.clickhouse.com.s3.eu-north-1.amazonaws.com/data/*.json',
    JSONEachRow
)
SETTINGS 
    mode = 'ordered', 
    s3queue_enable_logging_to_s3queue_log = 1;
```

Create a table with the `MergeTree` table engine

```sql
CREATE TABLE orders (
    orderDate DateTime, 
    gender String,
    customerId UUID,
    cost Float32,
    name String,
    creditCardNumber String,
    address String,
    orderId UUID
)
ENGINE = MergeTree 
ORDER BY (customerId, orderDate);
```

Create a materialized view that reads data from S3 and writes it into the `orders` table

```sql
CREATE MATERIALIZED VIEW ordersConsumer TO orders AS 
SELECT * 
FROM ordersQueue;
```

## Querying the data

We should now see data coming into the `orders` table.
We can check on the ingestion progress by writing the following query:

```sql
FROM logs 
SELECT count(), 
       formatReadableQuantity(count()) AS countFriendly, 
       now() 
Format PrettyNoEscapes;

```