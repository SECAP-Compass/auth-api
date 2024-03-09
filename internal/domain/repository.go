package domain

import "github.com/golang-jwt/jwt"

type IUserRepository interface {
	// Query
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)

	// Command
	Store(user *User) error
	Update(user *User) error
	Delete(id string) error
}

type IJwtRepository interface {
	// Query
	FindByID(jti string) (*jwt.Token, error)

	// Command
	Store(jwt *jwt.Token) error
}
