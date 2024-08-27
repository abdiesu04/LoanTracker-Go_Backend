package Repository

import (
	"LoanTracker/Domain"
	"errors"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(user *Domain.User) (*Domain.User, error)
	FindByEmail(email string) (*Domain.User, error)
	FindByUsername(username string) (*Domain.User, error)
	IsDbEmpty() (bool, error)
	Update(username string, updateFields bson.M) error
	ShowUsers() ([]Domain.User, error)
	FindByID(id string) (Domain.User, error)
	Delete(userID string) error
	// GetUserById(id string) (*Domain.User, error)
	// GetUserByUsername(username string) (*Domain.User, error)
	// GetAllUsers() ([]Domain.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection}
}


func (r *userRepository) FindByID(id string) (Domain.User, error) {
	var user Domain.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Domain.User{}, err
	}
	err = r.collection.FindOne(context.TODO(), bson.M{"id": objectID}).Decode(&user)
	if err != nil {
		return Domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(user *Domain.User) (*Domain.User, error) {
	ctx := context.TODO()
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*Domain.User, error) {
	ctx := context.TODO()
	var user Domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Return nil for both user and error if no document is found
		}
		return nil, err // Return the actual error if something else went wrong
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*Domain.User, error) {
	ctx := context.TODO()
	var user Domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Return nil for both user and error if no document is found
		}
		return nil, err // Return the actual error if something else went wrong
	}
	return &user, nil
}

func (r *userRepository) IsDbEmpty() (bool, error) {
	count, err := r.collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return false, err
	}
	return count == 0, nil
}



func (r *userRepository) Update(username string, updateFields bson.M) error {
	filter := bson.M{"username": username}

	// Only perform the update if there are fields to update
	if len(updateFields) == 0 {
		return nil // No update needed
	}

	_, err := r.collection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateFields})
	return err
}

func (r *userRepository) ShowUsers() ([]Domain.User, error) {

	filter := bson.M{"role": "user"}
	cur, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var users []Domain.User
	for cur.Next(context.TODO()) {
		var User Domain.User
		if err := cur.Decode(&User); err != nil {
			return nil, err
		}
		users = append(users, User)
	}

	return users, nil
}


func (r *userRepository) Delete(userID string) error{
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return nil
}