package handler

import (
	"fmt"
	"main/server/request"
	"main/server/response"
	"main/server/services/auth"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// GuestLoginService handles guest login and token generation.
//
// @Summary Guest Login
// @Description Perform guest login and generate access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param guestLoginRequest body request.GuestLoginRequest true "Guest Login Request"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /guest-login [post]
func GuestLoginHandler(ctx *gin.Context) {

	var guestLoginReq request.GuestLoginRequest
	err := utils.RequestDecoding(ctx, &guestLoginReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	err = guestLoginReq.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	auth.GuestLoginService(ctx, guestLoginReq)
}

// just for testing not required...
func AdminSignUpHandler() {
	auth.AdminSignUpService()
}

// @Summary Generates link
// @Description	Generates and send an email with reset password lin
// @Accept			json
// @Produce		json
// @Success		200			{object}	response.Success
// @Failure		400			{object}	response.Success
// @Param			adminEmail	body		request.ForgotPassRequest	true	"Admin registered email"
// @Tags			Authentication
// @Router			/forgot-password [post]
func ForgotPasswordHandler(ctx *gin.Context) {
	var forgotRequest request.ForgotPassRequest
	fmt.Println("origin is ", ctx.Request.Header.Get("Origin"))
	err := utils.RequestDecoding(ctx, &forgotRequest)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = forgotRequest.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)

		return
	}

	auth.ForgotPassService(ctx, forgotRequest)
}

// @Summary Reset's the password
// @Description	Reset password
// @Accept			json
// @Produce		json
// @Success		200			{object}	response.Success
// @Failure		400			{object}	response.Success
// @Param token query string true "Reset token"
// @Param			NewPassword	body		request.UpdatePasswordRequest	true	"Admins new password"
// @Tags			Authentication
// @Router			/reset-password [patch]
func ResetPasswordHandler(ctx *gin.Context) {
	tokenString := ctx.Request.URL.Query().Get("token")
	fmt.Println("token string is ", tokenString)
	var password request.ResetPasswordRequest

	err := utils.RequestDecoding(ctx, &password)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = password.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	auth.ResetPasswordService(ctx, tokenString, password)

}

// LoginService handles user login and token generation.
//
// @Summary Admin Login
// @Description Perform admin login and generate access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginDetails body request.AdminLoginRequest true "Login Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /login [post]
func LoginHandler(ctx *gin.Context) {
	var loginReq request.AdminLoginRequest
	err := utils.RequestDecoding(ctx, &loginReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = loginReq.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	fmt.Println("admin login")
	auth.AdminLoginService(ctx, loginReq)

}

// LoginService handles user login and token generation.
//
// @Summary Player Login
// @Description Perform user login and generate access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginDetails body request.PlayerLoginRequest true "Login Details"
// @Success 200 {object} response.Success "Login successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /player-login [post]
func PlayerLoginHandler(ctx *gin.Context) {
	var loginReq request.PlayerLoginRequest
	err := utils.RequestDecoding(ctx, &loginReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = loginReq.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	fmt.Println("player login")
	auth.PlayerLoginService(ctx, loginReq)

}

// UpdateEmailService updates the email of a player.
//
// @Summary Update Email
// @Description Update the email address of a player
// @Tags Player
// @Accept json
// @Produce json
// @Param email body request.UpdateEmailRequest true "Update Email Request"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} response.Success "Email updated successfully"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 401 {object} response.Success "Unauthorized"
// @Failure 404 {object} response.Success "Player not found"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /update-email [put]
func UpdateEmailHandler(ctx *gin.Context) {
	var email request.UpdateEmailRequest
	playerId, exists := ctx.Get(utils.PLAYERID)
	fmt.Println("player id is:", playerId, exists)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	err := utils.RequestDecoding(ctx, &email)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	err = email.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	auth.UpdateEmailService(ctx, email, playerId.(string))
}

// @Summary Updates the password of the admin
// @Description	Updates the password of the admin
// @Accept		json
// @Produce		json
// @Success		200			{object}	response.Success
// @Failure		400			{object}	response.Success
// @Failure		401			{object}	response.Success
// @Param Authorization header string true "Admin Access Token"
// @Param			newPassword	body		request.UpdatePasswordRequest	true	"New password of the player"
// @Tags			Authentication
// @Router			/update-pass [patch]
func UpdatePasswordHandler(ctx *gin.Context) {

	//get player id from the context that is passed from middleware
	playerID, exists := ctx.Get(utils.PLAYERID)
	fmt.Println("player id is :", playerID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	var password request.UpdatePasswordRequest
	err := utils.RequestDecoding(ctx, &password)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = password.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	auth.UpdatePasswordService(ctx, password, playerID.(string))
}

// GetAdminHandler retrieves the admin.
//
// @Summary Get Admins list
// @Description Retrieve the list of all admins
// @Tags			Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.Success "Admin Details fetched successfully"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /admin [get]
func GetAdminHandler(ctx *gin.Context) {
	auth.GetAdminService(ctx)
}

// @Summary Delete Player's account
// @Description	Delets the account of the player
// @Accept		json
// @Produce		json
// @Success		200			{object}	response.Success
// @Failure		400			{object}	response.Success
// @Failure		401			{object}	response.Success
// @Param Authorization header string true "Player Access Token"
// @Tags			Authentication
// @Router			/delete-account [delete]
func DeleteAccountHandler(ctx *gin.Context) {
	playerID, exists := ctx.Get(utils.PLAYERID)
	fmt.Println("player id is :", playerID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	auth.DeleteAccountService(ctx, playerID.(string))
}

// @Summary Logout Player
// @Description	Logs out a player
// @Accept			json
// @Produce		json
// @Param Authorization header string true "Player Access Token"
// @Success		200	{object}	response.Success
// @Failure		400	{object}	response.Success
// @Failure		401	{object}	response.Success
// @Tags			Authentication
// @Router			/logout [delete]
func LogutHandler(ctx *gin.Context) {
	playerID, exists := ctx.Get(utils.PLAYERID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	auth.LogoutService(ctx, playerID.(string))
}
