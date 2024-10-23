package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time"

    _ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() {
    var err error
    connStr := "user=postgres dbname=test_go host=127.0.0.1 sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Cannot connect to database:", err)
    }
    fmt.Println("Connected to the database!")
}

// Create table and column if they don't exist
func ensureTableExists() {
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS greetings (
        id SERIAL PRIMARY KEY,
        message TEXT NOT NULL
    );`
    _, err := db.Exec(createTableQuery)
    if err != nil {
        log.Fatal("Error creating table:", err)
    }

    // Check if the created_at column exists, and add it if not
    addColumnQuery := `
    DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                       WHERE table_name='greetings' AND column_name='created_at') THEN
            ALTER TABLE greetings ADD COLUMN created_at TIMESTAMP;
        END IF;
    END $$;`
    _, err = db.Exec(addColumnQuery)
    if err != nil {
        log.Fatal("Error adding column:", err)
    }

    fmt.Println("Table and column are ready!")
}

func insertGreeting() {
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    message := fmt.Sprintf("Welcome to my Go web server on %s", currentTime)
    
    _, err := db.Exec("INSERT INTO greetings (message, created_at) VALUES ($1, $2)", message, currentTime)
    if err != nil {
        log.Fatal("Error inserting greeting:", err)
    }
    fmt.Println("Inserted greeting:", message)
}

// Retrieve the latest greeting from the database
func getLatestGreeting() string {
    var message string
    err := db.QueryRow("SELECT message FROM greetings ORDER BY id DESC LIMIT 1").Scan(&message)
    if err != nil {
        log.Fatal("Error fetching message:", err)
    }
    return message
}

func handler(w http.ResponseWriter, r *http.Request) {
    message := getLatestGreeting()

    html := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>App</title>
        <style>
            body {
                background-color: #282c34;
                color: #61dafb;
                font-family: Arial, sans-serif;
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
                margin: 0;
            }
            .container {
                text-align: center;
            }
            h1 {
                font-size: 4rem;
                margin: 0;
            }
            p {
                font-size: 1.5rem;
                color: #fff;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Greetings!</h1>
            <p>` + message + `</p>
        </div>
    </body>
    </html>
    `
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, html)
}

func main() {
    // Connect to the database
    connectDB()

    // Ensure the greetings table and created_at column exist
    ensureTableExists()

    // Insert a dynamic greeting with date and time
    insertGreeting()

    // Set up HTTP handler
    http.HandleFunc("/", handler)
    fmt.Println("Serving on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
