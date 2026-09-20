package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Category        string             `json:"category" bson:"category"`
	PricePerKg      float64            `json:"price_per_kg" bson:"price_per_kg"`
	StockKg         int                `json:"stock_kg" bson:"stock_kg"`
	Unit            string             `json:"unit" bson:"unit"`
	OriginRegion    string             `json:"origin_region" bson:"origin_region"`
	FarmerID        int                `json:"farmer_id,omitempty" bson:"farmer_id,omitempty"`
	FarmerName      string             `json:"farmer_name" bson:"farmer_name"`
	FarmerAvatarURL string             `json:"farmer_avatar_url,omitempty" bson:"farmer_avatar_url,omitempty"`
	IsOrganic       bool               `json:"is_organic" bson:"is_organic"`
	Description     string             `json:"description" bson:"description"`
	ImageURL        string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}

type CreateProductDTO struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	PricePerKg      float64 `json:"price_per_kg"`
	StockKg         int     `json:"stock_kg"`
	Unit            string  `json:"unit"`
	OriginRegion    string  `json:"origin_region"`
	FarmerID        int     `json:"farmer_id,omitempty"`
	FarmerName      string  `json:"farmer_name"`
	FarmerAvatarURL string  `json:"farmer_avatar_url,omitempty"`
	IsOrganic       bool    `json:"is_organic"`
	Description     string  `json:"description"`
	ImageURL        string  `json:"image_url,omitempty"`
}

type UpdateProductDTO struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	PricePerKg      float64 `json:"price_per_kg"`
	StockKg         int     `json:"stock_kg"`
	Unit            string  `json:"unit"`
	OriginRegion    string  `json:"origin_region"`
	IsOrganic       bool    `json:"is_organic"`
	Description     string  `json:"description"`
	ImageURL        string  `json:"image_url,omitempty"`
}

type StockUpdateDTO struct {
	Quantity int `json:"quantity"`
}
