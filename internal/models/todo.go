package models

type Todo struct {
	ID       ID      `json:"id,omitempty"`
	Title    *string `json:"title,omitempty"`
	Text     string  `json:"text,omitempty"`
	Complete bool    `json:"complete,omitempty"`
	TaskID   ID      `json:"id_task,omitempty"`
}
