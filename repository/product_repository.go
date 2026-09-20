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

type ProductRepository interface {
	FindAll(ctx context.Context, search, category string) ([]models.Product, error)
	FindByID(ctx context.Context, id string) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, id string, product *models.Product) error
	Delete(ctx context.Context, id string) error
	DeductStock(ctx context.Context, id string, quantity int) error
}

type mongoProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &mongoProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *mongoProductRepository) FindAll(ctx context.Context, search, category string) ([]models.Product, error) {
	filter := bson.M{}

	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	if category != "" && category != "Semua" {
		filter["category"] = category
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	for cursor.Next(ctx) {
		var p models.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if products == nil {
		products = []models.Product{}
	}

	return products, nil
}

func (r *mongoProductRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product id format")
	}

	var product models.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *mongoProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *mongoProductRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product id format")
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("product not found to delete")
	}
	return nil
}

func (r *mongoProductRepository) DeductStock(ctx context.Context, id string, quantity int) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product id format")
	}

	filter := bson.M{
		"_id":      objID,
		"stock_kg": bson.M{"$gte": quantity},
	}
	update := bson.M{
		"$inc": bson.M{"stock_kg": -quantity},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("insufficient stock or product not found")
	}
	return nil
}

func (r *mongoProductRepository) Update(ctx context.Context, id string, product *models.Product) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product id format")
	}

	update := bson.M{
		"$set": bson.M{
			"name":          product.Name,
			"category":      product.Category,
			"price_per_kg":  product.PricePerKg,
			"stock_kg":      product.StockKg,
			"unit":          product.Unit,
			"origin_region": product.OriginRegion,
			"is_organic":    product.IsOrganic,
			"description":   product.Description,
			"image_url":     product.ImageURL,
		},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("product not found to update")
	}
	return nil
}

