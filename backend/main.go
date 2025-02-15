package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "backend/handlers"

    _ "github.com/go-sql-driver/mysql"
    "github.com/rs/cors"
)

var db *sql.DB

func connectDB() {
    var err error
    // Build connection string using environment variables
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"))

    // Open database connection
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }

    // Set connection pool parameters
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Add retry logic for initial connection
    for i := 0; i < 30; i++ {
        err = db.Ping()
        if err == nil {
            break
        }
        log.Printf("Waiting for database... Attempt %d/30", i+1)
        time.Sleep(2 * time.Second)
    }
    if err != nil {
        log.Fatal("Database is not reachable:", err)
    }

    fmt.Println("✅ Connected to MySQL successfully!")
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "Welcome to the Go backend!"}`))
}

func setupRoutes() {
    // Base routes
    http.HandleFunc("/", welcomeHandler)
    
    // User routes
    http.HandleFunc("/login", handlers.LoginHandler(db))
    http.HandleFunc("/signup", handlers.SignupHandler(db))
    http.HandleFunc("/users", handlers.GetUsersHandler(db))
    http.HandleFunc("/users/delete", handlers.DeleteUserHandler(db))
    
    // Post routes
    http.HandleFunc("/posts", handlers.GetPostsHandler(db))
    http.HandleFunc("/posts/create", handlers.CreatePostHandler(db))
    http.HandleFunc("/posts/delete", handlers.DeletePostHandler(db))
    http.HandleFunc("/posts/edit", handlers.EditPostHandler(db))
}

func main() {
    // Set default values for environment variables if not set
    if os.Getenv("DB_HOST") == "" {
        os.Setenv("DB_HOST", "mysql")
    }
    if os.Getenv("DB_PORT") == "" {
        os.Setenv("DB_PORT", "3306")
    }
    if os.Getenv("DB_USER") == "" {
        os.Setenv("DB_USER", "apptimus")
    }
    if os.Getenv("DB_PASSWORD") == "" {
        os.Setenv("DB_PASSWORD", "1q2w3e")
    }
    if os.Getenv("DB_NAME") == "" {
        os.Setenv("DB_NAME", "apptimus_db")
    }

    connectDB()
    defer db.Close()

    setupRoutes()

    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:3000",
            "http://localhost",
            "http://frontend:3000",
        },
        AllowedMethods: []string{
            "GET", "POST", "PUT", "DELETE", "OPTIONS",
        },
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
		},
        AllowCredentials: true,
        // Add debug logging for CORS issues
        Debug: true,
    }).Handler(http.DefaultServeMux)

    fmt.Println("🚀 Server is running on :8080")
    log.Fatal(http.ListenAndServe(":8080", corsHandler))
}