package music_data_response_dto

type MusicDataResponseDto struct {
	Url            string
	Title          string
	Author         string
	LengthSeconds  uint16
	AlabumCoverUrl string
	NextPageToken  string
}
