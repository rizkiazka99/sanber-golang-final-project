package models

import "time"

type Item struct {
	Id          int          `json:"id"`
	ItemName    string       `json:"item_name"`
	Images      []ItemImages `json:"images"`
	Description string       `json:"desc,omitempty"`
	Price       int          `json:"price"`
	Stock       int          `json:"stock,omitempty"`
	CreatedBy   string       `json:"created_by,omitempty"`
	CreatedAt   *time.Time   `json:"created_at,omitempty"`
	ModifiedBy  string       `json:"modified_by,omitempty"`
	ModifiedAt  *time.Time   `json:"modified_at,omitempty"`
}

type ItemImages struct {
	Id       int    `json:"id"`
	ItemId   int    `json:"item_id"`
	ImageUrl string `json:"image_url"`
}
