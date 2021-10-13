package models

type Team struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RequestTeam struct {
	Name *string `json:"name,omitempty"`
}
