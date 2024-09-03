package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "blog-app/routers"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
    var err error
    
    // Replace these placeholders with your RDS MySQL credentials
    rdsEndpoint := "database-2.c5qcmqo02oo1.us-east-2.rds.amazonaws.com"
    rdsPort := "3306"
    dbUser := "admin"
    dbPassword := "Mayank1234"
    dbName := "blog_api_go" 

    // Construct the connection string without specifying a database name
    dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        dbUser,
        dbPassword,
        rdsEndpoint,
        rdsPort,
        dbName,
    )

    db, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        log.Fatalf("Error opening database connection: %v", err)
    }
    defer db.Close()

    // Ping the database to verify the connection
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error pinging database: %v", err)
    }

    // Create the database if it doesn't exist
    _, err = db.Exec(`CREATE DATABASE IF NOT EXISTS blog_api_go`)
    if err != nil {
        log.Fatal(err)
    }

    // Select the newly created database
    _, err = db.Exec("USE blog_api_go")
    if err != nil {
        log.Fatal(err)
    }

    // Create the users table if it doesn't exist
    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(100) NOT NULL,
        name VARCHAR(100),
        email VARCHAR(100) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
`)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Database 'blog_api_go' and table 'users' created successfully.")

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS blogs (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        username VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )
`)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Table 'blogs' created successfully.")

    router := routers.InitRouter(db)
    fmt.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", router))

}
