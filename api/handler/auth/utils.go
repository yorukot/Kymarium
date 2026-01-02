package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/repository"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/encrypt"
	"github.com/yorukot/kymarium/utils/id"
	"golang.org/x/oauth2"
)

// +----------------------------------------------+
// | General auth part                            |
// +----------------------------------------------+

// GenerateUser generate a user and account for the register request
func GenerateUser(registerRequest registerRequest) (models.User, models.Account, error) {
	userID, err := id.GetID()
	if err != nil {
		return models.User{}, models.Account{}, fmt.Errorf("failed to generate user ID: %w", err)
	}

	accountID, err := id.GetID()
	if err != nil {
		return models.User{}, models.Account{}, fmt.Errorf("failed to generate account ID: %w", err)
	}

	// hash the password
	passwordHash, err := encrypt.CreateArgon2idHash(registerRequest.Password)
	if err != nil {
		return models.User{}, models.Account{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// create the user
	user := models.User{
		ID:           userID,
		PasswordHash: &passwordHash,
		DisplayName:  registerRequest.DisplayName,
		Verified:     false,
		VerifyCode:   nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create the account
	account := models.Account{
		ID:             accountID,
		Provider:       models.ProviderEmail,
		ProviderUserID: strconv.FormatInt(userID, 10),
		UserID:         userID,
		Email:          registerRequest.Email,
		IsPrimary:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return user, account, nil
}

// generateRefreshToken generates a refresh token for the user
func generateRefreshToken(userID int64, userAgent string, ip string) (models.RefreshToken, error) {
	refreshTokenID, err := id.GetID()
	if err != nil {
		return models.RefreshToken{}, fmt.Errorf("failed to generate refresh token ID: %w", err)
	}

	// Extract IP address, handling cases where port may or may not be present
	ipStr := ip
	if host, _, err := net.SplitHostPort(ip); err == nil {
		// Successfully split, use the host part
		ipStr = host
	}
	parsedIP := net.ParseIP(ipStr)

	refreshToken, err := encrypt.GenerateSecureRefreshToken()
	if err != nil {
		return models.RefreshToken{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return models.RefreshToken{
		ID:        refreshTokenID,
		UserID:    userID,
		Token:     refreshToken,
		UserAgent: &userAgent,
		IP:        parsedIP,
		UsedAt:    nil,
		CreatedAt: time.Now(),
	}, nil
}

// generateRefreshTokenCookie generates a refresh token cookie
func generateRefreshTokenCookie(refreshToken models.RefreshToken) http.Cookie {
	return http.Cookie{
		Name:     models.CookieNameRefreshToken,
		Path:     "/api/auth/refresh",
		Domain:   config.Env().FrontendDomain,
		Value:    refreshToken.Token,
		HttpOnly: true,
		Secure:   config.Env().AppEnv == config.AppEnvProd,
		Expires:  refreshToken.CreatedAt.Add(time.Duration(config.Env().RefreshTokenExpiresAt) * time.Second),
		SameSite: http.SameSiteLaxMode,
	}
}

func generateAccessToken(userID int64) (string, error) {
	accessTokenClaims := encrypt.JWTSecret{
		Secret: config.Env().JWTSecretKey,
	}

	return accessTokenClaims.GenerateAccessToken(
		config.Env().AppName,
		strconv.FormatInt(userID, 10),
		time.Now().Add(time.Duration(config.Env().AccessTokenExpiresAt)*time.Second),
	)
}

func generateAccessTokenCookieForUser(userID int64) (http.Cookie, error) {
	accessToken, err := generateAccessToken(userID)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	return generateAccessTokenCookie(accessToken), nil
}

// generateAccessTokenCookie generates an access token cookie
func generateAccessTokenCookie(accessToken string) http.Cookie {
	return http.Cookie{
		Name:     models.CookieNameAccessToken,
		Path:     "/api",
		Domain:   config.Env().FrontendDomain,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   config.Env().AppEnv == config.AppEnvProd,
		Expires:  time.Now().Add(time.Duration(config.Env().AccessTokenExpiresAt) * time.Second),
		SameSite: http.SameSiteLaxMode,
	}
}

func generateEmailVerificationToken(userID int64, email string) (string, error) {
	secret := encrypt.JWTSecret{
		Secret: config.Env().JWTSecretKey,
	}

	expiresAt := time.Now().Add(time.Duration(config.Env().EmailVerifyExpiresAt) * time.Second)
	return secret.GenerateEmailVerificationToken(userID, email, expiresAt)
}

func validateEmailVerificationToken(token string) (bool, encrypt.EmailVerificationClaims, error) {
	secret := encrypt.JWTSecret{
		Secret: config.Env().JWTSecretKey,
	}

	return secret.ValidateEmailVerificationToken(token)
}

func buildEmailVerificationURL(token string) (string, error) {
	base := strings.TrimSpace(config.Env().BackendURL)
	if base == "" {
		return "", fmt.Errorf("missing backend url")
	}

	separator := "?"
	if strings.Contains(base, "?") {
		separator = "&"
	}

	return strings.TrimRight(base, "/") + "/api/auth/verify" + separator + "token=" + url.QueryEscape(token), nil
}

func sendVerificationEmail(toEmail string, token string) error {
	verifyURL, err := buildEmailVerificationURL(token)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Verify your %s account", config.Env().AppName)

	body := fmt.Sprintf("Please verify your email by clicking the link below:\n\n%s\n\nIf you did not request this, you can ignore this email.", verifyURL)
	return config.SendEmail(toEmail, nil, nil, subject, body)
}

// generateTokenAndSaveRefreshToken generates a refresh token and saves it to the database
func generateTokenAndSaveRefreshToken(e echo.Context, repo repository.Repository, tx pgx.Tx, userID int64) (models.RefreshToken, error) {
	userAgent := e.Request().UserAgent()
	ip := e.RealIP()

	refreshToken, err := generateRefreshToken(userID, userAgent, ip)
	if err != nil {
		return models.RefreshToken{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := repo.CreateRefreshToken(e.Request().Context(), tx, refreshToken); err != nil {
		return models.RefreshToken{}, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return refreshToken, nil
}

// +----------------------------------------------+
// | OAuth part                                   |
// +----------------------------------------------+

// parseProvider parse the provider from the request
func parseProvider(provider string) (models.Provider, error) {
	switch provider {
	case string(models.ProviderGoogle):
		return models.ProviderGoogle, nil
	default:
		return "", fmt.Errorf("invalid provider: %s", provider)
	}
}

// oauthGenerateStateWithPayload generate the oauth state with the payload
func oauthGenerateStateWithPayload(redirectURI string, expiresAt time.Time, userID string) (string, string, error) {
	OAuthState, err := encrypt.GenerateRandomString(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random string: %w", err)
	}

	secret := encrypt.JWTSecret{
		Secret: config.Env().JWTSecretKey,
	}

	tokenString, err := secret.GenerateOAuthState(OAuthState, redirectURI, expiresAt, userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate oauth state: %w", err)
	}

	return tokenString, OAuthState, nil
}

// oauthValidateStateWithPayload validate the oauth state with the payload
func oauthValidateStateWithPayload(oauthState string) (bool, encrypt.OAuthStateClaims, error) {
	secret := encrypt.JWTSecret{
		Secret: config.Env().JWTSecretKey,
	}

	valid, payload, err := secret.ValidateOAuthStateAndGetClaims(oauthState)
	if err != nil {
		return false, encrypt.OAuthStateClaims{}, fmt.Errorf("failed to validate oauth state: %w", err)
	}

	if payload.ExpiresAt < time.Now().Unix() {
		return false, encrypt.OAuthStateClaims{}, fmt.Errorf("oauth state expired")
	}

	return valid, payload, nil
}

// oauthVerifyTokenAndGetUserInfo verifies the token for the OAuth flow
func oauthVerifyTokenAndGetUserInfo(ctx context.Context, rawIDToken string, token *oauth2.Token, oidcProvider *oidc.Provider, oauthConfig *oauth2.Config) (*oidc.UserInfo, error) {

	// Create verifier with client ID for audience validation
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: oauthConfig.ClientID})

	// Verify the ID token
	verifiedToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	// Extract claims from verified token
	var tokenClaims map[string]any
	if err := verifiedToken.Claims(&tokenClaims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	userInfo, err := oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, nil
}

// generateUserFromOAuthUserInfo generate the user and account from the oauth user info
func generateUserFromOAuthUserInfo(userInfo *oidc.UserInfo, provider models.Provider) (models.User, models.Account, error) {
	userID, err := id.GetID()
	if err != nil {
		return models.User{}, models.Account{}, fmt.Errorf("failed to generate user ID: %w", err)
	}

	// Get the picture from the user info
	var picture *string
	var displayName string
	var claims struct {
		Picture    string `json:"picture"`
		FamilyName string `json:"family_name"`
		GivenName  string `json:"given_name"`
	}
	if err := userInfo.Claims(&claims); err == nil && claims.Picture != "" {
		picture = &claims.Picture
	}

	displayName = fmt.Sprintf("%s %s", claims.GivenName, claims.FamilyName)

	if displayName == "" {
		displayName = encrypt.GenerateRandomUserDisplayName()
	}

	// create the user
	user := models.User{
		ID:           userID,
		PasswordHash: nil,
		DisplayName:  displayName,
		Avatar:       picture,
		Verified:     true,
		VerifyCode:   nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	accountID, err := id.GetID()
	if err != nil {
		return models.User{}, models.Account{}, fmt.Errorf("failed to generate account ID: %w", err)
	}

	// create the account
	account := models.Account{
		ID:             accountID,
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: userInfo.Subject,
		Email:          userInfo.Email,
		IsPrimary:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return user, account, nil
}

// generateUserAccountFromOAuthUserInfo generate the user and account from the oauth user info
func generateUserAccountFromOAuthUserInfo(userInfo *oidc.UserInfo, provider models.Provider, userID int64) (models.Account, error) {
	accountID, err := id.GetID()
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to generate account ID: %w", err)
	}

	// create the account
	account := models.Account{
		ID:             accountID,
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: userInfo.Subject,
		Email:          userInfo.Email,
		IsPrimary:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return account, nil
}

// generateSessionCookie generates a session cookie
func generateOAuthSessionCookie(session string) http.Cookie {
	oauthSessionCookie := http.Cookie{
		Name:     models.CookieNameOAuthSession,
		Value:    session,
		Domain:   config.Env().FrontendDomain,
		HttpOnly: true,
		Path:     "/api/auth/oauth",
		Secure:   config.Env().AppEnv == config.AppEnvProd,
		Expires:  time.Now().Add(time.Duration(config.Env().OAuthStateExpiresAt) * time.Second),
		SameSite: http.SameSiteLaxMode,
	}

	return oauthSessionCookie
}
