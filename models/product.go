package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Category     string             `json:"category" bson:"category"`
	PricePerKg   float64            `json:"price_per_kg" bson:"price_per_kg"`
	StockKg      int                `json:"stock_kg" bson:"stock_kg"`
	Unit         string             `json:"unit" bson:"unit"`
	OriginRegion string             `json:"origin_region" bson:"origin_region"`
	FarmerName   string             `json:"farmer_name" bson:"farmer_name"`
	IsOrganic    bool               `json:"is_organic" bson:"is_organic"`
	Description  string             `json:"description" bson:"description"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

type CreateProductDTO struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	PricePerKg   float64 `json:"price_per_kg"`
	StockKg      int     `json:"stock_kg"`
	Unit         string  `json:"unit"`
	OriginRegion string  `json:"origin_region"`
	FarmerName   string  `json:"farmer_name"`
	IsOrganic    bool    `json:"is_organic"`
	Description  string  `json:"description"`
}

type StockUpdateDTO struct {
	Quantity int `json:"quantity"`
}
