package model

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions []Permission `json:"permissions,omitempty"`
}
