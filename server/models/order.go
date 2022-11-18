package models

import "time"

type Order struct {
	OrderID         string        `json:"orderID"         bson:"orderID"`
	OrderContent    Content       `json:"orderContent"    bson:"orderContent"`
	ShippingAddress Address       `json:"shippingAddress" bson:"shippingAddress" validate:"required"`
	OrderStatus     []OrderUpdate `json:"orderStatus"     bson:"orderStatus"`
}

type OrderUpdate struct {
	Status    string    `json:"status" bson:"status"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
