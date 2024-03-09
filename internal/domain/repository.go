package domain

type IUserRepository interface {
	// Query
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)

	// Command
	Store(user *User) error
	Update(user *User) error
	Delete(id string) error
}
