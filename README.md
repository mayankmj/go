# Blog API

## Setup the go application
### git clone
### cd blog-api
### go mod tidy
### go run main.go

## Note: here , I am using aws-rds (mysql) for database

## User Endpoints
### Register User

- **Endpoint:** `/register`
- **Method:** `POST`
- **Description:** Register a new user.
- **Payload:**
    ```json
    {
      "name": "John Doe",
      "email": "john@example.com",
      "username": "johndoe",
      "password": "password123"
    }
    ```
- **cURL Example:**
    ```bash
    curl -X POST http://localhost:8080/register \
         -H "Content-Type: application/json" \
         -d '{
           "name": "John Doe",
           "email": "john@example.com",
           "username": "johndoe",
           "password": "password123"
         }'
    ```

### Login User

- **Endpoint:** `/login`
- **Method:** `POST`
- **Description:** Authenticate a user and return a token or success message.
- **Payload:**
    ```json
    {
      "email": "john@example.com",
      "password": "password123"
    }
    ```
- **cURL Example:**
    ```bash
    curl -X POST http://localhost:8080/login \
         -H "Content-Type: application/json" \
         -d '{
           "email": "john@example.com",
           "password": "password123"
         }'
    ```

## User Profile

### Get Profile

- **Endpoint:** `/profile`
- **Method:** `GET`
- **Description:** Fetch user profile details.
- **Response:**
    ```json
    {
      "username": "johndoe",
      "name": "John Doe"
    }
    ```
- **cURL Example:**
    ```bash
    curl -X GET http://localhost:8080/profile
    ```

### Update Profile

- **Endpoint:** `/profile`
- **Method:** `POST`
- **Description:** Update user profile details.
- **Payload:**
    ```json
    {
      "email": "john@example.com",
      "password": "password123",
      "new_name": "Mayank",
      "new_username": "new_username_test",
      "new_email": "johnnn@example.com"
    }
    ```
- **cURL Example:**
    ```bash
    curl -X POST http://localhost:8080/profile \
         -H "Content-Type: application/json" \
         -d '{
           "email": "john@example.com",
           "password": "password123",
           "new_name": "Mayank",
           "new_username": "new_username_test",
           "new_email": "johnnn@example.com"
         }'
    ```


## Blog Post Endpoints

### Get All Posts

- **Endpoint:** `/posts`
- **Method:** `GET`
- **Description:** Fetch all blog posts.
- **cURL Example:**
    ```bash
    curl -X GET http://localhost:8080/posts
    ```

### Get Post by ID

- **Endpoint:** `/posts/{id}`
- **Method:** `GET`
- **Description:** Fetch a blog post by its ID.
- **cURL Example:**
    ```bash
    curl -X GET http://localhost:8080/posts/1
    ```

### Create Post

- **Endpoint:** `/createpost`
- **Method:** `POST`
- **Description:** Create a new blog post.
- **Payload:**
    ```json
    {
      "email": "john@example.com",
      "password": "password123",
      "name": "Post Author",
      "title": "Post Title",
      "content": "This is the post content."
    }
    ```
- **cURL Example:**
    ```bash
    curl -X POST http://localhost:8080/createpost \
         -H "Content-Type: application/json" \
         -d '{
           "email": "john@example.com",
           "password": "password123",
           "name": "Post Author",
           "title": "Post Title",
           "content": "This is the post content."
         }'
    ```

### Update Post

- **Endpoint:** `/updatepost/{id}`
- **Method:** `PUT`
- **Description:** Update a blog post by its ID.
- **Payload:**
    ```json
    {
      "email": "john@example.com",
      "password": "password123",
      "title": "Updated Title",
      "content": "Updated content for the post."
    }
    ```
- **cURL Example:**
    ```bash
    curl -X PUT http://localhost:8080/updatepost/1 \
         -H "Content-Type: application/json" \
         -d '{
           "email": "john@example.com",
           "password": "password123",
           "title": "Updated Title",
           "content": "Updated content for the post."
         }'
    ```

### Delete Post

- **Endpoint:** `/deletepost/{id}`
- **Method:** `DELETE`
- **Description:** Delete a blog post by its ID.
- **Payload:**
    ```json
    {
      "email": "john@example.com",
      "password": "password123"
    }
    ```
- **cURL Example:**
    ```bash
    curl -X DELETE http://localhost:8080/deletepost/1 \
         -H "Content-Type: application/json" \
         -d '{
           "email": "john@example.com",
           "password": "password123"
         }'
    ```

