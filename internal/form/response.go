package form

import (
	"time"
)

type PaginationResponsePartial struct {
	PageIdx  int64 `json:"pageIdx"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
}

type UserLocation struct {
	Id        string    `json:"id"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserId    string    `json:"userId"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
}

type GetLocationsResponse struct {
	Pagination PaginationResponsePartial `json:"pagination"`
	Data       []UserLocation            `json:"data"`
}
