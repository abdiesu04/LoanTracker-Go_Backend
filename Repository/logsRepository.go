package Repository

import (
	"LoanTracker/Domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository interface {
	CreateLog(log *Domain.SystemLog) error
	GetLogs() ([]Domain.SystemLog, error)
}

type logRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(collection *mongo.Collection) LogRepository {
	return &logRepository{collection}
}

func (r *logRepository) CreateLog(log *Domain.SystemLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, log)
	return err
}

func (r *logRepository) GetLogs() ([]Domain.SystemLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set options to sort by the timestamp in descending order and limit to 100 logs
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}). // Sort by latest timestamp
		SetLimit(100) // Limit to the first 100 results

	cursor, err := r.collection.Find(ctx, bson.D{}, opts) // Use bson.D{} to represent an empty filter
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []Domain.SystemLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	fmt.Println(logs)
	return logs, nil

}

