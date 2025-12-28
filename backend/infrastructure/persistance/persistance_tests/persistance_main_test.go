package persistance_tests

import (
	"testing"

	"xsedox.com/main/test_helpers/integration_tests"
)

func TestMain(m *testing.M) {
	integration_tests.InitializeDatabaseContainer()
	integration_tests.RunTestsWithDatabase(m)
}
