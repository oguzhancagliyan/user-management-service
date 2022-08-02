package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddUserModel struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddUserResponseModel struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserModel struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UpdateUserResponseModel struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type DeleteUserModel struct {
	Id string `json:"id"`
}

type GetUserModel struct {
	Id string `json:"id"`
}

type GetUserResponseModel struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}

type ErrorModel struct {
	Error      string `json:"error"`
	StatusCode int    `json:"-"`
}

type UserEntity struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `json:"name" bson:"Name"`
	Password string             `json:"password" bson:"Password"`
	Email    string             `json:"email" bson:"Email"`
}
