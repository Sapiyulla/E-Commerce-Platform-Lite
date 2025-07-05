package repository

import (
	"context"
	dto "payment-service/internal/application/DTO"
	"payment-service/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	database *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{database: db}
}

func (r *MongoRepository) Pay(ctx context.Context, payment domain.Payment) dto.PaymentStatus {
	_, err := r.database.
		Collection("payments").
		InsertOne(ctx, payment)
	if err != nil {
		return dto.PaymentStatus{Status: "error", Message: err.Error()}
	}
	return dto.PaymentStatus{Status: "success", Message: payment.ID + " payed"}
}
