package song_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query_song_data_dto"
	"github.com/XsedoX/RoomPlay/presentation/controllers/song_controller"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/presentation/response"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
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
		NextPageToken: new("whatever"),
		PageSize:      10,
	}

	params := url.Values{}
	params.Add("query", query.Query)
	params.Add("pageSize", strconv.Itoa(query.PageSize))
	params.Add("nextPageToken", *query.NextPageToken)

	url, parsingErr := url.ParseRequestURI(constants.ApiBasePath + song_controller.SongBasePath + song_controller.SearchSongPath)
	require.NoError(t, parsingErr)
	url.RawQuery = params.Encode()

	req := httptest.NewRequest(
		http.MethodGet,
		url.String(),
		nil,
	)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[[]search_song_query_song_data_dto.SearchSongQuerySongDataDto]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)
	response := responseWrapper.Data

	require.Equal(t, response[0].VideoId, seeder.ExternalSongData.Songs[0].VideoId)
	require.Equal(t, response[1].VideoId, seeder.ExternalSongData.Songs[1].VideoId)
	require.Equal(t, response[2].VideoId, seeder.ExternalSongData.Songs[2].VideoId)

	require.Equal(t, response[0].Title, seeder.ExternalSongData.Songs[0].Title)
	require.Equal(t, response[1].Title, seeder.ExternalSongData.Songs[1].Title)
	require.Equal(t, response[2].Title, seeder.ExternalSongData.Songs[2].Title)

	require.Equal(t, response[0].Author, seeder.ExternalSongData.Songs[0].Author)
	require.Equal(t, response[1].Author, seeder.ExternalSongData.Songs[1].Author)
	require.Equal(t, response[2].Author, seeder.ExternalSongData.Songs[2].Author)

	require.Equal(t, response[0].AlbumCoverUrl, seeder.ExternalSongData.Songs[0].AlbumCoverUrl)
	require.Equal(t, response[1].AlbumCoverUrl, seeder.ExternalSongData.Songs[1].AlbumCoverUrl)
	require.Equal(t, response[2].AlbumCoverUrl, seeder.ExternalSongData.Songs[2].AlbumCoverUrl)

	require.Equal(t, responseWrapper.Meta.NextPageToken, seeder.ExternalSongData.PageMetaDto.NextPageToken)
	require.Equal(t, responseWrapper.Meta.PageSize, seeder.ExternalSongData.PageMetaDto.PageSize)
	require.Equal(t, responseWrapper.Meta.PreviousPageToken, seeder.ExternalSongData.PageMetaDto.PreviousPageToken)
	require.Equal(t, responseWrapper.Meta.HasNextPage, seeder.ExternalSongData.PageMetaDto.HasNextPage)
}

func TestSearchSongsTooShortQuery(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	query := search_song_query.SearchSongQuery{
		Query:         "t",
		NextPageToken: new("whatever"),
		PageSize:      10,
	}
	params := url.Values{}
	params.Add("query", query.Query)
	params.Add("pageSize", strconv.Itoa(query.PageSize))
	params.Add("nextPageToken", *query.NextPageToken)

	url, parsingErr := url.ParseRequestURI(constants.ApiBasePath + song_controller.SongBasePath + song_controller.SearchSongPath)
	require.NoError(t, parsingErr)
	url.RawQuery = params.Encode()

	req := httptest.NewRequest(
		http.MethodGet,
		url.String(),
		nil,
	)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var responseWrapper response.ProblemDetails
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	require.Equal(t, responseWrapper.Description, "One or more fields are not correctly filled")
	require.Equal(t, responseWrapper.Title, "Validation error occurred.")
	require.Equal(t, responseWrapper.Status, http.StatusUnprocessableEntity)
	require.Equal(t, responseWrapper.Instance, url.String())
	require.Equal(t, "SearchSongsByQuery.Validation", responseWrapper.Type)
	validationErrors := responseWrapper.ValidationErrors
	require.Len(t, validationErrors, 1)
	require.Equal(t, "Song Query must be greater than or equal to 2.", validationErrors["query"])
}

func TestSearchSongsTooLongQuery(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	query := search_song_query.SearchSongQuery{
		Query:         gofakeit.LetterN(51),
		NextPageToken: new("whatever"),
		PageSize:      10,
	}

	params := url.Values{}
	params.Add("query", query.Query)
	params.Add("pageSize", strconv.Itoa(query.PageSize))
	params.Add("nextPageToken", *query.NextPageToken)

	url, parsingErr := url.ParseRequestURI(constants.ApiBasePath + song_controller.SongBasePath + song_controller.SearchSongPath)
	require.NoError(t, parsingErr)
	url.RawQuery = params.Encode()

	req := httptest.NewRequest(
		http.MethodGet,
		url.String(),
		nil,
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var responseWrapper response.ProblemDetails
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	require.Equal(t, responseWrapper.Description, "One or more fields are not correctly filled")
	require.Equal(t, responseWrapper.Title, "Validation error occurred.")
	require.Equal(t, responseWrapper.Status, http.StatusUnprocessableEntity)
	require.Equal(t, responseWrapper.Instance, url.String())
	require.Equal(t, "SearchSongsByQuery.Validation", responseWrapper.Type)
	validationErrors := responseWrapper.ValidationErrors
	require.Len(t, validationErrors, 1)
	require.Equal(t, "Song Query must be less than or equal to 50.", validationErrors["query"])
}
