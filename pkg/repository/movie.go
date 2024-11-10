package repository

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type MovieRepo struct {
	// insert struct dependencies (db connections, etc.)
}

func NewMovieRepo() *MovieRepo {
	return &MovieRepo{}
}

func (repo *MovieRepo) List() ([]Movie, error) {
	return []Movie{
		{ID: "terrifier-2-2023", Title: "Terrifier 2"},
		{ID: "terrifier-3-2024", Title: "Terrifier 3"},
	}, nil
}
