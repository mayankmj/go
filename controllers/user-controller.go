package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "blog-app/models"
    "golang.org/x/crypto/bcrypt"
    "fmt"
    "regexp"
    "blog-app/utils"
)
func validatePassword (pass string) bool {
    var password = regexp.MustCompile(`^(?=.*[A-Z].*[A-Z])(?=.*[!@#$&*])(?=.*[0-9].*[0-9])(?=.*[a-z].*[a-z].*[a-z]).{8,}$`)
    return password.MatchString(pass)
}
func Register(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            fmt.Println("Error decoding JSON:", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Validate the password length
        // reg = "^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$"
        if len(user.Password) < 6 {
            fmt.Println("Password too short:", user.Password)
            http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
            return
        }

        if !validatePassword(user.Password) {
            fmt.Println("Password should match the minimum requirements")
            http.Error(w,"Recheck the passwrod",http.StatusBadRequest)
            return
        } 

        // Hash the password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            fmt.Println("Error hashing password:", err)
            http.Error(w, "Error creating user", http.StatusInternalServerError)
            return
        }
        user.Password = string(hashedPassword)

        // Create the user in the database
        err = user.CreateUser(db)
        if err != nil {
            fmt.Println("Error creating user in database:", err)
            http.Error(w, "Error creating user", http.StatusInternalServerError)
            return
        }

        // Return a success response
        fmt.Println("User registered successfully:", user)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

func Login(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var loginRequest models.LoginRequest
        err := json.NewDecoder(r.Body).Decode(&loginRequest)
        if err != nil {
            fmt.Println("Error decoding JSON:", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        fmt.Println("User login data:", loginRequest)

        // Fetch user by email
        storedUser, err := models.GetUserByEmail(db, loginRequest.Email)
        if err != nil {
            fmt.Println("Error fetching user by email:", err)
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        fmt.Println("Stored user data:", storedUser)

        // Compare the hashed password
        err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginRequest.Password))
        if err != nil {
            fmt.Println("Password comparison failed:", err)
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        // Generate JWT token
        token, err := utils.GenerateJWT(storedUser)
        if err != nil {
            fmt.Println("Error generating JWT:", err)
            http.Error(w, "Error generating token", http.StatusInternalServerError)
            return
        }

        fmt.Println("Login successful, token generated:", token)

        // Send the token in the response
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}










