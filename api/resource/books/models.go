package books

type bookEntity struct {
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
	Title         string `json:"title"`
	Author        string `json:"author"`
	Genre         string `json:"genre"`
	NumberOfPages int    `json:"numberOfPages"`
	Price         int    `json:"price"`
	ReleaseYear   int    `json:"releaseYear"`
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
