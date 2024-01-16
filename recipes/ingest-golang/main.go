package main

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"io"
	"log"
	"os"
	"strconv"
)

func connect() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"localhost:9000"},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: "default",
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func main() {
	gzFile, err := os.Open("performance.csv.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		log.Fatal(err)
	}
	defer gzReader.Close()

	csvReader := csv.NewReader(gzReader)
	rowChan := make(chan []string)

	go func() {
		defer close(rowChan)

		if _, err := csvReader.Read(); err != nil {
			log.Fatal(err)
		}

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			rowChan <- record
		}
	}()

	conn, err := connect()
	if err != nil {
		panic((err))
	}

	ctx := context.Background()

	newBatch := func() driver.Batch {
		batch, err := conn.PrepareBatch(ctx, "INSERT INTO performance")
		if err != nil {
			panic(err)
		}
		return batch
	}

	batch := newBatch()

	batchSize := 0
	for row := range rowChan {
		tileX, _ := strconv.ParseFloat(row[2], 32)
		tileY, _ := strconv.ParseFloat(row[3], 32)
		downloadSpeedKbps, _ := strconv.ParseUint(row[4], 10, 32)
		uploadSpeedKbps, _ := strconv.ParseUint(row[5], 10, 32)
		latencyMs, _ := strconv.ParseUint(row[6], 10, 32)
		downloadLatencyMs, _ := strconv.ParseUint(row[7], 10, 32)
		uploadLatencyMs, _ := strconv.ParseUint(row[8], 10, 32)
		tests, _ := strconv.ParseUint(row[9], 10, 32)
		devices, _ := strconv.ParseUint(row[10], 10, 16)

		batchErr := batch.Append(
			tileX, tileY,
			downloadSpeedKbps, uploadSpeedKbps,
			latencyMs, downloadLatencyMs, uploadLatencyMs,
			tests, devices,
		)
		if batchErr != nil {
			panic(batchErr)
		}

		batchSize++

		if batchSize >= 1000 {
			if err := batch.Send(); err != nil {
				panic(err)
			}
			batch = newBatch()
			batchSize = 0
		}
	}

	if batchSize > 0 {
		if err := batch.Send(); err != nil {
			panic(err)
		}
	}
}
