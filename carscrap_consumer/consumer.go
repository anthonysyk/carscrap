package main

import (
	"carscrap/pkg/carscrap_redis"
	"carscrap/store"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	// Init MongoDB
	collections := store.Init()
	_ = collections

	// Init Redis
	svc := carscrap_redis.New(&carscrap_redis.NewInput{
		RedisURL: "0.0.0.0:6379",
	})

	// Subscribe to worker
	reply := make(chan []byte)
	err := svc.Subscribe("worker", reply)
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to report
	report := make(chan []byte)
	err = svc.Subscribe("report", report)
	if err != nil {
		log.Fatal(err)
	}

	// Start Listening
	fmt.Println("Listening to Redis broker ...")

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		for {
			msg := <-reply
			log.Printf("received: %q", string(msg))
			//TODO: Download Page
			//TODO: Extract Data
			//TODO: Store in MongoDB
		}
	}()
	go func() {
		for {
			msg := <-report
			log.Printf("report %q", string(msg))

			// Stop consumer when done
			if string(msg) == "done" {
				close(ch)
			}
		}
	}()

	// Graceful Shutdown
	<- ch
	fmt.Println("\nClosing MongoDB connection ...")
	collections.GetClient().Disconnect(context.TODO())
	fmt.Println("Closing Redis connection ...")
	svc.Close()
	fmt.Println("End of Program")
}
