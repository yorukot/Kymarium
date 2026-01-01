package auth

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	"go.uber.org/zap"
)

// +----------------------------------------------+
// | OAuth Callback                               |
// +----------------------------------------------+

// OAuthCallback godoc
// @Summary OAuth callback handler
// @Description Handles OAuth provider callback, processes authorization code, creates/links user accounts, and issues authentication tokens
// @Tags oauth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth provider (e.g., google, github)"
// @Param code query string true "Authorization code from OAuth provider"
// @Param state query string true "OAuth state parameter for CSRF protection"
// @Success 307 {string} string "Redirect to success URL with authentication cookies set"
// @Failure 400 {object} response.ErrorResponse "Invalid provider, oauth state, or verification failed"
// @Failure 500 {object} response.ErrorResponse "Internal server error during user creation or token generation"
// @Router /auth/oauth/{provider}/callback [get]
func (h *Handler) OAuthCallback(c echo.Context) error {
	// Get the oauth state from the query params
	oauthState := c.QueryParam("state")
	code := c.QueryParam("code")

	// Get the oauth session cookie
	oauthSessionCookie, err := c.Cookie(models.CookieNameOAuthSession)
	if err != nil {
		zap.L().Debug("OAuth session cookie not found", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "OAuth session not found")
	}

	// Parse the provider
	provider, err := parseProvider(c.Param("provider"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid provider")
	}

	// No need to check if the provider is valid because it's checked in the parseProvider function
	oauthConfig := h.OAuthConfig.Providers[provider]
	// Get the oidc provider
	oidcProvider := h.OAuthConfig.OIDCProviders[provider]

	// Validate the oauth state
	valid, payload, err := oauthValidateStateWithPayload(oauthSessionCookie.Value)
	if err != nil || !valid || oauthState != payload.State {
		zap.L().Warn("OAuth state validation failed",
			zap.String("ip", c.RealIP()),
			zap.String("user_agent", c.Request().UserAgent()),
			zap.String("provider", string(provider)),
			zap.String("oauth_state", oauthState),
			zap.String("payload_state", payload.State))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid oauth state")
	}

	// Get the user ID from the session cookie
	var userID int64
	var accountID int64
	if payload.Subject != "" {
		userID, err = strconv.ParseInt(payload.Subject, 10, 64)
		if err != nil {
			zap.L().Error("Failed to parse user ID", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID in session")
		}
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	// Exchange the code for a token
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to exchange code")
	}

	// Get the raw ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		zap.L().Error("Failed to get id token from oauth response", zap.Any("token", token))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get id token")
	}

	// Verify the token
	userInfo, err := oauthVerifyTokenAndGetUserInfo(c.Request().Context(), rawIDToken, token, oidcProvider, oauthConfig)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to verify token")
	}

	// Begin the transaction
	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction", err)
	}

	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	// Get the account and user by the provider and user ID for checking if the user is already linked/registered
	account, user, err := h.Repo.GetAccountWithUserByProviderUserID(c.Request().Context(), tx, provider, userInfo.Subject)
	if err != nil {
		zap.L().Error("Failed to get account", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get account")
	}

	emailUserID, err := h.Repo.GetUserIDByEmail(c.Request().Context(), tx, userInfo.Email)
	if err != nil {
		zap.L().Error("Failed to check email ownership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check email ownership")
	}

	if emailUserID != nil && userID != 0 && *emailUserID != userID {
		return echo.NewHTTPError(http.StatusConflict, "Email is already linked to another user")
	}

	if emailUserID != nil && account != nil && *emailUserID != account.UserID {
		return echo.NewHTTPError(http.StatusConflict, "Email is already linked to another user")
	}

	// If the account is not found and the userID is not zero, it means the user is already registered
	// so we need to link the account to the user
	if user == nil && userID == 0 && emailUserID != nil {
		userID = *emailUserID
	}

	if user == nil && userID != 0 {
		// Link the account to the user
		newAccount, err := generateUserAccountFromOAuthUserInfo(userInfo, provider, userID)
		if err != nil {
			zap.L().Error("Failed to link account", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate account")
		}

		accountID = newAccount.ID

		// Create the account
		if err = h.Repo.CreateAccount(c.Request().Context(), tx, newAccount); err != nil {
			zap.L().Error("Failed to create account", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create account")
		}

		zap.L().Info("OAuth link account successful",
			zap.String("provider", string(provider)),
			zap.Int64("user_id", userID),
			zap.String("ip", c.RealIP()))
	} else if account == nil && userID == 0 {
		// Generate the full user object
		newUser, newAccount, err := generateUserFromOAuthUserInfo(userInfo, provider)
		if err != nil {
			zap.L().Error("Failed to generate user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate user and account")
		}

		// Create the user and account
		if err = h.Repo.CreateUserAndAccount(c.Request().Context(), tx, newUser, newAccount); err != nil {
			zap.L().Error("Failed to create user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user and account")
		}

		accountID = newAccount.ID
		userID = newUser.ID

		zap.L().Info("OAuth new user registered",
			zap.String("provider", string(provider)),
			zap.Int64("user_id", userID),
			zap.String("ip", c.RealIP()))
	} else {
		accountID = account.ID
		userID = user.ID

		zap.L().Info("OAuth login successful",
			zap.String("provider", string(provider)),
			zap.Int64("user_id", userID),
			zap.String("ip", c.RealIP()))
	}

	// If the user ID is zero, it means something went wrong (it should not happen)
	if userID == 0 {
		zap.L().Error("User ID is zero", zap.Any("user", user))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user and account")
	}

	// Create the oauth token
	err = h.Repo.CreateOAuthToken(c.Request().Context(), tx, models.OAuthToken{
		AccountID:    accountID,
		AccessToken:  token.AccessToken,
		RefreshToken: &token.RefreshToken,
		Expiry:       token.Expiry,
		TokenType:    token.TokenType,
		Provider:     provider,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		zap.L().Error("Failed to create oauth token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create oauth token")
	}

	// Generate the refresh token
	refreshToken, err := generateTokenAndSaveRefreshToken(c, h.Repo, tx, userID)
	if err != nil {
		zap.L().Error("Failed to create refresh token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create refresh token")
	}

	accessTokenCookie, err := generateAccessTokenCookieForUser(userID)
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

	// Redirect to the redirect URI
	return c.Redirect(http.StatusTemporaryRedirect, payload.RedirectURI)
}
