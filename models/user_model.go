package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PetSize struct{
		Height	float64	`json:"height,omitempty"`
		Weight	float64	`json:"weight,omitempty"`
}

type Pet struct {
    Id       primitive.ObjectID `json:"id,omitempty"`
		Name		string							`json:"name,omitempty" validate:"required"`
    DOB     string             	`json:"dob,omitempty" validate:"required"`
    Owner 	string             	`json:"owner,omitempty" validate:"required"`
		Animal	string							`json:"animal,omitempty" validate:"required"`
		Breed		string							`json:"breed,omitempty"`
    Size   	PetSize             `json:"size,omitempty"`
		Toy			string							`json:"toy,omitempty" validate:"required"`
}