package song_controller

import (
	"encoding/json"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_query_handler"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query_response"
	"github.com/XsedoX/RoomPlay/presentation/response"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
)

const (
	SongBasePath   = "/song"
	SearchSongPath = "/search"
)

type SongController struct {
	searchSongQueryHandler i_query_handler.IQueryHandlerWithRequest[*search_song_query.SearchSongQuery, []search_song_query_response.SearchSongQueryResponse]
}

func NewSongController(
	searchSongQueryHandler i_query_handler.IQueryHandlerWithRequest[*search_song_query.SearchSongQuery, []search_song_query_response.SearchSongQueryResponse],
) *SongController {
	return &SongController{
		searchSongQueryHandler: searchSongQueryHandler,
	}
}

func (sc *SongController) SearchSongsByQuery(w http.ResponseWriter, r *http.Request) {
	var query search_song_query.SearchSongQuery
	bodyDecodeErr := json.NewDecoder(r.Body).Decode(&query)
	if bodyDecodeErr != nil {
		response.WriteJsonFailure(w,
			"SearchSongsByQuery.Decoding",
			"Problem with decoding request body",
			bodyDecodeErr.Error(),
			r.URL.RequestURI(),
			http.StatusBadRequest,
		)
	}

	validationErr := setup_validation.ValidatorInstance.Struct(query)
	if validationErr != nil {
		response.WriteJsonValidationFailure(w,
			"SearchSongsByQuery.Validation",
			r.URL.RequestURI(),
			validationErr,
		)
	}

	searchSongResponse, searchSongResponseErr := sc.searchSongQueryHandler.Handle(
		r.Context(),
		&query,
	)
	if searchSongResponseErr != nil {
		response.WriteJsonApplicationFailure(w,
			searchSongResponseErr,
			r.URL.RequestURI(),
		)
	}

	response.WriteJsonSuccess(w,
		http.StatusOK,
		searchSongResponse,
	)
}
