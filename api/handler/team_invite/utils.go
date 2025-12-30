package team_invite

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils/config"
)

func normalizeInviteRole(role string) (models.MemberRole, error) {
	if role == "" {
		return models.MemberRoleMember, nil
	}

	switch models.MemberRole(role) {
	case models.MemberRoleAdmin, models.MemberRoleMember, models.MemberRoleViewer:
		return models.MemberRole(role), nil
	default:
		return "", fmt.Errorf("invalid role")
	}
}

func buildInviteAcceptURL(token string) (string, error) {
	base := strings.TrimSpace(config.Env().BackendURL)
	if base == "" {
		return "", fmt.Errorf("missing backend url")
	}

	return strings.TrimRight(base, "/") + "/api/invite-tokens/" + url.PathEscape(token), nil
}

func buildInviteEmailBody(teamName string, acceptURL string, expiresAt time.Time) string {
	return fmt.Sprintf(
		"You have been invited to join %s.\n\nAccept invite: %s\n\nThis invite expires on %s.\n",
		teamName,
		acceptURL,
		expiresAt.UTC().Format(time.RFC3339),
	)
}
