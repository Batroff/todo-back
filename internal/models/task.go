package models

type Task struct {
	ID       ID     `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Priority *uint  `json:"priority,omitempty"`
	UserID   ID     `json:"id_user,omitempty"`
	TeamID   *ID    `json:"id_team,omitempty"`
}

// NewTask returns new *Task, generated from params
func NewTask(title string, priority *uint, userID ID, teamID *ID) *Task {
	return &Task{
		ID:       NewID(),
		Title:    title,
		Priority: priority,
		UserID:   userID,
		TeamID:   teamID,
	}
}

type RequestTask struct {
	Title    *string `json:"title,omitempty"`
	Priority *uint   `json:"priority,omitempty"`
	UserID   *ID     `json:"id_user,omitempty"`
	TeamID   *ID     `json:"id_team,omitempty"`
}
