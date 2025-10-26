package v1handler

import (
	"net/http"

	v1dto "dnk.com/hoc-golang/internal/dto/v1"
	v1service "dnk.com/hoc-golang/internal/service/v1"
	"dnk.com/hoc-golang/internal/utils"
	"dnk.com/hoc-golang/internal/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1service.AuthService
}

func NewAuthHandler(service v1service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}

}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var input v1dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		// ctx.JSON(http.StatusBadRequest, validation.HandleValidationErrors(err))
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiresIn, err := ah.service.Login(ctx, input.Email, input.Password)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	response := v1dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Login succcessfully", response)
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	var input v1dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {

		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := ah.service.Logout(ctx, input.RefreshToken); err != nil {
		utils.ResponseError(ctx, err)
		return

	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Logout succcessfully")
}
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {

	var input v1dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {

		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiresIn, err := ah.service.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	response := v1dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Refresh token generate succcessfully", response)
}

func (ah *AuthHandler) RequestForgotPassword(ctx *gin.Context) {
	var input v1dto.RefreshPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {

		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	
	err := ah.service.RequestForgotPassword(ctx, input.Email)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Reset Link sent to email")
}

func (ah *AuthHandler) ResetPassword(ctx *gin.Context) {
	var input v1dto.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {

		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	err := ah.service.ResetPassword(ctx, input.Token, input.NewPassword)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Password reset successfully")
}
