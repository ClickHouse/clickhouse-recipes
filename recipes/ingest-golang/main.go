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
	"sync"
)

func connectToClickHouse(host string) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{host},
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

func ingestRecords(wg *sync.WaitGroup, rowChan <-chan []string, conn driver.Conn, batchSize int) {
	defer wg.Done()

	newBatch := func() driver.Batch {
		ctx := context.Background()
		batch, err := conn.PrepareBatch(ctx, "INSERT INTO performance")
		if err != nil {
			panic(err)
		}
		return batch
	}
	batch := newBatch()
	recordsProcessed := 0
	for row := range rowChan {
		quadKey := row[0]
		tileX, _ := strconv.ParseFloat(row[2], 32)
		tileY, _ := strconv.ParseFloat(row[3], 32)
		downloadSpeedKbps, _ := strconv.ParseUint(row[4], 10, 32)
		uploadSpeedKbps, _ := strconv.ParseUint(row[5], 10, 32)
		latencyMs, _ := strconv.ParseUint(row[6], 10, 32)
		downloadLatencyMs, _ := strconv.ParseUint(row[7], 10, 32)
		uploadLatencyMs, _ := strconv.ParseUint(row[8], 10, 32)
		tests, _ := strconv.ParseUint(row[9], 10, 32)
		devices, _ := strconv.ParseUint(row[10], 10, 16)

		err := batch.Append(
			quadKey,
			tileX, tileY,
			downloadSpeedKbps, uploadSpeedKbps,
			latencyMs, downloadLatencyMs, uploadLatencyMs,
			tests, devices,
		)
		if err != nil {
			log.Fatal(err)
		}

		recordsProcessed++

		if recordsProcessed%batchSize == 0 {
			if err := batch.Send(); err != nil {
				log.Fatal(err)
			}
			batch = newBatch()
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal(err)
	}
}

func readCSVToChannel(filePath string, rowChan chan<- []string) {
	gzFile, err := os.Open(filePath)
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

	defer close(rowChan)

	if _, err := csvReader.Read(); err != nil { // Skip header or handle error
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
}

func main() {
	rowChan := make(chan []string)
	go readCSVToChannel("performance.csv.gz", rowChan)

	conn, err := connectToClickHouse("localhost:9000")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	const numWorkers = 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go ingestRecords(&wg, rowChan, conn, 10_000)
	}

	wg.Wait()
}
