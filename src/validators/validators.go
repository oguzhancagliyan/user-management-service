package validators

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/mail"
	"user-management-service/src/models"
)

type IUserValidator interface {
	ValidateAddUserModel(model models.AddUserModel) *models.ErrorModel
	ValidateUpdateUserModel(model models.UpdateUserModel) *models.ErrorModel
	ValidateDeleteUserModel(model models.DeleteUserModel) *models.ErrorModel
	ValidateGetUserModel(model models.GetUserModel) *models.ErrorModel
}

type UserValidator struct {
	logger *logrus.Logger
}

func NewUserValidator(logger *logrus.Logger) *UserValidator {
	return &UserValidator{logger: logger}
}

func (v *UserValidator) ValidateAddUserModel(model models.AddUserModel) *models.ErrorModel {
	if model.Name == "" || model.Password == "" || model.Email == "" {
		v.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserValidator").
			WithField("Method", "ValidateAddUserModel").
			Warn("Name or Password or Email empty")
		return &models.ErrorModel{
			StatusCode: http.StatusBadRequest,
			Error:      models.BadRequestErrorMessage,
		}
	} else if model.Email != "" {
		_, err := mail.ParseAddress(model.Email)

		if err != nil {
			v.logger.
				WithField("RequestModel", model).
				WithField("Service", "UserValidator").
				WithField("Operation", "ParseAddress").
				WithField("Method", "ValidateAddUserModel").
				WithField("Error", err.Error()).
				Warn("Name or Password or Email empty")
			return &models.ErrorModel{
				StatusCode: http.StatusBadRequest,
				Error:      models.BadRequestErrorMessage,
			}
		}
	}
	return nil
}

func (v *UserValidator) ValidateUpdateUserModel(model models.UpdateUserModel) *models.ErrorModel {
	if model.Name == "" || model.Password == "" || model.Id == "" {
		v.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserValidator").
			WithField("Method", "ValidateUpdateUserModel").
			Warn("Name or Password or Id empty")
		return &models.ErrorModel{
			StatusCode: http.StatusBadRequest,
			Error:      models.BadRequestErrorMessage,
		}
	}
	return nil
}

func (v *UserValidator) ValidateDeleteUserModel(model models.DeleteUserModel) *models.ErrorModel {
	if model.Id == "" || model.Id == primitive.NilObjectID.Hex() {
		v.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserValidator").
			WithField("Method", "ValidateDeleteUserModel").
			Warn("Id is not valid or empty")
		return &models.ErrorModel{
			StatusCode: http.StatusBadRequest,
			Error:      models.BadRequestErrorMessage,
		}
	}
	return nil
}

func (v *UserValidator) ValidateGetUserModel(model models.GetUserModel) *models.ErrorModel {
	if model.Id == "" || model.Id == primitive.NilObjectID.Hex() {
		v.logger.
			WithField("RequestModel", model).
			WithField("Service", "UserValidator").
			WithField("Method", "ValidateGetUserModel").
			Warn("Id is not valid or empty")
		return &models.ErrorModel{
			StatusCode: http.StatusBadRequest,
			Error:      models.BadRequestErrorMessage,
		}
	}
	return nil
}
