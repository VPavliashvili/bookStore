package books

import (
	"encoding/json"
)

type bookEntity struct {
	ID            int
	Title         string
	Author        string
	Genre         string
	NumberOfPages int
	Price         int
	ReleaseYear   int
}

func (b bookEntity) ToDto() bookDTO {
	return bookDTO(b)
}

type bookDTO struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Genre         string `json:"genre"`
	NumberOfPages int    `json:"numberOfPages"`
	Price         int    `json:"price"`
	ReleaseYear   int    `json:"releaseYear"`
}

func (b bookDTO) ToEntity() bookEntity {
	return bookEntity(b)
}

type ActionResponse struct {
	ResourceId int `json:"resourceId"`
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	json, _ := json.Marshal(e)
	return string(json[:])
}

type internalErr struct {
	message string
}

func (e internalErr) Error() string {
	return e.message
}

type notfoundErr struct {
	message string
}

func (e notfoundErr) Error() string {
	return e.message
}

type badreqErr struct {
	message string
}

func (e badreqErr) Error() string {
	return e.message
}
