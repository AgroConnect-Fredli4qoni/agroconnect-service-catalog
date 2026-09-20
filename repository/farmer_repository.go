package repository

import (
	"context"
	"errors"
	"time"

	"github.com/agroconnect/service-catalog/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FarmerRepository interface {
	FindAll(ctx context.Context) ([]models.Farmer, error)
	FindBySlug(ctx context.Context, slug string) (*models.Farmer, error)
	Create(ctx context.Context, farmer *models.Farmer) error
}

type mongoFarmerRepository struct {
	collection *mongo.Collection
}

func NewFarmerRepository(db *mongo.Database) FarmerRepository {
	return &mongoFarmerRepository{
		collection: db.Collection("farmers"),
	}
}

func (r *mongoFarmerRepository) FindAll(ctx context.Context) ([]models.Farmer, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var farmers []models.Farmer
	for cursor.Next(ctx) {
		var f models.Farmer
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		farmers = append(farmers, f)
	}

	if farmers == nil {
		farmers = []models.Farmer{}
	}

	return farmers, nil
}

func (r *mongoFarmerRepository) FindBySlug(ctx context.Context, slug string) (*models.Farmer, error) {
	var farmer models.Farmer
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&farmer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("farmer not found")
		}
		return nil, err
	}

	return &farmer, nil
}

func (r *mongoFarmerRepository) Create(ctx context.Context, farmer *models.Farmer) error {
	farmer.ID = primitive.NewObjectID()
	farmer.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, farmer)
	return err
}
