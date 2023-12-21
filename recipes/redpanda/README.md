# Ingesting data from Redpanda to ClickHouse

In this recipe, we'll learn how to ingest data from Redpanda into ClickHouse.

## Setup Redpanda

Launch cluster

```bash
docker compose up
```

Create topic

```bash
rpk topic describe wiki_events -p
```

Ingest Wikimedia recent changes stream

```bash
curl -N https://stream.wikimedia.org/v2/stream/recentchange |
awk '/^data: /{gsub(/^data: /, ""); print}' |
kcat -P -b localhost:9092 -t wiki_events -Kø
```

## Setup ClickHouse

Launch ClickHouse Local

```bash
./clickhouse local -m --path wiki.chdb

Create database

```sql
CREATE DATABASE wiki;
```

```sql
USE wiki;
```

Create a table to connect to Redpanda

```sql
CREATE TABLE wikiQueue(
    id UInt32,
    type String,
    title String,
    title_url String,
    comment String,
    timestamp UInt64,
    user String,
    bot Boolean,
    server_url String,
    server_name String,
    wiki String,
    meta Tuple(uri String, id String, stream String, topic String, domain String)
)
ENGINE = Kafka('localhost:9092', 'wiki_events', 'consumer-group-wiki', 'JSONEachRow');
```

Create a table to store the data

```sql
CREATE TABLE wiki (
    dateTime DateTime64(3, 'UTC'),
    type String,
    title String,
    title_url String,
    id String,
    stream String,
    topic String,
    user String,
    bot Boolean, 
    server_name String,
    wiki String
) 
ENGINE = MergeTree 
ORDER BY dateTime;
```

Create a materialized view that ingests the data

```sql
CREATE MATERIALIZED VIEW wiki_mv TO wiki AS 
SELECT toDateTime(timestamp) AS dateTime,
       type, title, title_url, 
       tupleElement(meta, 'id') AS id, 
       tupleElement(meta, 'stream') AS stream, 
       tupleElement(meta, 'topic') AS topic, 
       user, bot, server_name, wiki
FROM wikiQueue;
```

## Querying Clickhouse

Count the number of records

```sql
FROM wiki SELECT count();
```

Find the most active users

```sql
FROM wiki
SELECT user, bot, COUNT(*) AS updates
GROUP BY user, bot
ORDER BY updates DESC
LIMIT 10;
```

With a nice plot

```sql
WITH users AS (
    FROM wiki
    SELECT user, bot, COUNT(*) AS updates
    GROUP BY user, bot
    ORDER BY updates DESC
)
SELECT
    user, bot,
    updates,
    bar(updates, 0, (SELECT max(updates) FROM users), 30) AS plot
FROM users
LIMIT 10;    
```