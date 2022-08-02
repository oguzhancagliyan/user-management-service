package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"user-management-service/src/models"
	"user-management-service/src/services"
)

type UserController struct {
	userService services.IUserService
	logger      *logrus.Logger
}

func NewUserController(userService services.IUserService, logger *logrus.Logger) *UserController {
	return &UserController{userService: userService, logger: logger}
}

// AddUser godoc
// @Summary      AddUser
// @description  Adds the user
// @Tags         user
// @Success      200     {object}  models.AddUserResponseModel
// @Failure      400              {string}  string    "error"
// @Param        model  body    models.AddUserModel  true  "AddUserModel"
// @Router       /users [post]
func (c *UserController) AddUser(context *gin.Context) {
	var model models.AddUserModel
	err := context.ShouldBindJSON(&model)

	if err != nil {
		errorModel := models.ErrorModel{
			Error:      models.BadRequestErrorMessage,
			StatusCode: http.StatusBadRequest,
		}

		context.JSON(errorModel.StatusCode, errorModel.Error)
		return
	}
	result, error := c.userService.AddUser(context.Request.Context(), model)

	if error != nil {
		context.JSON(error.StatusCode, error.Error)
		return
	}

	context.JSON(http.StatusOK, result)
}

// UpdateUser godoc
// @Summary      UpdateUser
// @description  updates the user
// @Tags         user
// @Success      200     {object}  models.UpdateUserResponseModel
// @Failure      400              {string}  string    "error"
// @Param        model  body    models.UpdateUserModel  true  "UpdateUserModel"
// @Router       /users [patch]
func (c *UserController) UpdateUser(context *gin.Context) {
	var model models.UpdateUserModel
	err := context.ShouldBindJSON(&model)

	if err != nil {
		errorModel := models.ErrorModel{
			Error:      models.BadRequestErrorMessage,
			StatusCode: http.StatusBadRequest,
		}

		context.JSON(errorModel.StatusCode, errorModel.Error)
		return
	}
	result, error := c.userService.UpdateUser(context.Request.Context(), model)

	if error != nil {
		context.JSON(error.StatusCode, error.Error)
		return
	}

	context.JSON(http.StatusOK, result)
}

// DeleteUser godoc
// @Summary      DeleteUser
// @description  deletes the user
// @Tags         user
// @Success      200     {object}  models.AddUserResponseModel
// @Failure      400              {string}  string    "error"
// @Param        id   path      string  true  "id"
// @Router       /users/{id} [delete]
func (c *UserController) DeleteUser(context *gin.Context, id string) {
	model := models.DeleteUserModel{Id: id}

	errorModel := c.userService.DeleteUser(context.Request.Context(), model)

	if errorModel != nil {
		context.JSON(errorModel.StatusCode, errorModel.Error)
		return
	}

	context.JSON(http.StatusOK, nil)
}

// GetUser godoc
// @Summary      GetUser
// @description  retrieves the user
// @Tags         user
// @Success      200     {object}  models.GetUserResponseModel
// @Failure      400              {string}  string    "error"
// @Param        id   path      string  true  "id"
// @Router       /users/{id} [get]
func (c *UserController) GetUser(context *gin.Context, id string) {

	model := models.GetUserModel{Id: id}

	response, errorModel := c.userService.GetUser(context.Request.Context(), model)

	if errorModel != nil {
		context.JSON(errorModel.StatusCode, errorModel.Error)
		return
	}

	context.JSON(http.StatusOK, response)
}

// GetAllUser godoc
// @Summary      GetAllUser
// @description  retrieves the user
// @Tags         user
// @Success      200     {object}  []models.GetUserResponseModel
// @Failure      400              {string}  string    "error"
// @Router       /users [get]
func (c *UserController) GetAllUser(context *gin.Context) {
	response, errorModel := c.userService.GetAllUsers(context)

	if errorModel != nil {
		context.JSON(errorModel.StatusCode, errorModel.Error)
		return
	}

	context.JSON(http.StatusOK, response)
}
