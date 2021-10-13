package models

type UserTeamRel struct {
	UserID ID `json:"id_user,omitempty"`
	TeamID ID `json:"id_team,omitempty"`
}
