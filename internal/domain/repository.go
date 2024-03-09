package domain

import "context"

type IUserRepository interface {
	// Query
	FindByEmail(context.Context, string) (*User, error)
	FindByID(context.Context, string) (*User, error)

	// Command
	Store(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, string) error
}

type IJtiRecordRepository interface {
	// Query
	FindByID(context.Context, string) (*JtiRecord, error)
	FindByUserID(context.Context, string) (*JtiRecord, error)

	// Command
	Store(context.Context, *JtiRecord) error
	Delete(context.Context, string) error
}
