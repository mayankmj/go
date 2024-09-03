package models

import (
    "database/sql"
    "errors"
    "time"
    "fmt"
)

// for User
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    CreatedAt time.Time `json:"created_at"`
}

// for blog post
type Post struct {
    ID         int       `json:"id"`
    Name       string    `json:"name"`
    Title      string    `json:"title"`
    Content    string    `json:"content"`
    Username   string    `json:"username"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func (u *User) CreateUser(db *sql.DB) error {
    stmt, err := db.Prepare("INSERT INTO users (username, name, email, password) VALUES (?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    result, err := stmt.Exec(u.Username, u.Name, u.Email, u.Password)
    if err != nil {
        return err
    }

    userId, err := result.LastInsertId()
    if err != nil {
        return err
    }
    u.ID = int(userId)

    return nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
    var user User
    query := `SELECT id, name, username, email, password, created_at FROM users WHERE email = ?`
    row := db.QueryRow(query, email)
    fmt.Println("Executing query:", query)
    err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    return &user, nil
}

func GetUserByID(db *sql.DB, userID int) (*User, error) {
    var user User
    query := `SELECT id, name, username, email, password, created_at FROM users WHERE id = ?`
    row := db.QueryRow(query, userID)
    fmt.Println("Executing query:", query)
    err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    return &user, nil
}


func (user *User) UpdateProfile(db *sql.DB) error {
    query := `UPDATE users SET name = ?, username = ?, email = ? WHERE id = ?`
    _, err := db.Exec(query, user.Name, user.Username, user.Email, user.ID)
    if err != nil {
        // You can log the error here for debugging purposes
        fmt.Println("Error updating user profile:", err)
        return err
    }
    return nil
}


// CreatePost inserts a new post into the database
func CreatePost(db *sql.DB, post *Post) error {
    query := `INSERT INTO blogs (title, name, content, username, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
    _, err := db.Exec(query, post.Title, post.Name, post.Content, post.Username, post.CreatedAt, post.UpdatedAt)
    if err != nil {
        fmt.Println("Error executing query:", err)
        return err
    }
    return nil
}




// GetAllPosts retrieves all posts from the database
func GetAllPosts(db *sql.DB) ([]Post, error) {
    query := `SELECT id, name, title, content, username, created_at, updated_at FROM blogs`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.Name, &post.Title, &post.Content, &post.Username, &post.CreatedAt, &post.UpdatedAt)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return posts, nil
}

// GetPostByID retrieves a post by its ID
func GetPostByID(db *sql.DB, postID string) (*Post, error) {
    var post Post
    query := `SELECT id, name, title, content, username, created_at, updated_at FROM blogs WHERE id = ?`
    row := db.QueryRow(query, postID)
    err := row.Scan(&post.ID, &post.Name, &post.Title, &post.Content, &post.Username, &post.CreatedAt, &post.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("post not found")
        }
        return nil, err
    }
    return &post, nil
}


// UpdatePost updates an existing post in the database
func UpdatePost(db *sql.DB, post *Post) error {
    query := `UPDATE blogs SET title = ?, content = ?, updated_at = ? WHERE id = ? AND username = ?`
    _, err := db.Exec(query, post.Title, post.Content, post.UpdatedAt, post.ID, post.Username)
    return err
}

// DeletePost deletes a post from the database
func DeletePost(db *sql.DB, postID string) error {
    _, err := db.Exec(`DELETE FROM blogs WHERE id = ?`, postID) 
    return err
}




