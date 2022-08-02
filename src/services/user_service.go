package services

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"user-management-service/src/helpers"
	"user-management-service/src/models"
	"user-management-service/src/validators"
)

type IUserService interface {
	AddUser(context context.Context, model models.AddUserModel) (responseModel models.
		AddUserResponseModel,
		errorModel *models.ErrorModel)
	UpdateUser(context context.Context, model models.UpdateUserModel) (responseModel models.
		UpdateUserResponseModel,
		errorModel *models.ErrorModel)
	GetAllUsers(ctx context.Context) (responseModel []models.GetUserResponseModel,
		errorModel *models.ErrorModel)
	GetUser(context context.Context, model models.GetUserModel) (responseModel models.
		GetUserResponseModel,
		errorModel *models.ErrorModel)
	DeleteUser(context context.Context, model models.DeleteUserModel) (
		errorModel *models.ErrorModel)
}

type UserService struct {
	validator validators.IUserValidator
	logger    *logrus.Logger
}

func NewUserService(validator validators.IUserValidator, logger *logrus.Logger) *UserService {
	return &UserService{validator: validator, logger: logger}
}

func (c *UserService) AddUser(context context.Context, model models.AddUserModel) (responseModel models.
	AddUserResponseModel,
	errorModel *models.ErrorModel) {

	error := c.validator.ValidateAddUserModel(model)

	if error != nil {
		return responseModel, error
	}

	var user models.UserEntity

	err := helpers.UserCollection.FindOne(context, bson.D{{"Email", model.Email}}).Decode(&user)

	if err == nil {
		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "AddUser").
			WithField("Operation", "CountDocuments").
			Warn("User already exist for the email")
		return responseModel, &models.ErrorModel{
			Error:      models.EmailExistMessage,
			StatusCode: http.StatusForbidden,
		}
	}

	if err != mongo.ErrNoDocuments {
		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "AddUser").
			WithField("Operation", "FindOne").
			WithField("Error", err.Error()).
			Error("")
		return responseModel, &models.ErrorModel{
			Error:      models.InternalErrorMessage,
			StatusCode: http.StatusInternalServerError,
		}
	}

	userEntity := models.UserEntity{
		Id:       primitive.NewObjectID(),
		Name:     model.Name,
		Password: model.Password,
		Email:    model.Email,
	}

	_, err = helpers.UserCollection.InsertOne(context, userEntity)

	if err != nil {
		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "AddUser").
			WithField("Operation", "InsertOne").
			WithField("Error", err.Error()).
			Error("")
		return responseModel, &models.ErrorModel{
			Error:      models.InternalErrorMessage,
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.logger.
		WithField("RequestModel", model).
		WithField("Service", "UserService").
		WithField("Method", "AddUser").
		Info("User Created")

	resp := models.AddUserResponseModel{
		Id:    userEntity.Id.Hex(),
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}

	return resp, nil
}

func (c *UserService) UpdateUser(context context.Context, model models.UpdateUserModel) (responseModel models.
	UpdateUserResponseModel,
	errorModel *models.ErrorModel) {
	error := c.validator.ValidateUpdateUserModel(model)

	if error != nil {
		return responseModel, error
	}

	var userEntity models.UserEntity

	objID, _ := primitive.ObjectIDFromHex(model.Id)

	err := helpers.UserCollection.FindOne(context, bson.D{{"_id", objID}}).Decode(&userEntity)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.logger.
				WithField("RequestModel", model).
				WithField("Service", "UserService").
				WithField("Method", "UpdateUser").
				WithField("Operation", "FindOne").
				Warn("UserNotFound")
			return responseModel, &models.ErrorModel{
				Error:      models.UserNotFoundErrorMessage,
				StatusCode: http.StatusNotFound,
			}
		}

		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "UpdateUser").
			WithField("Operation", "FindOne").
			WithField("Error", err.Error()).
			Error("")
		return responseModel, &models.ErrorModel{
			Error:      models.InternalErrorMessage,
			StatusCode: http.StatusInternalServerError,
		}
	}

	update := bson.D{{"$set",
		bson.D{
			{"Name", model.Name},
			{"Password", model.Password}}}}

	updateResult, err := helpers.UserCollection.UpdateByID(context, objID, update)

	if err != nil {
		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "UpdateUser").
			WithField("Operation", "UpdateByID").
			WithField("Error", err.Error()).
			Error("")
		return responseModel, &models.ErrorModel{
			Error:      models.InternalErrorMessage,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if updateResult.MatchedCount == 0 {
		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "UpdateUser").
			WithField("Operation", "UpdateByID").
			Warn("User not found")
		return responseModel, &models.ErrorModel{
			Error:      models.UserNotFoundErrorMessage,
			StatusCode: http.StatusNotFound,
		}
	}

	return models.UpdateUserResponseModel{
		Id:    model.Id,
		Email: userEntity.Email,
		Name:  model.Name,
	}, nil

}

func (c *UserService) DeleteUser(context context.Context, model models.DeleteUserModel) (
	errorModel *models.ErrorModel) {
	error := c.validator.ValidateDeleteUserModel(model)

	if error != nil {
		return error
	}

	deleteResult, err := helpers.UserCollection.DeleteOne(context, bson.D{{"_id", model.Id}})

	if err != nil {
		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "DeleteUser").
			WithField("Operation", "DeleteOne").
			WithField("Error", err.Error()).
			Error("")
		return &models.ErrorModel{
			Error:      models.InternalErrorMessage,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if deleteResult.DeletedCount == 0 {
		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "DeleteUser").
			WithField("Operation", "DeleteOne").
			Warn("User Not Found")
		return &models.ErrorModel{
			Error:      models.UserNotFoundErrorMessage,
			StatusCode: http.StatusNotFound,
		}
	}
	return nil
}

func (c *UserService) GetUser(context context.Context, model models.GetUserModel) (responseModel models.
	GetUserResponseModel,
	errorModel *models.ErrorModel) {

	error := c.validator.ValidateGetUserModel(model)

	if error != nil {
		return responseModel, error
	}

	var userEntity models.UserEntity

	objID, _ := primitive.ObjectIDFromHex(model.Id)

	err := helpers.UserCollection.FindOne(context, bson.D{{"_id", objID}}).Decode(&userEntity)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.logger.
				WithField("RequestModel", model).
				WithField("Service", "UserService").
				WithField("Method", "GetUser").
				WithField("Operation", "FindOne").
				Warn("UserNotFound")
			return responseModel, &models.ErrorModel{
				Error:      models.UserNotFoundErrorMessage,
				StatusCode: http.StatusNotFound,
			}
		}

		c.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserService").
			WithField("Method", "GetUser").
			WithField("Operation", "FindOne").
			WithField("Error", err.Error()).
			Error("")
		return responseModel, &models.ErrorModel{
			Error:      models.InternalErrorMessage,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return models.GetUserResponseModel{
		Id:    userEntity.Id,
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}, nil

}

func (c *UserService) GetAllUsers(ctx context.Context) (responseModel []models.GetUserResponseModel,
	errorModel *models.ErrorModel) {

	cursor, err := helpers.UserCollection.Find(context.TODO(), bson.D{})
	if err = cursor.All(ctx, &responseModel); err != nil {
		c.logger.
			WithField("Service", "UserService").
			WithField("Method", "GetAllUsers").
			WithField("Operation", "Find").
			WithField("Error", err.Error()).
			Error("")
		return nil, &models.ErrorModel{
			Error:      models.InternalErrorMessage,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return responseModel, nil
}
