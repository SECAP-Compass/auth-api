package domain

import "context"

type IUserQueryRepository interface {
	FindByEmail(context.Context, string) (*User, error)
	FindByID(context.Context, uint) (*User, error)
}

type IUserCommandRepository interface {
	Store(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, string) error
}

type IJtiRecordQueryRepository interface {
	FindByID(context.Context, string) (*JtiRecord, error)
	FindByUserID(context.Context, uint) (*JtiRecord, error)
}

type IJtiRecordCommandRepository interface {
	Store(context.Context, *JtiRecord) error
	Delete(context.Context, string) error
}

type ICityQueryRepository interface {
	FindByID(context.Context, uint) (*City, error)
}
