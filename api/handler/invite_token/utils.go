package invitetoken

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/yorukot/kymarium/utils/config"
)

func buildFrontendInviteURL(result string, teamID int64, email string) string {
	base := strings.TrimSpace(config.Env().FrontendDomain)
	if base == "" {
		return "/"
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}

	base = strings.TrimRight(base, "/")

	if result == "success" && teamID != 0 {
		return base + "/" + strconv.FormatInt(teamID, 10)
	}

	path := "/invites/" + result

	values := url.Values{}
	if teamID != 0 {
		values.Set("team_id", strconv.FormatInt(teamID, 10))
	}
	if email != "" {
		values.Set("email", email)
	}

	if len(values) == 0 {
		return base + path
	}

	return base + path + "?" + values.Encode()
}
