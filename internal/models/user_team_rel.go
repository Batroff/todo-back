package models

type UserTeamRel struct {
	UserID ID `json:"id_user,omitempty"`
	TeamID ID `json:"id_team,omitempty"`
}

func NewUserTeamRel(userID, teamID ID) *UserTeamRel {
	return &UserTeamRel{
		UserID: userID,
		TeamID: teamID,
	}
}
