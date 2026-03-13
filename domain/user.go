package domain

// User adalah entitas inti dalam domain bisnis.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UserRepository mendefinisikan kontrak akses data untuk entitas User.
type UserRepository interface {
	Save(user User) error
	FindByID(id string) (*User, error)
	Delete(id string) error
}