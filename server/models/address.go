package models

type Address struct {
	Details  string   `json:"details"  bson:"details"`
	Building string   `json:"building" bson:"building" validate:"required"`
	Street   string   `json:"street"   bson:"street" validate:"required"`
	Area     Location `json:"area"     bson:"area" validate:"required"`
}

type Location struct {
	City    string `json:"city"      bson:"city" validate:"required"`
	Country string `json:"country"   bson:"country" validate:"required"`
}
