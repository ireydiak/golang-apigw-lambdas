package repository

type User struct {
	Name string `json:"name"`
}

type UserRepository struct {
	// db connection, etc.
}

func NewUserRepo() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) List() []User {
	return []User{
		{Name: "Jane Doe"},
		{Name: "John Doe"},
	}
}
