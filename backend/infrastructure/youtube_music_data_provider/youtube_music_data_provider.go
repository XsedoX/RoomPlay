package youtube_music_data_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
)

type YoutubeMusicDataProvider struct{}

func NewYoutubeMusicDataProvider() *YoutubeMusicDataProvider {
	return &YoutubeMusicDataProvider{}
}

func (musicDataProvider *YoutubeMusicDataProvider) GetSongById(ctx context.Context, accessToken, songId string) (*music_data_response_dto.SongDataByIdResponseDto, error) {
	youtubeUrl, _ := url.ParseRequestURI("https://www.googleapis.com/youtube/v3/videos")
	params := url.Values{}
	params.Add("part", "snippet,contentDetails")
	params.Add("maxResults", "1")
	params.Add("videoCategoryId", "10") // Music category
	params.Add("id", songId)
	youtubeUrl.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		youtubeUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	type youTubeContentDetails struct {
		Duration string `json:"duration"`
	}
	type youtubeGetByIdResponseElement struct {
		Snippet struct {
			Title        string                          `json:"title"`
			Thumbnails   map[string]youtubeThumbnailData `json:"thumbnails"`
			ChannelTitle string                          `json:"channelTitle"`
		} `json:"snippet"`
		ContentDetails youTubeContentDetails `json:"contentDetails"`
		Id             string                `json:"id"`
	}
	type youtubeGetByIdResponse struct {
		Items []youtubeGetByIdResponseElement `json:"items"`
	}
	var response youtubeGetByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Items) < 0 {
		return nil, application_error.NewApplicationError(
			"YoutubeMusicDataProvider.GetSongById.NoSongByThisId",
			"No song found with the provided ID",
			nil,
			application_error_type.Unexpected,
		)
	}
	songDurationSeconds, err := parseISODuration(response.Items[0].ContentDetails.Duration)
	if err != nil {
		return nil, application_error.NewApplicationError(
			"YoutubeMusicDataProvider.GetSongById.InvalidDuration",
			"Invalid duration format for the song",
			err,
			application_error_type.Unexpected,
		)
	}
	result := &music_data_response_dto.SongDataByIdResponseDto{
		VideoId:       response.Items[0].Id,
		Title:         response.Items[0].Snippet.Title,
		Author:        response.Items[0].Snippet.ChannelTitle,
		AlbumCoverUrl: response.Items[0].Snippet.Thumbnails["default"].Url,
		LengthSeconds: uint16(songDurationSeconds),
		MusicProvider: music_provider.YouTube,
		Isrc:          nil, // YouTube API does not provide ISRC directly
	}
	return result, nil
}

func (musicDataProvider *YoutubeMusicDataProvider) SearchSongsByQuery(ctx context.Context, accessToken, query string, nextPageToken *string, pageSize uint8) (*music_data_response_dto.MusicDataResponseDto, error) {
	youtubeUrl, _ := url.ParseRequestURI("https://www.googleapis.com/youtube/v3/search")
	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("maxResults", strconv.Itoa(int(pageSize)))
	params.Add("q", query)
	params.Add("videoCategoryId", "10") // Music category
	params.Add("type", "video")
	if nextPageToken != nil {
		nextPageTokenString := *nextPageToken
		params.Add("pageToken", nextPageTokenString)
	}
	youtubeUrl.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		youtubeUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	type youtubeSearchResponseElement struct {
		Snippet struct {
			Title        string                          `json:"title"`
			Thumbnails   map[string]youtubeThumbnailData `json:"thumbnails"`
			ChannelTitle string                          `json:"channelTitle"`
		} `json:"snippet"`
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
	}
	type youtubeSearchResponse struct {
		Items         []youtubeSearchResponseElement `json:"items"`
		NextPageToken string                         `json:"nextPageToken"`
		PrevPageToken string                         `json:"prevPageToken"`
	}
	var response youtubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	result := &music_data_response_dto.MusicDataResponseDto{}
	result.Songs = make([]music_data_response_dto.SearchSongDataResponseDto, 0, len(response.Items))

	for _, item := range response.Items {
		if err != nil {
			log.Printf("Error parsing duration for video %s: %v", item.Id.VideoId, err)
			continue
		}
		result.Songs = append(result.Songs, music_data_response_dto.SearchSongDataResponseDto{
			VideoId:       item.Id.VideoId,
			Title:         item.Snippet.Title,
			Author:        item.Snippet.ChannelTitle,
			AlbumCoverUrl: item.Snippet.Thumbnails["default"].Url,
			MusicProvider: music_provider.YouTube,
			Isrc:          nil, // YouTube API does not provide ISRC directly
		})
	}
	result.PageMetaDto = page_meta_dto.PageMetaDto{
		NextPageToken:     &response.NextPageToken,
		PreviousPageToken: &response.PrevPageToken,
		HasNextPage:       response.NextPageToken != "",
		PageSize:          uint8(len(response.Items)),
	}
	return result, nil
}

func parseISODuration(d string) (int, error) {
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	m := re.FindStringSubmatch(d)
	if m == nil {
		return 0, fmt.Errorf("invalid duration: %s", d)
	}
	var secs int
	for i, unit := range []int{3600, 60, 1} {
		if m[i+1] != "" {
			v, _ := strconv.Atoi(m[i+1])
			secs += v * unit
		}
	}
	return secs, nil
}

type youtubeThumbnailData struct {
	Url string `json:"url"`
}
