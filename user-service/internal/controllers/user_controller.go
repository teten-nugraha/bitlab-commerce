package controllers

import (
	"net/http"
	"time"

	"user-service/internal/services"
	"user-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
	timeout     time.Duration
}

func NewUserController(userService *services.UserService, timeout time.Duration) *UserController {
	return &UserController{
		userService: userService,
		timeout:     timeout,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := c.userService.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrEmailAlreadyInUse {
			status = http.StatusConflict
		}
		response.Error(ctx, status, err)
		return
	}

	response.Success(ctx, http.StatusCreated, gin.H{
		"user": user,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := c.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrUserNotFound || err == services.ErrInvalidPassword {
			status = http.StatusUnauthorized
		}
		response.Error(ctx, status, err)
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{
		"token": token,
	})
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil)
		return
	}

	user, err := c.userService.GetUser(ctx, userID.(string))
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrUserNotFound {
			status = http.StatusNotFound
		}
		response.Error(ctx, status, err)
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{
		"user": user,
	})
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil)
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := c.userService.UpdateUser(ctx, userID.(string), req.FirstName, req.LastName)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrUserNotFound {
			status = http.StatusNotFound
		}
		response.Error(ctx, status, err)
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{
		"user": user,
	})
}
