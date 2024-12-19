package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware for logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        
        // Log incoming request
        log.Printf("Incoming %s request from %s to %s", r.Method, r.RemoteAddr, r.URL.Path)
        
        next(w, r)
        
        // Log request completion
        log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(startTime))
    }
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    log.Printf("Processing request at root endpoint")
    fmt.Fprintf(w, "Hello from Go!")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    log.Printf("Health check requested")
    fmt.Fprintf(w, "Service is healthy!")
}

func main() {
    // Configure logging
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    
    // Register routes with logging middleware
    http.HandleFunc("/", loggingMiddleware(handleRoot))
    http.HandleFunc("/health", loggingMiddleware(handleHealth))
    
    port := ":8080"
    log.Printf("Server starting on port %s...", port)
    
    // Use log.Fatal to catch and log startup errors
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal("Server failed to start: ", err)
    }
}
