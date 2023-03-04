package models

type Menu struct {
	MenuID    int         `json:menu_id`
	Menu      string      `json:menu`
	Icon      string      `json:icon`
	Position  int         `json:position`
	MenuFor   string      `json:menu_for`
	MenuURL   string      `json:menu_url`
	Submenus  interface{} `json:submenus`
	IsActive  string      `json:is_active`
	CreatedAt string      `json:created_at`
	UpdatedAt string      `json:updated_at`
}
