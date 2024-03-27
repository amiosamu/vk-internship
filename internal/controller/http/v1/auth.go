package v1

import (
	"net/http"

	"github.com/amiosamu/vk-internship/internal/service"
	"github.com/amiosamu/vk-internship/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authRoutes struct {
	authService service.Auth
}

type signUpInput struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type signInInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func newAuthRoutes(router *gin.RouterGroup, authService service.Auth) {
	r := &authRoutes{
		authService: authService,
	}

	router.POST("/sign-up", r.signUp)
	router.POST("/sign-in", r.signIn)
}

// @Summary Sign Up
// @Description User sign up
// @Tags auth
// @Accept json
// @Param input body signUpInput true "input"
// @Success 201 {object} v1.authRoutes.signUp.response
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/sign-up [post]
func (r *authRoutes) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if !validator.ValidateName(input.Name) {
		newErrorResponse(c, http.StatusBadRequest, "invalid name")
		return
	}
	if !validator.ValidateSurname(input.Surname) {
		newErrorResponse(c, http.StatusBadRequest, "invalid surname")
		return
	}
	if !validator.ValidatePassword(input.Password) {
		newErrorResponse(c, http.StatusBadRequest, "invalid password")
		return
	}

	id, err := r.authService.RegisterUser(c.Request.Context(), service.AuthCreateUserInput{
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			newErrorResponse(c, http.StatusConflict, "user already exists")
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	type response struct {
		ID uuid.UUID `json:"id"`
	}

	c.JSON(http.StatusCreated, response{
		ID: id,
	})

}

// @Summary Sign in
// @Description Sign in
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signInInput true "input"
// @Success 200 {object} v1.authRoutes.signIn.response
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/sign-in [post]
func (r *authRoutes) signIn(c *gin.Context) {

	var input signInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := r.authService.GenerateToken(c.Request.Context(), service.AuthGenerateTokenInput{
		Email: input.Email,
	})

	if err != nil {
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.Header("Authorization", "Bearer "+token)

	type response struct {
		Token string `json:"token"`
	}

	c.JSON(http.StatusOK, response{
		Token: token,
	})
}
