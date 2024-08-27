package Domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp time.Time          `bson:"timestamp"`
	UserID    string             `bson:"user_id"`
	Action    string             `bson:"action"`
	Details   string             `bson:"details"`
}
