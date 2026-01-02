package testutil

import (
	"os"
	"sync"
	"testing"

	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/id"
	"go.uber.org/zap"
)

var (
	initOnce sync.Once
	initErr  error
)

// InitTestEnv seeds env vars, config, ID generator, and logger for handler tests.
func InitTestEnv(t *testing.T) {
	t.Helper()

	initOnce.Do(func() {
		seedEnvVars()

		if _, err := config.InitConfig(); err != nil {
			initErr = err
			return
		}

		if err := id.Init(); err != nil {
			initErr = err
			return
		}

		zap.ReplaceGlobals(zap.NewNop())
	})

	if initErr != nil {
		t.Fatalf("failed to initialize test environment: %v", initErr)
	}
}

func seedEnvVars() {
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "test")
	_ = os.Setenv("DB_PASSWORD", "test")
	_ = os.Setenv("DB_NAME", "test")
	_ = os.Setenv("DB_SSL_MODE", "disable")
	_ = os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	_ = os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	_ = os.Setenv("GOOGLE_REDIRECT_URL", "http://localhost/callback")
	_ = os.Setenv("JWT_SECRET_KEY", "secret")
	_ = os.Setenv("APP_MACHINE_ID", "1")
}
