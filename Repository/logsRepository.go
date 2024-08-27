package Repository

import (
	"LoanTracker/Domain"
	"context"
	"time"

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
	
	opts := options.Find().SetSort(map[string]int{"timestamp": -1}) // Sort logs by newest first
	cursor, err := r.collection.Find(ctx, map[string]interface{}{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var logs []Domain.SystemLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
