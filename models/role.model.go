package models

import "database/sql"

type Role struct {
	RoleID    string         `json:"role_id"`
	Role      string         `json:"role"`
	RoleGroup string         `json:"role_group"`
	AddedInfo sql.NullString `json:"added_info"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}
