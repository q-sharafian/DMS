package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"fmt"
)

type SessionDAL interface {
	// Create a login for specified user and return its id
	CreateSession(loginInfo *m.Session) (*m.ID, error)
	// Delete a session by sessionID
	DeleteSession(sessionID m.ID) error
	// Returns true if the id of the user who owns the specified session matches the claimed user id.
	IsMatchSessionUserID(sessionID, claimedUserID m.ID) (bool, error)
	// If both error and session be nil, means there's not any matched session or
	// it's disabled/deleted.
	GetSessionByID(sessionID m.ID) (*m.Session, error)
}

type psqlSessionDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlSessionDAL(db *db.PSQLDB, logger l.Logger) *psqlSessionDAL {
	return &psqlSessionDAL{db, logger}
}

func (p *psqlSessionDAL) CreateSession(loginInfo *m.Session) (*m.ID, error) {
	session := db.Session{
		UserID:    *modelID2DBID(&loginInfo.UserID),
		UserAgent: loginInfo.UserAgent,
	}
	result := p.db.Create(&session)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create session for userID %s (%s)", session.UserID.ToString(), result.Error)
	}
	return dbID2ModelID(&session.ID), nil
}

func (p *psqlSessionDAL) DeleteSession(userID m.ID) error {
	result := p.db.Where(&db.Session{
		BaseModel: db.BaseModel{ID: *modelID2DBID(&userID)}}).
		Delete(&db.Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete session for userID %s (%s)", userID.ToString(), result.Error)
	}
	return nil
}

func (p *psqlSessionDAL) IsMatchSessionUserID(sessionID, claimedUserID m.ID) (bool, error) {
	var session db.Session
	result := p.db.Where(&db.Session{
		BaseModel: db.BaseModel{ID: *modelID2DBID(&sessionID)}}).
		Find(&session)
	if result.Error != nil {
		return false, fmt.Errorf("failed to get session by id %s (%s)", sessionID.ToString(), result.Error)
	}
	return session.UserID == *modelID2DBID(&claimedUserID), nil
}

func (p *psqlSessionDAL) GetSessionByID(sessionID m.ID) (*m.Session, error) {
	var session db.Session
	result := p.db.Where(&db.Session{
		BaseModel: db.BaseModel{ID: *modelID2DBID(&sessionID)}}).
		Find(&session)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get session by id %s (%s)", sessionID.ToString(), result.Error)
	} else if result.RowsAffected == 0 {
		return nil, nil
	}
	return &m.Session{
		ID:          *dbID2ModelID(&session.ID),
		UserID:      *dbID2ModelID(&session.UserID),
		UserAgent:   session.UserAgent,
		IssuedAt:    session.CreatedAt.Unix(),
		ExpiredAt:   session.ExpiredAt,
		LastUsageAt: session.LastUsageAt,
	}, nil
}
