package monitor

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/models/monitorm"
	"github.com/yorukot/kymarium/utils/response"
)

func validateStruct(c echo.Context, v any, message string) error {
	if err := validator.New().Struct(v); err != nil {
		if errResponse, ok := response.ValidationErrorResponse(err, message, "VALIDATION_ERROR"); ok {
			return c.JSON(http.StatusBadRequest, errResponse)
		}

		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: message,
		})
	}

	return nil
}

func validateMonitorConfig(c echo.Context, monitorType models.MonitorType, configRaw json.RawMessage) error {
	switch monitorType {
	case models.MonitorTypeHTTP:
		var config monitorm.HTTPMonitorConfig
		if err := json.Unmarshal(configRaw, &config); err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Message: "Invalid monitor config",
			})
		}
		return validateStruct(c, config, "Invalid monitor config")
	case models.MonitorTypePing:
		var config monitorm.PingMonitorConfig
		if err := json.Unmarshal(configRaw, &config); err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Message: "Invalid monitor config",
			})
		}
		return validateStruct(c, config, "Invalid monitor config")
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Unsupported monitor type")
	}
}
