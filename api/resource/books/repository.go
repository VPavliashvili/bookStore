package books

type IBooksRepo interface {
	GetBooks() ([]bookEntity, error)
	GetBookById(int) (bookEntity, error)
}

type BooksRepo struct{}

func (repo *BooksRepo) GetBooks() ([]bookEntity, error) {
	return []bookEntity{
		{
			Title:         "The Fellowship of the Ring",
			Author:        "JRR Tolkien",
			Price:         20,
			NumberOfPages: 432,
			Genre:         "fantasy",
			ReleaseYear:   1954,
		},
	}, nil
}

func (repo *BooksRepo) GetBookById(id int) (bookEntity, error) {
	return bookEntity{
		Title:         "The Two Towers",
		Author:        "JRR Tolkien",
		Genre:         "fantasy",
		NumberOfPages: 352,
		Price:         20,
		ReleaseYear:   1954,
	}, nil
}
