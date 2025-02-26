To allow the **Repository** to use multiple types of databases (e.g., PostgreSQL, MySQL, or even in-memory databases), you can leverage **interfaces** in Go. By defining an interface for the data access operations, you can create different implementations of the **DAL** for each database type. The **Repository** will then depend on the interface rather than a specific implementation, making it database-agnostic.

Here’s how you can modify the example to support multiple database types:

---

### 1. **Define an Interface for Data Access**
Create an interface that defines the methods required for user data access. Both the **Repository** and the **DAL** implementations will depend on this interface.

```go
package dal

// User represents a user entity.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserDataAccess defines the interface for user data access operations.
type UserDataAccess interface {
	CreateUser(name, email string) error
	GetUserByID(id int) (*User, error)
}
```

---

### 2. **Implement the Interface for Different Databases**
Create separate implementations of the `UserDataAccess` interface for each database type.

#### PostgreSQL Implementation
```go
package dal

import (
	"database/sql"
	"fmt"
)

// PostgresUserDAL implements UserDataAccess for PostgreSQL.
type PostgresUserDAL struct {
	DB *sql.DB
}

// NewPostgresUserDAL creates a new PostgresUserDAL instance.
func NewPostgresUserDAL(db *sql.DB) *PostgresUserDAL {
	return &PostgresUserDAL{DB: db}
}

// CreateUser inserts a new user into the PostgreSQL database.
func (d *PostgresUserDAL) CreateUser(name, email string) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2)`
	_, err := d.DB.Exec(query, name, email)
	if err != nil {
		return fmt.Errorf("failed to create user in PostgreSQL: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by their ID from the PostgreSQL database.
func (d *PostgresUserDAL) GetUserByID(id int) (*User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`
	row := d.DB.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from PostgreSQL: %w", err)
	}
	return user, nil
}
```

#### In-Memory Implementation (for testing or prototyping)
```go
package dal

import (
	"errors"
	"sync"
)

// InMemoryUserDAL implements UserDataAccess for an in-memory database.
type InMemoryUserDAL struct {
	users map[int]*User
	mu    sync.Mutex
	nextID int
}

// NewInMemoryUserDAL creates a new InMemoryUserDAL instance.
func NewInMemoryUserDAL() *InMemoryUserDAL {
	return &InMemoryUserDAL{
		users: make(map[int]*User),
		nextID: 1,
	}
}

// CreateUser adds a new user to the in-memory database.
func (d *InMemoryUserDAL) CreateUser(name, email string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	user := &User{
		ID:    d.nextID,
		Name:  name,
		Email: email,
	}
	d.users[d.nextID] = user
	d.nextID++
	return nil
}

// GetUserByID retrieves a user by their ID from the in-memory database.
func (d *InMemoryUserDAL) GetUserByID(id int) (*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	user, exists := d.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
```

---

### 3. **Update the Repository to Use the Interface**
The **Repository** will now depend on the `UserDataAccess` interface instead of a specific implementation.

```go
package repository

import (
	"example.com/dal"
)

// UserRepository provides a domain-centric interface for user operations.
type UserRepository struct {
	userDAL dal.UserDataAccess
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(userDAL dal.UserDataAccess) *UserRepository {
	return &UserRepository{userDAL: userDAL}
}

// CreateUser creates a new user using the repository.
func (r *UserRepository) CreateUser(name, email string) error {
	return r.userDAL.CreateUser(name, email)
}

// GetUserByID retrieves a user by their ID using the repository.
func (r *UserRepository) GetUserByID(id int) (*dal.User, error) {
	return r.userDAL.GetUserByID(id)
}
```

---

### 4. **Usage Example**
You can now use the **Repository** with different database implementations.

```go
package main

import (
	"database/sql"
	"fmt"
	"example.com/dal"
	"example.com/repository"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Example 1: Use PostgreSQL
	db, err := sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	postgresDAL := dal.NewPostgresUserDAL(db)
	userRepo := repository.NewUserRepository(postgresDAL)

	// Create and retrieve a user
	_ = userRepo.CreateUser("John Doe", "john.doe@example.com")
	user, _ := userRepo.GetUserByID(1)
	fmt.Printf("User from PostgreSQL: %+v\n", user)

	// Example 2: Use In-Memory Database
	inMemoryDAL := dal.NewInMemoryUserDAL()
	userRepo = repository.NewUserRepository(inMemoryDAL)

	_ = userRepo.CreateUser("Jane Doe", "jane.doe@example.com")
	user, _ = userRepo.GetUserByID(1)
	fmt.Printf("User from In-Memory: %+v\n", user)
}
```

---

### Key Benefits:
1. **Flexibility**: The **Repository** can work with any database implementation that satisfies the `UserDataAccess` interface.
2. **Testability**: You can use an in-memory implementation for unit testing without needing a real database.
3. **Decoupling**: The business logic (Repository) is decoupled from the specific database implementation.

This approach is highly scalable and adheres to the **Dependency Inversion Principle** (one of the SOLID principles), making your code more modular and maintainable.