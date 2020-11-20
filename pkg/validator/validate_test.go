package validator

import (
	"fmt"
	"os"
	"testing"

	"github.com/microlib/simple"
)

var (
	logger *simple.Logger
)

// TestEnvars - entry for validator tests
func TestEnvars(t *testing.T) {
	logger := &simple.Logger{Level: "trace"}

	t.Run("ValidateEnvars : should fail", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "")
		err := ValidateEnvars(logger)
		if err == nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with no error - got (%v) wanted (%v)", "ValidateEnvars", err, nil))
		}
	})

	t.Run("ValidateEnvars : should pass", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "info")
		os.Setenv("ADV_URL", "/test")
		os.Setenv("VERSION", "1.0.3")
		os.Setenv("NAME", "test")
		os.Setenv("CRON", "0 0/9 * * *")
		os.Setenv("SLEEP", "3600")
		os.Setenv("USER", "1234")
		os.Setenv("AWS_ACCOUNT", "1234")
		os.Setenv("AWS_REGION", "eu-east-1")
		os.Setenv("AWS_BUCKET", "test1234")
		os.Setenv("AWS_ACCESS_KEY_ID", "test1234")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "test1234")
		os.Setenv("AWS_USER", "test1234")
		os.Setenv("BASE_DIR", "tests")
		err := ValidateEnvars(logger)
		if err != nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with error - got (%v) wanted (%v)", "ValidateEnvars", err, nil))
		}
	})
}
