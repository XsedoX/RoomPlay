package page_meta_dto

type PageMetaDto struct {
	NextPageToken     *string `json:"nextPageToken" swaggertype:"string" extensions:"x-nullable"`
	PreviousPageToken *string `json:"previousPageToken" swaggertype:"string" extensions:"x-nullable"`
	PageSize          uint8   `json:"pageSize" swaggertype:"integer"`
	HasNextPage       bool    `json:"hasNextPage" swaggertype:"boolean"`
}
