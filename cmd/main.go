package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/tomcase/syncopation-api/controllers"
	"github.com/tomcase/syncopation-api/data"
	"github.com/tomcase/syncopation-api/middleware"
	"github.com/tomcase/syncopation-api/sync"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	db := &data.Db{}
	err := db.Migrate()
	if err != nil {
		log.Fatalf("Failed to Migrate: %v", err)
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minutes().Do(func() {
		log.Println("Starting Sync")
		syncErr := sync.Go(db)
		if syncErr != nil {
			log.Printf("Failed to sync: %v\n", syncErr)
		}
		log.Println("Finished Sync")
	})
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	s.StartAsync()

	r := mux.NewRouter().StrictSlash(true)

	prefix := "/api"
	controllers.RegisterHandlers(r, prefix, db)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware)

	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("API_PORT is has not been defined.")
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	log.Println("Api Starting!")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println(fmt.Sprintf("HTTP Server is running on http://localhost:%s", port))

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down...")
	os.Exit(0)
}
