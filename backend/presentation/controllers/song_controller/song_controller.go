package song_controller

import (
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_query_handler"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query_dto"
	"github.com/XsedoX/RoomPlay/presentation/response"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/gorilla/schema"
)

const (
	SongBasePath   = "/song"
	SearchSongPath = "/search"
)

type SongController struct {
	searchSongQueryHandler i_query_handler.IQueryHandlerWithRequest[*search_song_query.SearchSongQuery, *search_song_query_dto.SearchSongQueryDto]
}

func NewSongController(
	searchSongQueryHandler i_query_handler.IQueryHandlerWithRequest[*search_song_query.SearchSongQuery, *search_song_query_dto.SearchSongQueryDto],
) *SongController {
	return &SongController{
		searchSongQueryHandler: searchSongQueryHandler,
	}
}

func (sc *SongController) SearchSongsByQuery(w http.ResponseWriter, r *http.Request) {
	var query search_song_query.SearchSongQuery
	decoder := schema.NewDecoder()
	paramsDecodeErr := decoder.Decode(&query, r.URL.Query())
	if paramsDecodeErr != nil {
		response.WriteJsonDecodingFailure(
			w,
			"SearchSongsByQuery.Decoding",
			paramsDecodeErr,
			r.URL.RequestURI(),
		)
		return
	}

	validationErr := setup_validation.ValidatorInstance.Struct(query)
	if validationErr != nil {
		response.WriteJsonValidationFailure(w,
			"SearchSongsByQuery.Validation",
			r.URL.RequestURI(),
			validationErr,
		)
		return
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
		return
	}

	response.WriteJsonSuccess(w,
		searchSongResponse.Songs,
		searchSongResponse.PageMetaDto,
	)
}
