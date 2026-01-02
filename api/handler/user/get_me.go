package user

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

type userResponse struct {
	ID          int64     `json:"id,string" example:"175928847299117063"`
	DisplayName string    `json:"display_name" example:"John Doe"`
	Avatar      *string   `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"`
	Verified    bool      `json:"verified" example:"false"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-02T15:30:00Z"`
}

func newUserResponse(user *models.User) userResponse {
	return userResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Avatar:      user.Avatar,
		Verified:    user.Verified,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// GetMe godoc
// @Summary Get current user
// @Description Retrieves the authenticated user's profile
// @Tags users
// @Produce json
// @Success 200 {object} response.SuccessResponse "User retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me [get]
func (h *Handler) GetMe(c echo.Context) error {
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

	user, err := h.Repo.GetUserByID(c.Request().Context(), tx, *userID)
	if err != nil {
		zap.L().Error("Failed to get user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("User retrieved successfully", newUserResponse(user)))
}
