package unit_tests

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-management-service/src/helpers"
	"user-management-service/src/models"
	"user-management-service/src/services"
	"user-management-service/src/validators"
)

func TestValidateAddUserModel_Should_Not_Validate(t *testing.T) {
	logger := log.New()
	validator := validators.NewUserValidator(logger)

	model := models.AddUserModel{
		Name:     "",
		Email:    "oguzhan@gmail.com",
		Password: "123",
	}

	result := validator.ValidateAddUserModel(model)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)

	model.Name = "oguzhan"
	model.Email = ""
	result = validator.ValidateAddUserModel(model)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)

	model.Email = "oguzhan@gmail.com"
	model.Password = ""
	result = validator.ValidateAddUserModel(model)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)
}

func TestValidateAddUserModel_Should_Validate(t *testing.T) {
	logger := log.New()
	validator := validators.NewUserValidator(logger)

	model := models.AddUserModel{
		Name:     "oguzhan",
		Email:    "oguzhan@gmail.com",
		Password: "123",
	}

	result := validator.ValidateAddUserModel(model)

	assert.Nil(t, result)
}

func TestValidateUpdateUserModel_Should_Not_Validate(t *testing.T) {
	logger := log.New()
	validator := validators.NewUserValidator(logger)

	model := models.UpdateUserModel{
		Id:       "12",
		Name:     "",
		Password: "123",
	}

	result := validator.ValidateUpdateUserModel(model)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)

	model.Name = "oguzhan"
	model.Id = ""
	result = validator.ValidateUpdateUserModel(model)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)

	model.Password = ""
	result = validator.ValidateUpdateUserModel(model)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)
}

func TestUpdateUser_Should_Return_UserNotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("user not found", func(mt *mtest.T) {
		model := models.UpdateUserModel{
			Name:     "oguzhan",
			Password: "123",
			Id:       primitive.NewObjectID().Hex(),
		}
		logger := log.New()
		validator := validators.NewUserValidator(logger)
		helpers.UserCollection = mt.Coll
		userService := services.NewUserService(validator, logger)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		first := mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch)

		mt.AddMockResponses(first)

		_, message := userService.UpdateUser(c, model)
		assert.NotNil(t, message)
		assert.Equal(t, models.UserNotFoundErrorMessage, message.Error)
	})
}

func TestValidateDeleteUserModel_Should_Not_Validate(t *testing.T) {
	logger := log.New()
	validator := validators.NewUserValidator(logger)

	model := models.DeleteUserModel{
		Id: "",
	}

	result := validator.ValidateDeleteUserModel(model)

	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)
}

func TestValidateGetUserModel_Should_Not_Validate(t *testing.T) {
	logger := log.New()
	validator := validators.NewUserValidator(logger)

	model := models.GetUserModel{
		Id: "",
	}

	result := validator.ValidateGetUserModel(model)

	assert.NotNil(t, result)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, models.BadRequestErrorMessage, result.Error)
}

func TestGetUser_Should_Return_UserNotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("user not found", func(mt *mtest.T) {
		model := models.GetUserModel{
			Id: primitive.NewObjectID().Hex(),
		}
		logger := log.New()
		validator := validators.NewUserValidator(logger)
		helpers.UserCollection = mt.Coll
		userService := services.NewUserService(validator, logger)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		first := mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch)

		mt.AddMockResponses(first)

		_, message := userService.GetUser(c, model)
		assert.NotNil(t, message)
		assert.Equal(t, models.UserNotFoundErrorMessage, message.Error)
	})
}
