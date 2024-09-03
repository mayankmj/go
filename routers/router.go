package routers

import (
    "database/sql"
    "blog-app/controllers"
    "github.com/gorilla/mux"
)

func InitRouter(db *sql.DB) *mux.Router {
    router := mux.NewRouter()

    // User endpoints
    router.HandleFunc("/register", controllers.Register(db)).Methods("POST")
    router.HandleFunc("/login", controllers.Login(db)).Methods("POST")
    router.HandleFunc("/profile", controllers.ProfileHandler(db)).Methods("GET", "POST")

    // Blog post endpoints
    router.HandleFunc("/posts", controllers.GetAllPosts(db)).Methods("GET") // Fetch all posts
    router.HandleFunc("/posts/{id:[0-9]+}", controllers.GetPostByID(db)).Methods("GET") // Fetch post by ID
    router.HandleFunc("/createpost", controllers.CreatePost(db)).Methods("POST") // Create a post
    router.HandleFunc("/updatepost/{id:[0-9]+}", controllers.UpdatePost(db)).Methods("PUT") // Update a post
    router.HandleFunc("/deletepost/{id:[0-9]+}", controllers.DeletePost(db)).Methods("DELETE") // Delete a post

    return router
}

