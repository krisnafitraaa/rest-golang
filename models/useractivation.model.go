package models

import "database/sql"

type UserActivation struct {
	ID              int            `json:"id"`
	UserID          string         `json:"user_id"`
	ActivationToken string         `json:"activation_token"`
	ValidUntil      string         `json:"valid_until"`
	ActivatedAt     sql.NullString `json:"activated_at"`
	IpAddress       sql.NullString `json:"ip_address"`
	UserAgent       sql.NullString `json:"user_agent"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}
