package models

// Contains details of a session
type Session struct {
	// Session id
	ID ID `json:"id"`
	// Unix time the session is created
	IssuedAt int64 `json:"issued_at"`
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	// Unix time the session expires
	ExpiredAt int64 `json:"expired_at"`
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	// This session is belong to the user id
	UserID ID `json:"user_id"`
	// Details of the device on which the user is logged in.
	UserAgent string `json:"user_agent"`
	// Last usage time of the session
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	// If it be 0 means it's not used.
	LastUsageAt int64
}

type JWT struct {
	UserID *ID `json:"sub"`
	JPID   *ID `json:"jp_id"`
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	IAT int64 `json:"iat"`
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	EXP int64 `json:"exp"`
	// JTI: (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed
	// (allows a token to be used only once). Each session has its own unique JTI.
	// It's stored as a Unix timestamp. (In seconds and UTC time zone)
	JTI *ID `json:"jti"`
}
