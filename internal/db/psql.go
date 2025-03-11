package db

import (
	l "DMS/internal/logger"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ID uint64

func (d *ID) ToInt64() int64 {
	return int64(*d)
}

func (d *ID) IsNull() bool {
	return *d == 0
}

type BaseModel struct {
	ID        ID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	IsDisabled    uint8 = 1
	IsNotDisabled uint8 = 0
)

type User struct {
	BaseModel
	Name        string
	PhoneNumber string `gorm:"not null;unique"`
	// value 0 means it's not disabled and value 1 means it's disabled.
	IsDisabled uint8
	// The user created the current user
	CreatedByID *ID
	CreatedBy   *User `gorm:"foreignKey:CreatedByID"`
	// Each user could have many job positions
	JobPosition []JobPosition `gorm:"foreignKey:UserID"`
}

type Event struct {
	BaseModel
	// event name
	Name string
	// The id of job position who created the document
	CreatedByID ID `gorm:"not null"`
	Description string
	Doc         []Doc `gorm:"foreignKey:EventID"`
}

type Doc struct {
	BaseModel
	// The id of job position who created the document
	CreatedByID ID `gorm:"not null"`
	// The id of event the document is for that
	EventID    ID `gorm:"not null"`
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
	DocID ID `gorm:"not null"`
	Type  MediaType
	// Full path and file name (contains type too)
	Src string
	// Just contains filename and its type
	FileName string
}

type JobPosition struct {
	BaseModel
	// ID of the user the JP is for that.
	UserID   ID `gorm:"not null"`
	Title    string
	RegionID ID
	// ID of parent job position the current job position is for that
	ParentID      *ID
	Parent        *JobPosition  `gorm:"foreignKey:ParentID"`
	JPPermission  JPPermission  `gorm:"foreignKey:JPID"`
	Event         []Event       `gorm:"foreignKey:CreatedByID"`
	Doc           []Doc         `gorm:"foreignKey:CreatedByID"`
	HierarchyTree HierarchyTree `gorm:"foreignKey:JPID"`
}

// Permissions of a job position
type JPPermission struct {
	BaseModel
	JPID            ID `gorm:"not null;unique"`
	IsAllowCreateJP bool
}

// Customize name of the table
// TableName overrides the default table name.
func (JPPermission) TableName() string {
	return "jp_permissions"
}

// For each job position we have a row that contains job position id and its child
// job positions.
type HierarchyTree struct {
	BaseModel
	JPID ID `gorm:"not null;unique"`
	// TODO: Create foreign ket for elements of this array
	ChildJPsID *[]ID
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

	logger.Infof("Connected to database '%s' successfully.", conn.DB)
	return *db
}

// Migrate from schema to database and update the database scheme.
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Event{}, &Doc{}, &JobPosition{}, &JPPermission{})
}
