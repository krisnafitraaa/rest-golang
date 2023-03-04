package models

import "database/sql"

type User struct {
	UserID       string         `json:"user_id" validate="required"`
	RoleID       string         `json:"role_id" validate="required"`
	VolunteerID  string         `json:"volunteer_id" validate="required"`
	Fullname     string         `json:"fullname" validate="required"`
	Email        string         `json:"email" validate="required,email"`
	Password     string         `json:"password" validate="required"`
	Photo        sql.NullString `json:"photo"`
	IsLoggedIn   int            `json:"is_loggedin"`
	LastLoggedIn sql.NullString `json:"last_loggedin"`
	Provider     string         `json:"provider" validate="required"`
	IsActive     int            `json:"is_active"`
	VerifiedUser int            `json:"verified_user"`
	CreatedBy    sql.NullString `json:"created_by"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}
