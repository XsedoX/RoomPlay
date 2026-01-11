import (
	"testing"

	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
)

func TestLogoutSuccess(t *testing.T) {
	testServer := integration_tests.TestServer
	r := testServer.Router()
	
	req := httptest.NewRequest(http.MethodGet, helpers.ApiBasePath+controllers.RoomBasePath, nil)
	w := httptest.NewRecorder()
	w.
}
