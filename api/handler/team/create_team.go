package team

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/id"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type createTeamRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// CreateTeam godoc
// @Summary Create a team
// @Description Creates a new team and assigns the requesting user as owner
// @Tags teams
// @Accept json
// @Produce json
// @Param request body createTeamRequest true "Team create request"
// @Success 200 {object} response.SuccessResponse "Team created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams [post]
func (h *Handler) CreateTeam(c echo.Context) error {
	var req createTeamRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

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

	now := time.Now()

	teamID, err := id.GetID()
	if err != nil {
		zap.L().Error("Failed to generate team ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate team ID")
	}

	team := models.Team{
		ID:        teamID,
		Name:      req.Name,
		UpdatedAt: now,
		CreatedAt: now,
	}

	memberID, err := id.GetID()
	if err != nil {
		zap.L().Error("Failed to generate team member ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate team member ID")
	}

	member := models.TeamMember{
		ID:        memberID,
		TeamID:    teamID,
		UserID:    *userID,
		Role:      models.MemberRoleOwner,
		UpdatedAt: now,
		CreatedAt: now,
	}

	if err := h.Repo.CreateTeam(c.Request().Context(), tx, team); err != nil {
		zap.L().Error("Failed to create team", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create team")
	}

	if err := h.Repo.CreateTeamMember(c.Request().Context(), tx, member); err != nil {
		zap.L().Error("Failed to create team member", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create team member")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Team created successfully", team))
}
