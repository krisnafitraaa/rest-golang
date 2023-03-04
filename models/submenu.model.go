package models

type Submenu struct {
	SubmenuID  int    `json:"submenu_id"`
	MenuID     int    `json:"menu_id"`
	Submenu    string `json:"submenu"`
	Icon       string `json:"sm_icon"`
	Position   int    `json:"sm_position"`
	SubmenuURL string `json:submenu_url`
	IsActive   string `json:sm_is_active`
	CreatedAt  string `json:created_at`
	UpdatedAt  string `json:updated_at`
}
