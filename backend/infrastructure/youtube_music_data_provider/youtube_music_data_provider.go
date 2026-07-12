package youtube_music_data_provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
)

type YoutubeMusicDataProvider struct{}

func NewYoutubeMusicDataProvider() *YoutubeMusicDataProvider {
	return &YoutubeMusicDataProvider{}
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

	type YouTubeSearchThumbnailData struct {
		Url    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	type YouTubeSearchResponseElement struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		Id   struct {
			Kind       string `json:"kind"`
			VideoId    string `json:"videoId"`
			ChannelId  string `json:"channelId"`
			PlaylistId string `json:"playlistId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt          time.Time                             `json:"publishedAt"`
			ChannelId            string                                `json:"channelId"`
			Title                string                                `json:"title"`
			Description          string                                `json:"description"`
			Thumbnails           map[string]YouTubeSearchThumbnailData `json:"thumbnails"`
			ChannelTitle         string                                `json:"channelTitle"`
			LiveBroadcastContent string                                `json:"liveBroadcastContent"`
		} `json:"snippet"`
	}
	type YoutubeSearchResponse struct {
		Kind          string `json:"kind"`
		Etag          string `json:"etag"`
		NextPageToken string `json:"nextPageToken"`
		PrevPageToken string `json:"prevPageToken"`
		RegionCode    string `json:"regionCode"`
		PageInfo      struct {
			TotalResults   int `json:"totalResults"`
			ResultsPerPage int `json:"resultsPerPage"`
		} `json:"pageInfo"`
		Items []YouTubeSearchResponseElement `json:"items"`
	}
	var response YoutubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	result := &music_data_response_dto.MusicDataResponseDto{}
	result.Songs = make([]music_data_response_dto.SongDataResponseDto, 0, len(response.Items))

	for _, item := range response.Items {
		result.Songs = append(result.Songs, music_data_response_dto.SongDataResponseDto{
			VideoId:       item.Id.VideoId,
			Title:         item.Snippet.Title,
			Author:        item.Snippet.ChannelTitle,
			AlbumCoverUrl: item.Snippet.Thumbnails["default"].Url,
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
