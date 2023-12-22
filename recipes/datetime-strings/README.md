# Parsing DateTime Strings

In this recipe, we'll learn how to parse DateTime strings in ClickHouse.

## Download Clickhouse

```bash
curl https://clickhouse.com/ | sh
```

## Setup ClickHouse

Launch ClickHouse Local

```bash
./clickhouse local -m
```

## Parsing DateTime strings

### Using Format String

See [syntax documentation](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#formatDateTime).

```sql
WITH '2023-12-21T00:01:13.607111Z' AS dateString
SELECT
    dateString,
    parseDateTime(dateString, '%Y-%m-%dT%H:%i:%s.%fZ') AS date;
```

```sql
WITH 'Mon 17 December 2023 00:01:13' AS dateString
SELECT
    dateString,
    parseDateTime(dateString, '%a %e %M %Y %H:%i:%s') AS date;
```

```sql
WITH '12/17/2023 00:01:13 -0500' AS dateString
SELECT
    dateString,
    parseDateTime(dateString, '%m/%d/%Y %H:%i:%s %z') AS date;
```

```sql
WITH '2023-12-17' AS dateString
SELECT
    dateString,
    parseDateTime(dateString, '%F') AS date;
```

### Using Joda syntax

See [syntax documentation](https://joda-time.sourceforge.net/apidocs/org/joda/time/format/DateTimeFormat.html)

[NOTE]
====
Doesn't support fractional seconds
====

```sql
WITH '2023-12-21T00:01:13' AS dateString
SELECT
    dateString,
    parseDateTimeInJodaSyntax(dateString, 'YYYY-MM-dd''T''HH:mm:ss') AS date;
```

```sql
WITH '2023-12-21 12:01:13 AM' AS dateString
SELECT
    dateString,
    parseDateTimeInJodaSyntax(dateString, 'YYYY-MM-dd hh:mm:ss a') AS date;
```

## Best effort parsing

```sql
WITH '2023-12-21T00:01:13.607111Z' AS dateString
SELECT
    dateString,
    parseDateTime32BestEffort(dateString) AS date32,
    parseDateTime64BestEffort(dateString) AS date64;
```