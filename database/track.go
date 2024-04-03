package database

type Track struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
}
