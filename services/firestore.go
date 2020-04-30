package services

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
)

// FirestoreClient describe a Google Cloud Firestore Client
type FirestoreClient struct {
	c *firestore.Client
}

// User describe an application user
type User struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Users describe multiple users
type Users []*User

var s *FirestoreClient

const (
	usersCollection string = "users"
)

// NewFirestoreClient declare a new Firestore Client
func NewFirestoreClient() (*FirestoreClient, error) {
	if s != nil {
		return s, nil
	}

	client, err := firestore.NewClient(context.Background(), os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		return nil, err
	}

	s = &FirestoreClient{client}
	return s, nil
}

// GetUsers return all saved users in database
func (s *FirestoreClient) GetUsers() (Users, error) {
	// List documents in collection
	data, err := s.c.Collection(usersCollection).Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}

	// Parse data into our struct
	var users Users = make(Users, 0)
	for _, u := range data {
		var user User
		err := mapstructure.Decode(u.Data(), &user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
