package models

import (
	"time"
)

type Cart struct {
	Id            int        `json:"id"`
	UserId        int        `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	TotalPrice    int        `json:"total_price"`
	CartItems     []CartItem `json:"items"`
	PaymentMethod string     `json:"payment_method"`
	PaymentStatus string     `json:"payment_status"`
}

type CartItem struct {
	Id       int   `json:"id"`
	CartId   int   `json:"cart_id"`
	ItemId   int   `json:"item_id"`
	Quantity int   `json:"quantity"`
	Item     *Item `json:"item"`
}

type PostCartBody struct {
	Id            int        `json:"id"`
	UserId        int        `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	Items         []CartItem `json:"items"`
	TotalPrice    int        `json:"total_price"`
	PaymentMethod string     `json:"payment_method"`
	PaymentStatus string     `json:"payment_status"`
}

type CartResponse struct {
	Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
}

type CartItemUpdate struct {
	ItemId   int64 `json:"item_id"`
	Quantity int   `json:"quantity"`
}

type CartPayment struct {
	PaymentToken string `json:"payment_token"`
}
