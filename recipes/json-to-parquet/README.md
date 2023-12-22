# Converting JSON to Parquet

In this recipe, we'll learn a variety of ways to convert JSON to Parquet in ClickHouse.

## Download Clickhouse

```bash
curl https://clickhouse.com/ | sh
```

## Setup ClickHouse

Launch ClickHouse Local

```bash
./clickhouse local -m
```

## Querying JSON

Count the number of records

```sql
FROM 'data/movies.json.gz' SELECT count();
```

Explore the record structure

```sql
DESCRIBE 'data/movies.json.gz'
SETTINGS describe_compact_output=1;
```

```sql
FROM 'data/movies.json.gz' 
SELECT *
LIMIT 1
Format Vertical;
```

## JSON to Parquet

One file

```sql
FROM 'data/movies.json.gz' 
SELECT *
INTO OUTFILE 'data/movies.parquet' 
FORMAT Parquet;
```

Partition by language

```sql
INSERT INTO FUNCTION 
file('data/movies_lang_{_partition_id}.parquet', 'Parquet') 
PARTITION BY original_language
select *
from file('data/movies.json.gz');
```

Partition by vote average buckets

```sql
INSERT INTO FUNCTION 
file('data/movies_vote_{_partition_id}.parquet', 'Parquet') 
PARTITION BY multiIf(
    vote_average = 10,
    '9-10',
    vote_average = floor(vote_average),
    toString(vote_average) || '-' || toString(vote_average +1),
    toString(floor(vote_average)) || '-' || toString(ceil(vote_average))
)
select *
from file('data/movies.json.gz');
```