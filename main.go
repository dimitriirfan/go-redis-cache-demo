package main

import (
	"context"
	"go-redis/cache"
	"go-redis/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctx := context.Background()
	logger := log.New(os.Stdout, "main-api", log.LstdFlags)
	mux := http.NewServeMux()

	redisCache := cache.NewRedisCache(logger, "localhost:6379", "", 0, 1*time.Hour, ctx)
	handlerPhoto := handler.NewPhotoHandler(logger, redisCache)

	mux.Handle("/photos", handlerPhoto)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Println("Server started on port 8080")
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)
	// http.ListenAndServe(":8080", serveMux)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)
}
