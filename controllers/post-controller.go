package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "blog-app/models"
    "golang.org/x/crypto/bcrypt"
    "fmt"
    "time"
    "github.com/gorilla/mux"
)
func ProfileHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getProfile(db, w, r)
        case http.MethodPost:
            updateProfile(db, w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }
}

func getProfile(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Decode the request body
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Fetch the user by email
    user, err := models.GetUserByEmail(db, req.Email)
    if err != nil {
        fmt.Println("Error fetching user by email:", err)
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Check the password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        fmt.Println("Password comparison failed:", err)
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    // Return the user profile (only name, username, and email)
    response := struct {
        Name     string `json:"name"`
        Username string `json:"username"`
        Email    string `json:"email"`
    }{
        Name:     user.Name,
        Username: user.Username,
        Email:    user.Email,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func updateProfile(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email       string `json:"email"`
        Password    string `json:"password"`
        NewName     string `json:"new_name"`
        NewUsername string `json:"new_username"`
        NewEmail    string `json:"new_email"`
    }

    // Decode the request body
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Fetch the user by email
    user, err := models.GetUserByEmail(db, req.Email)
    if err != nil {
        fmt.Println("Error fetching user by email:", err)
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Check the password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        fmt.Println("Password comparison failed:", err)
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    // Update the user profile
    updateQuery := `
        UPDATE users
        SET name = COALESCE(?, name),
            username = COALESCE(?, username),
            email = COALESCE(?, email)
        WHERE email = ?`

    _, err = db.Exec(updateQuery, req.NewName, req.NewUsername, req.NewEmail, req.Email)
    if err != nil {
        fmt.Println("Error updating user profile:", err)
        http.Error(w, "Error updating user profile", http.StatusInternalServerError)
        return
    }

    fmt.Println("Profile updated successfully")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func authenticateUser(db *sql.DB, email, password string) (int, error) {
    user, err := models.GetUserByEmail(db, email)
    if err != nil {
        return 0, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    // if err != nil {
    //     return 0, errors.New("invalid email or password")
    // }

    return user.ID, nil
}

func CreatePost(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var requestData map[string]interface{}
        err := json.NewDecoder(r.Body).Decode(&requestData)
        if err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        email, _ := requestData["email"].(string)
        password, _ := requestData["password"].(string)
        name, _ := requestData["name"].(string)
        title, _ := requestData["title"].(string)
        content, _ := requestData["content"].(string)

        if email == "" || password == "" || title == "" || content == "" {
            http.Error(w, "Email, password, title, and content are required", http.StatusBadRequest)
            return
        }

        // Authenticate the user
        userID, err := authenticateUser(db, email, password)
        if err != nil {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        // Fetch user details based on userID
        user, err := models.GetUserByID(db, userID)
        if err != nil {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }

        // Create the post
        post := &models.Post{
            Name:      name,
            Title:     title,
            Content:   content,
            Username:  user.Username,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }

        err = models.CreatePost(db, post)
        if err != nil {
            http.Error(w, "Error creating post", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(post)
    }
}

func GetAllPosts(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        posts, err := models.GetAllPosts(db)
        if err != nil {
            http.Error(w, "Error fetching posts", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(posts)
    }
}

func GetPostByID(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        postIDStr, ok := vars["id"]
        if !ok {
            http.Error(w, "Post ID not provided", http.StatusBadRequest)
            return
        }

        post, err := models.GetPostByID(db, postIDStr)
        if err != nil {
            http.Error(w, "No post exists with this post_id, please try again!", http.StatusNotFound)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(post)
    }
}


func UpdatePost(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract postID from URL parameters
        vars := mux.Vars(r)
        postIDStr, ok := vars["id"]
        if !ok {
            http.Error(w, "Post ID not provided", http.StatusBadRequest)
            return
        }

        // Decode JSON request body
        var requestData map[string]interface{}
        err := json.NewDecoder(r.Body).Decode(&requestData)
        if err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Extract other fields from request body
        email, _ := requestData["email"].(string)
        password, _ := requestData["password"].(string)
        newTitle, _ := requestData["title"].(string)
        newContent, _ := requestData["content"].(string)

        // Authenticate user and get user ID
        userID, err := authenticateUser(db, email, password)
        if err != nil {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        // Fetch user data using the userID
        user, err := models.GetUserByID(db, userID)
        if err != nil {
            http.Error(w, "Error fetching user data", http.StatusInternalServerError)
            return
        }

        // Get existing post by ID
        post, err := models.GetPostByID(db, postIDStr)
        if err != nil {
            http.Error(w, "No post exists with this post_id, please try again!", http.StatusNotFound)
            return
        }

        // Check authorization by comparing usernames
        if post.Username != user.Username {
            http.Error(w, "You are not authorized to update this post", http.StatusForbidden)
            return
        }

        // Update post details
        post.Title = newTitle
        post.Content = newContent
        post.UpdatedAt = time.Now()

        // Save changes
        err = models.UpdatePost(db, post)
        if err != nil {
            http.Error(w, "Error updating post", http.StatusInternalServerError)
            return
        }

        // Respond with updated post
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(post)
    }
}



func DeletePost(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract postID from URL parameters
        vars := mux.Vars(r)
        postIDStr, ok := vars["id"]
        if !ok {
            http.Error(w, "Post ID not provided", http.StatusBadRequest)
            return
        }

        // Decode JSON request body
        var requestData map[string]interface{}
        err := json.NewDecoder(r.Body).Decode(&requestData)
        if err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Extract user credentials from request body
        email, _ := requestData["email"].(string)
        password, _ := requestData["password"].(string)

        // Authenticate user
        userID, err := authenticateUser(db, email, password)
        if err != nil {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        // Get the authenticated user's information (assuming you have a function for that)
        user, err := models.GetUserByID(db, userID)
        if err != nil {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }

        // Get existing post by ID
        post, err := models.GetPostByID(db, postIDStr)
        if err != nil {
            http.Error(w, "Post not found", http.StatusNotFound)
            return
        }

        // Check authorization using Username instead of UserID
        if post.Username != user.Username {
            http.Error(w, "You are not authorized to delete this post", http.StatusForbidden)
            return
        }

        // Delete the post
        err = models.DeletePost(db, postIDStr)
        if err != nil {
            http.Error(w, "Error deleting post", http.StatusInternalServerError)
            return
        }

        // Respond with success message
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
    }
}
