package controllers_tests

import (
	"testing"

	"xsedox.com/main/test_helpers/integration_tests"
)

func TestMain(m *testing.M) {
	integration_tests.InitializeApiServer(m)
}
