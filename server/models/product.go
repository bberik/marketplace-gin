package models

type Product struct {
	ProductID   string    `json:"pID" bson:"pID"`
	ProductName *string   `json:"productName" bson:"productName"`
	Description *string   `json:"description" bson:"description"`
	Category    *string   `json:"category" bson:"category"`
	Images      []*string `json:"image" bson:"image"`
	Items       []*Item   `json:"items" bson:"items"`
}

type Item struct {
	ItemID   string   `json:"itemID" bson:"itemID"`
	Color    *string  `json:"color" bson:"color"`
	Size     *string  `json:"size" bson:"size"`
	InStock  *uint    `json:"instock" bson:"instock" validate:"required"`
	Price    *float64 `json:"price" bson:"price, omit" validate:"required"`
	Discount *float64 `json:"discount" bson:"discount"`
}

type Content struct {
	ProductID *string `json:"pID" bson:"pID" validate:"required"`
	ItemID    *string `json:"itemID" bson:"itemID" validate:"required"`
	Quantity  *uint   `json:"quantity" bson:"quantity" validate:"required"`
}
