Creating a clean MVC architecture in Go with the Gin framework, while avoiding tight coupling, involves careful planning and adherence to best practices. Here's a comprehensive approach:

**1. Project Structure:**

```
your-project/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── controllers/
│   │   ├── user_controller.go
│   │   └── product_controller.go
│   ├── models/
│   │   ├── user.go
│   │   └── product.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   └── product_repository.go
│   ├── services/
│   │   ├── user_service.go
│   │   └── product_service.go
├── pkg/
│   └── database/
│       └── database.go
├── go.mod
├── go.sum
```

**2. Model Layer (Models):**

* Define structs that represent your data entities.
* Keep models pure data structures, without any business logic.

```go
// internal/models/user.go
package models

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

**3. Repository Layer (Repositories):**

* Implement data access logic.
* Use interfaces to define repository contracts, enabling dependency injection.
* This layer will interact with the database (or other data sources).

```go
// internal/repositories/user_repository.go
package repositories

import (
    "your-project/internal/models"
    "your-project/pkg/database"
)

type UserRepository interface {
    GetUserByID(id int) (*models.User, error)
    CreateUser(user *models.User) error
    // ... other methods ...
}

type userRepository struct {
    db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
    // Database interaction using r.db
    // ...
    return &models.User{ID: id, Name: "Example", Email: "example@email.com"}, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
    // Database interaction using r.db
    // ...
    return nil
}
```

**4. Service Layer (Services):**

* Implement business logic.
* Use interfaces to define service contracts, enabling dependency injection.
* This layer will use repositories to access data.

```go
// internal/services/user_service.go
package services

import (
    "your-project/internal/models"
    "your-project/internal/repositories"
)

type UserService interface {
    GetUser(id int) (*models.User, error)
    CreateNewUser(user *models.User) error
    // ... other methods ...
}

type userService struct {
    userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}

func (s *userService) GetUser(id int) (*models.User, error) {
    return s.userRepo.GetUserByID(id)
}

func (s *userService) CreateNewUser(user *models.User) error {
    return s.userRepo.CreateUser(user)
}
```

**5. Controller Layer (Controllers):**

* Handle HTTP requests and responses.
* Use services to perform business logic.
* Keep controllers thin and focused on request/response handling.

```go
// internal/controllers/user_controller.go
package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "your-project/internal/services"
)

type UserController struct {
    userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
    return &UserController{userService: userService}
}

func (c *UserController) GetUser(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    user, err := c.userService.GetUser(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, user)
}

// ... other controller methods ...
```

**6. Database Package (`pkg/database`):**

* Encapsulate database connection and interaction logic.
* Provide an interface for database operations.

```go
// pkg/database/database.go
package database

import "database/sql"

type DB struct {
    Conn *sql.DB
}

func NewDB(dsn string) (*DB, error) {
    conn, err := sql.Open("mysql", dsn) // Replace with your DB driver
    if err != nil {
        return nil, err
    }
    return &DB{Conn: conn}, nil
}
```

**7. Main Application (`cmd/api/main.go`):**

* Initialize dependencies and configure Gin.
* Set up routes and connect controllers.

```go
// cmd/api/main.go
package main

import (
    "github.com/gin-gonic/gin"
    "your-project/internal/controllers"
    "your-project/internal/repositories"
    "your-project/internal/services"
    "your-project/pkg/database"
)

func main() {
    db, err := database.NewDB("your-dsn")
    if err != nil {
        panic(err)
    }
    userRepo := repositories.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userController := controllers.NewUserController(userService)

    router := gin.Default()
    router.GET("/users/:id", userController.GetUser)

    router.Run(":8080")
}
```

**Key Principles to Avoid Tight Coupling:**

* **Dependency Injection:** Inject dependencies (repositories, services) into controllers and services using interfaces.
* **Interfaces:** Define interfaces for repositories and services to abstract away concrete implementations.
* **Layered Architecture:** Separate concerns into distinct layers (models, repositories, services, controllers).
* **Avoid Direct Database Access in Controllers:** Controllers should only interact with services.
* **DTOs (Data Transfer Objects):** If needed, use DTOs to decouple models from API responses.
* **Error Handling:** Implement robust error handling to prevent cascading failures.
* **Configuration:** Use environment variables or configuration files for database connection strings and other settings.

This architecture provides a good balance between maintainability, testability, and scalability. Remember to adapt the example to your specific needs and project requirements.
