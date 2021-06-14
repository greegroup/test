package main

import (
   "context"
   "fmt"
   "os"
   "time"

   "github.com/jackc/pgx/v4"
   "github.com/jackc/pgx/v4/pgxpool"
)

func main() {
   /********************************************/
   /* Connect using Connection Pool            */
   /********************************************/
   ctx := context.Background()
   connStr := "yourConnectionStringHere"
   dbpool, err := pgxpool.Connect(ctx, connStr)
   if err != nil {
       fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
       os.Exit(1)
   }
   defer dbpool.Close()

   // Generate data to insert

   //SQL query to generate sample data
   queryDataGeneration := `
       SELECT generate_series(now() - interval '24 hour', now(), interval '5 minute') AS time,
       floor(random() * (3) + 1)::int as sensor_id,
       random()*100 AS temperature,
       random() AS cpu
       `
   //Execute query to generate samples for sensor_data hypertable
   rows, err := dbpool.Query(ctx, queryDataGeneration)
   if err != nil {
       fmt.Fprintf(os.Stderr, "Unable to generate sensor data: %v\n", err)
       os.Exit(1)
   }
   defer rows.Close()
   fmt.Println("Successfully generated sensor data\n")

   //Store data generated in slice results
   type result struct {
       Time        time.Time
       SensorId    int
       Temperature float64
       CPU         float64
   }
   var results []result
   for rows.Next() {
       var r result
       err = rows.Scan(&r.Time, &r.SensorId, &r.Temperature, &r.CPU)
       if err != nil {
           fmt.Fprintf(os.Stderr, "Unable to scan %v\n", err)
           os.Exit(1)
       }
       results = append(results, r)
   }
   // Any errors encountered by rows.Next or rows.Scan will be returned here
   if rows.Err() != nil {
       fmt.Fprintf(os.Stderr, "rows Error: %v\n", rows.Err())
       os.Exit(1)
   }

   // Check contents of results slice
   /*fmt.Println("Contents of RESULTS slice")
   for i := range results {
       var r result
       r = results[i]
       fmt.Printf("Time: %s | ID: %d | Temperature: %f | CPU: %f |\n", &r.Time, r.SensorId, r.Temperature, r.CPU)
   }*/

   //Insert contents of results slice into TimescaleDB
   //SQL query to generate sample data
   queryInsertTimeseriesData := `
   INSERT INTO sensor_data (time, sensor_id, temperature, cpu) VALUES ($1, $2, $3, $4);
   `

   /********************************************/
   /* Batch Insert into TimescaleDB            */
   /********************************************/
   //create batch
   batch := &pgx.Batch{}
   numInserts := len(results)
   //load insert statements into batch queue
   for i := range results {
       var r result
       r = results[i]
       batch.Queue(queryInsertTimeseriesData, r.Time, r.SensorId, r.Temperature, r.CPU)
   }
   batch.Queue("select count(*) from sensor_data")

   //send batch to connection pool
   br := dbpool.SendBatch(ctx, batch)
   //execute statements in batch queue
   for i := 0; i < numInserts; i++ {
       _, err := br.Exec()
       if err != nil {
           fmt.Fprintf(os.Stderr, "Unable to execute statement in batch queue %v\n", err)
           os.Exit(1)
       }
   }
   fmt.Println("Successfully batch inserted data n")

   //Compare length of results slice to size of table
   fmt.Println("size of results: %d\n", len(results))
   //check size of table for number of rows inserted
   // result of last SELECT statement
   var rowsInserted int
   err = br.QueryRow().Scan(&rowsInserted)
   fmt.Println("size of table: %d\n", rowsInserted)

   err = br.Close()
   if err != nil {
       fmt.Fprintf(os.Stderr, "Unable to closer batch %v\n", err)
       os.Exit(1)
   }

}
