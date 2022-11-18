package models

type User struct {
	UserID      string    `json:"userID" bson:"userID"`
	Username    *string   `json:"username" validate:"required"`
	FullName    *string   `json:"fullName" validate:"required"`
	Email       *string   `json:"email" validate:"email,required"`
	UserAddress []Address `json:"userAddress" bson:"userAddress"`
	Cart        []Content `json:"cart" bson:"cart"`
	Orders      []Order   `json:"orders" bson:"orders"`
	Password    *string   `json:"password" validate:"required"`
}

type Shop struct {
	UserID      string    `json:"userId" bson:"userId"`
	ShopName    string    `json:"shopName" validate:"required"`
	ShopAddress Address   `json:"shopAddress" validate:"required"`
	Products    []*string `json:"products" bson:"products"`
	Orders      []*string `json:"orders" bson:"orders"`
}
