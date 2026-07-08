package song_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query_response"
	"github.com/XsedoX/RoomPlay/presentation/controllers/song_controller"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests_initializer.InitializeApiServer(m)
}

func TestSearchSongsSuccess(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	query := search_song_query.SearchSongQuery{
		Query:         "song",
		NextPageToken: nil,
		PageSize:      uint8(10),
	}
	body, jsonErr := json.Marshal(query)
	require.NoError(t, jsonErr)

	req := httptest.NewRequest(
		http.MethodGet,
		constants.ApiBasePath+song_controller.SongBasePath+song_controller.SearchSongPath,
		bytes.NewReader(body),
	)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[[]search_song_query_response.SearchSongQueryResponse]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)
	response := responseWrapper.Data

	require.Equal(t, response[0].VideoId, seeder.ExternalSongData[0].VideoId)
	require.Equal(t, response[1].VideoId, seeder.ExternalSongData[1].VideoId)
	require.Equal(t, response[2].VideoId, seeder.ExternalSongData[2].VideoId)

	require.Equal(t, response[0].Title, seeder.ExternalSongData[0].Title)
	require.Equal(t, response[1].Title, seeder.ExternalSongData[1].Title)
	require.Equal(t, response[2].Title, seeder.ExternalSongData[2].Title)

	require.Equal(t, response[0].Author, seeder.ExternalSongData[0].Author)
	require.Equal(t, response[1].Author, seeder.ExternalSongData[1].Author)
	require.Equal(t, response[2].Author, seeder.ExternalSongData[2].Author)

	require.Equal(t, response[0].AlbumCoverUrl, seeder.ExternalSongData[0].AlbumCoverUrl)
	require.Equal(t, response[1].AlbumCoverUrl, seeder.ExternalSongData[1].AlbumCoverUrl)
	require.Equal(t, response[2].AlbumCoverUrl, seeder.ExternalSongData[2].AlbumCoverUrl)

	require.Equal(t, response[0].NextPageToken, seeder.ExternalSongData[0].NextPageToken)
	require.Equal(t, response[1].NextPageToken, seeder.ExternalSongData[1].NextPageToken)
	require.Equal(t, response[2].NextPageToken, seeder.ExternalSongData[2].NextPageToken)
}
