package controllerstests

import (
	"testing"

	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
)

func TestMain(m *testing.M) {
	integration_tests.InitializeApiServer(m)
}
