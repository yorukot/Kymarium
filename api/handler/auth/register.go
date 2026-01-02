package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

// +----------------------------------------------+
// | Register                                     |
// +----------------------------------------------+

type registerRequest struct {
	Email       string `json:"email" validate:"required,email,max=255" example:"user@example.com"`
	Password    string `json:"password" validate:"required,min=8,max=255" example:"password123"`
	DisplayName string `json:"display_name" validate:"required,min=3,max=255" example:"John Doe"`
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with email and password; sends verification email when SMTP is enabled
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Registration request"
// @Success 200 {object} response.SuccessResponse "Registration successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or email already in use"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	// Decode the request body
	var registerRequest registerRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate the request body
	if err := validator.New().Struct(registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Begin the transaction
	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}

	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	// Get the account by email
	checkedAccount, err := h.Repo.GetAccountByEmail(c.Request().Context(), tx, registerRequest.Email)
	if err != nil {
		zap.L().Error("Failed to check if user already exists", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if user already exists")
	}

	// If the account is found, return an error
	if checkedAccount != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "This email is already in use")
	}

	// Generate the user and account
	user, account, err := GenerateUser(registerRequest)
	if err != nil {
		zap.L().Error("Failed to generate user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate user")
	}

	smtpEnabled := config.Env().SMTPEnabled
	var verifyToken string
	if smtpEnabled {
		verifyToken, err = generateEmailVerificationToken(user.ID, registerRequest.Email)
		if err != nil {
			zap.L().Error("Failed to generate verification token", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate verification token")
		}
		user.Verified = false
		user.VerifyCode = &verifyToken
	} else {
		user.Verified = true
		user.VerifyCode = nil
	}

	// Create the user and account in the database
	if err = h.Repo.CreateUserAndAccount(c.Request().Context(), tx, user, account); err != nil {
		zap.L().Error("Failed to create user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	if smtpEnabled {
		if err := sendVerificationEmail(registerRequest.Email, verifyToken); err != nil {
			zap.L().Error("Failed to send verification email", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send verification email")
		}

		if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
			zap.L().Error("Failed to commit transaction", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
		}

		return c.JSON(http.StatusOK, response.SuccessMessage("Verification email sent"))
	}

	// Generate the refresh token
	refreshToken, err := generateTokenAndSaveRefreshToken(c, h.Repo, tx, user.ID)
	if err != nil {
		zap.L().Error("Failed to generate refresh token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token")
	}

	accessTokenCookie, err := generateAccessTokenCookieForUser(user.ID)
	if err != nil {
		zap.L().Error("Failed to generate access token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
	}

	// Commit the transaction
	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	// Generate the refresh token cookie
	refreshTokenCookie := generateRefreshTokenCookie(refreshToken)
	c.SetCookie(&refreshTokenCookie)
	c.SetCookie(&accessTokenCookie)

	// Respond with the success message
	return c.JSON(http.StatusOK, response.SuccessMessage("User registered successfully"))
}
