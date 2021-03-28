package graph

import "go.mongodb.org/mongo-driver/mongo"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	database mongo.Database
}

func (r *Resolver) GetDatabase() mongo.Database {
	return r.database
}

func NewResolver(db mongo.Database) Resolver {
	return Resolver{db}
}
