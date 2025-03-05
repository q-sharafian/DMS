package db

import (
	l "DMS/internal/logger"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ID uint64

type model struct {
	ID        ID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	model
	Name        string
	PhoneNumber string
	// value 0 means it's not disabled and value 1 means it's disabled.
	IsDisabled uint8
	// The user created the current user
	CreatedByID ID `json:"created_by_id"`
}

type Event struct {
	model
	// event name
	Name        string
	CreatedByID ID
	Description string
}

type Doc struct {
	model
	CreatedByID ID
	EventID     ID
	Context     string
}

type MediaType uint8

const (
	Image MediaType = iota
	Video
	Audio
)

type Multimedia struct {
	model
	DocID ID
	Type  MediaType
	// Full path and file name (contains type too)
	Src string
	// Just contains filename and its type
	FileName string
}

// A type alias for PostgreSQL database type
type PSQLDB = gorm.DB

type PsqlConnDetails struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       string
}

// Create a new PostgreSQL database instance. If occured any error during connecting
// to database, panic.
func NewPsqlConn(conn *PsqlConnDetails, logger *l.Logger) PSQLDB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tehran",
		conn.Host, conn.Username, conn.Password, conn.DB, conn.Port,
	)
	var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(
			fmt.Sprintf(
				"Failed to create a connection to PostgreSQL database. DB-name: '%s', username: '%s', port: %d, host: '%s'\nError: %s\n",
				conn.DB, conn.Username, conn.Port, conn.Host, err,
			),
		)
	}
	return *db
}
