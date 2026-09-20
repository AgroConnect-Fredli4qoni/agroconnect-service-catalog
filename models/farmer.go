package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Farmer struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Slug           string             `json:"slug" bson:"slug"`
	Name           string             `json:"name" bson:"name"`
	OriginRegion   string             `json:"origin_region" bson:"origin_region"`
	AvatarURL      string             `json:"avatar_url" bson:"avatar_url"`
	BannerURL      string             `json:"banner_url" bson:"banner_url"`
	Description    string             `json:"description" bson:"description"`
	Phone          string             `json:"phone" bson:"phone"`
	Address        string             `json:"address" bson:"address"`
	OperatingHours string             `json:"operating_hours" bson:"operating_hours"`
	LandArea       string             `json:"land_area" bson:"land_area"`
	IsVerified     bool               `json:"is_verified" bson:"is_verified"`
	FarmingMethods []string           `json:"farming_methods" bson:"farming_methods"`
	Certifications []string           `json:"certifications" bson:"certifications"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

type CreateFarmerDTO struct {
	Slug           string   `json:"slug"`
	Name           string   `json:"name"`
	OriginRegion   string   `json:"origin_region"`
	AvatarURL      string   `json:"avatar_url"`
	BannerURL      string   `json:"banner_url"`
	Description    string   `json:"description"`
	Phone          string   `json:"phone"`
	Address        string   `json:"address"`
	OperatingHours string   `json:"operating_hours"`
	LandArea       string   `json:"land_area"`
	IsVerified     bool     `json:"is_verified"`
	FarmingMethods []string `json:"farming_methods"`
	Certifications []string `json:"certifications"`
}
