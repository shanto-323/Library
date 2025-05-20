package userservice

import (
	"context"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/shanto-323/Library/books"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *UserModel) error
	UpdateUser(ctx context.Context, user *UserModel) error
	DeleteUser(ctx context.Context, id string) error

	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	GetUserById(ctx context.Context, id string) (*UserModel, error)
	GetUsers(ctx context.Context, limit int, offset int) (*UserModelList, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepostiry(db *gorm.DB) Repository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, u *UserModel) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if err != nil {
		books.LogError(slog.LevelError, "DATABASE", err, "failed to create user in database")
		return err
	}
	books.LogInfo(slog.LevelInfo, "DATABASE", "database created")
	return nil
}

func (r *userRepository) GetUserById(ctx context.Context, id string) (*UserModel, error) {
	var user UserModel
	err := r.db.WithContext(ctx).Where("id= ?", id).First(&user).Error
	if err != nil {
		books.LogError(slog.LevelError, "DATABASE", err, "failed to get user from database")
		return nil, err
	}
	books.LogInfo(slog.LevelInfo, "DATABASE", "got user by id")
	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	var user UserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		books.LogError(slog.LevelError, "DATABASE", err, "failed to get user from database")
		return nil, err
	}
	books.LogInfo(slog.LevelInfo, "DATABASE", "got user by email")
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, u *UserModel) error {
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", u.ID).Updates(&u).Error
	if err != nil {
		books.LogError(slog.LevelError, "DATABASE", err, "failed to update user in database")
		return err
	}
	books.LogInfo(slog.LevelInfo, "DATABASE", "user info updated")
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{}).Error
	if err != nil {
		books.LogError(slog.LevelError, "DATABASE", err, "failed to delete user in database")
		return err
	}
	books.LogInfo(slog.LevelInfo, "DATABASE", "user deleted")
	return nil
}

func (r *userRepository) GetUsers(ctx context.Context, limit int, offset int) (*UserModelList, error) {
	var users []UserModel
	var totalUser int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&totalUser).Error; err != nil {
		books.LogError(slog.LevelError, "DATABASE", err, "failed to count users in database")
		return nil, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		books.LogError(slog.LevelError, "DATABASE", err, "failed to get users from database")
		return nil, err
	}
	totalPage := int(totalUser) / limit
	if int(totalUser)%limit != 0 {
		totalPage++
	}

	books.LogInfo(slog.LevelInfo, "DATABASE", "got all users")
	return &UserModelList{
		TotalPage: totalPage,
		TotalUser: int(totalUser),
		UserModel: users,
	}, nil
}
