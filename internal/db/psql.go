package db

import (
	l "DMS/internal/logger"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ID uuid.UUID

func (d ID) ToString() string {
	return uuid.UUID(d).String()
}
func (d *ID) FromString(strID string) error {
	id, err := uuid.Parse(strID)
	if err != nil {
		return err
	}
	*d = ID(id)
	return nil
}

// Convert input value to the ID data type. It's used by the GORM.
func (d *ID) Scan(value any) error {
	strVal, isStr := value.(string)
	if !isStr {
		return fmt.Errorf("cannot convert %v of type %T to db.ID", value, value)
	}
	err := d.FromString(strVal)
	return err
}

// Return the ID to in a suitable format for the database driver. In this case, it's a string.
func (d ID) Value() (driver.Value, error) {
	return d.ToString(), nil
}

type BaseModel struct {
	ID        ID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Disability int8

const (
	IsDisabled    Disability = 1
	IsNotDisabled Disability = 0
)

type User struct {
	BaseModel
	Name        string
	PhoneNumber string `gorm:"not null;unique"`
	IsDisabled  Disability
	// The user created the current user
	CreatedByID *ID   `gorm:"type:uuid"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID"`
	// Each user could have many job positions
	JobPosition []JobPosition `gorm:"foreignKey:UserID"`
	Session     []Session     `gorm:"foreignKey:UserID"`
}

type Event struct {
	BaseModel
	// event name
	Name string
	// The id of job position who created the event
	CreatedByID ID `gorm:"type:uuid;default:uuid_generate_v4();not null"`
	Description string
	Doc         []Doc `gorm:"foreignKey:EventID"`
}

type Doc struct {
	BaseModel
	// The id of job position who created the document
	CreatedByID ID `gorm:"type:uuid;default:uuid_generate_v4();not null"`
	// The id of event the document is for that
	EventID    ID `gorm:"type:uuid;default:uuid_generate_v4();not null"`
	Context    *string
	Multimedia *[]Multimedia `gorm:"foreignKey:DocID"`
}

type MediaType uint8

const (
	MediaImage MediaType = iota
	MediaVideo
	MediaAudio
)

type Multimedia struct {
	BaseModel
	DocID ID `gorm:"type:uuid;default:uuid_generate_v4();not null"`
	Type  MediaType
	// Full path and file name (contains type too)
	Src string
	// Just contains filename and its type
	FileName string
}

type JobPosition struct {
	BaseModel
	// ID of the user the JP is for that.
	UserID   ID `gorm:"type:uuid;not null"`
	Title    string
	RegionID ID `gorm:"type:uuid;default:uuid_generate_v4()"`
	// ID of parent job position the current job position is for that
	ParentID     *ID          `gorm:"type:uuid"`
	Parent       *JobPosition `gorm:"foreignKey:ParentID"`
	JPPermission JPPermission `gorm:"foreignKey:JpID"`
	Event        []Event      `gorm:"foreignKey:CreatedByID"`
	Doc          []Doc        `gorm:"foreignKey:CreatedByID"`
}

// Permissions of a job position
type JPPermission struct {
	BaseModel
	JpID            ID `gorm:"type:uuid;not null;unique"`
	IsAllowCreateJP bool
}

// Customize name of the table
// TableName overrides the default table name.
func (JPPermission) TableName() string {
	return "jp_permissions"
}

// Store details of each login by users
type Session struct {
	BaseModel
	UserID    ID     `gorm:"type:uuid;not null"`
	UserAgent string `gorm:"not null"`
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	ExpiredAt int64 `gorm:"not null"`
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	LastUsageAt int64
}

// A type alias for PostgreSQL database type
type PSQLDB = gorm.DB

type PsqlConnDetails struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       string
	// Maximum connection lifetime
	MaxConnLifetime time.Duration
	// Maximum number of idle connections
	MaxIdleConns int
	// Maximum number of open connections
	MAxOpenConns int
}

// Create a new PostgreSQL database instance. If occured any error during connecting
// to database, panic.
// After creating the instance, config its options and migrate schema to the database.
func NewPsqlConn(conn *PsqlConnDetails, doAutoMigrate bool, logger l.Logger) PSQLDB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tehran",
		conn.Host, conn.Username, conn.Password, conn.DB, conn.Port,
	)
	logger.Infof("Trying to connect to PSQL database \"%s\" ", conn.DB)
	var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Panicf(
			"Failed to create a connection to PostgreSQL database. DB-name: '%s', username: '%s', port: %d, host: '%s' (%s)",
			conn.DB, conn.Username, conn.Port, conn.Host, err,
		)
	}
	if doAutoMigrate {
		switch err := autoMigrate(db); err {
		case nil:
			logger.Info("Migrated schema to the database")
		default:
			logger.Errorf("Failed to migrate schema to the database. (%s)", err)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panicf("Failed to getting created database '%s' instance", conn.DB)
	}
	sqlDB.SetMaxIdleConns(conn.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conn.MAxOpenConns)
	sqlDB.SetConnMaxLifetime(conn.MaxConnLifetime)

	if result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"); result.Error != nil {
		logger.Panicf("Failed to create extension 'uuid-ossp' in database '%s'", conn.DB)
	}
	logger.Infof("Connected to database '%s' successfully.", conn.DB)
	return *db
}

// Migrate from schema to database and update the database scheme.
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Event{}, &Doc{}, &JobPosition{}, &JPPermission{},
		&Multimedia{}, &Session{})
}
