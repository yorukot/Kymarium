package user

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type accountResponse struct {
	ID        int64     `json:"id,string" example:"175928847299117063"`
	Provider  string    `json:"provider" example:"email"`
	Email     string    `json:"email" example:"user@example.com"`
	IsPrimary bool      `json:"is_primary" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-02T15:30:00Z"`
}

func newAccountResponse(accounts []accountResponse) []accountResponse {
	return accounts
}

// GetAccount godoc
// @Summary Get current user account info
// @Description Retrieves the authenticated user's account details
// @Tags users
// @Produce json
// @Success 200 {object} response.SuccessResponse "Account retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me/account [get]
func (h *Handler) GetAccount(c echo.Context) error {
	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	accounts, err := h.Repo.ListAccountsByUserID(c.Request().Context(), tx, *userID)
	if err != nil {
		zap.L().Error("Failed to list accounts", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list accounts")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	responseAccounts := make([]accountResponse, 0, len(accounts))
	for _, account := range accounts {
		responseAccounts = append(responseAccounts, accountResponse{
			ID:        account.ID,
			Provider:  string(account.Provider),
			Email:     account.Email,
			IsPrimary: account.IsPrimary,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, response.Success("Account retrieved successfully", newAccountResponse(responseAccounts)))
}
