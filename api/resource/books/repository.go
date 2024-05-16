package books

type IBooksRepo interface {
	GetBooks() ([]bookEntity, error)
}

type BooksRepo struct{}

func (repo *BooksRepo) GetBooks() ([]bookEntity, error) {
	return []bookEntity{
		{
			Title:         "The Fellowship of the Ring",
			Author:        "JRR Tolkien",
			Price:         20,
			NumberOfPages: 432,
			Genre:         "",
			ReleaseYear:   1954,
		},
	}, nil
}
