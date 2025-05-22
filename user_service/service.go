package userservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/ksuid"
	"github.com/shanto-323/Library/books"
)

type Service interface {
	SignUpNewUser(ctx context.Context, name string, password string, email string, phone string, u_type string) (*UserModel, error)
	LoginWithEmail(ctx context.Context, email string, password string) (*UserModel, error)
	Logout(ctx context.Context, id string) error

	UpdateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	DeleteUserById(ctx context.Context, id string) error
	NewToken(ctx context.Context, id string, r_token string) (string, error)

	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	GetUserById(ctx context.Context, id string) (*UserModel, error)
	GetAllUser(ctx context.Context, limit int64, offset int64) (*UserModelList, error)
}

type userService struct {
	userRepository Repository
}

func NewUserService(repository Repository) Service {
	return &userService{
		userRepository: repository,
	}
}

func (s userService) SignUpNewUser(ctx context.Context, name string, password string, email string, phone string, u_type string) (*UserModel, error) {
	user := &UserModel{
		ID:       ksuid.New().String(),
		Name:     name,
		Email:    email,
		Phone:    phone,
		UserType: u_type,
	}

	token, r_token, err := CreateTokens(user.Email, user.ID, user.UserType)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "error creating tokens")
		return nil, err
	}
	user.Token = token
	user.RefreshToken = r_token
	user.Password, _ = CreateNewHashPassword(password)

	err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "error createing user for signup")
		return nil, err
	}
	books.LogInfo(slog.LevelInfo, "SERVICE", "signed up succesful")
	return user, nil
}

func (s userService) LoginWithEmail(ctx context.Context, email string, password string) (*UserModel, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "user not found")
		return nil, err
	}

	err = CompareWithHash(password, user.Password)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "password not matched")
		return nil, err
	}

	token, r_token, err := CreateTokens(user.Email, user.ID, user.UserType)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "error creating tokens")
		return nil, err
	}
	user.Token = token
	user.RefreshToken = r_token

	return user, nil
}

func (s userService) Logout(ctx context.Context, id string) error {
	user, _ := s.userRepository.GetUserById(ctx, id)
	user.Token = ""
	user.RefreshToken = ""
	err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "error logging out")
		return err
	}
	books.LogInfo(slog.LevelInfo, "SERVICE", "log out success")
	return nil
}

func (s userService) GetUserById(ctx context.Context, uid string) (*UserModel, error) {
	return s.userRepository.GetUserById(ctx, uid)
}

func (s userService) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	return s.userRepository.GetUserByEmail(ctx, email)
}

func (s userService) UpdateUser(ctx context.Context, user *UserModel) (*UserModel, error) {
	dbUser, err := s.GetUserById(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if dbUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	dbUser, err = mutationHelper(dbUser, user)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.UpdateUser(ctx, dbUser)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "error updating user")
		return nil, err
	}
	books.LogInfo(slog.LevelInfo, "SERVICE", "user updated succesfully")
	return user, nil
}

func (s userService) DeleteUserById(ctx context.Context, uid string) error {
	return s.userRepository.DeleteUser(ctx, uid)
}

func (s userService) NewToken(ctx context.Context, id string, r_token string) (string, error) {
	books.LogInfo(slog.LevelInfo, "client", id)
	user, err := s.userRepository.GetUserById(ctx, id)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("no user found")
	}

	if user.RefreshToken == "" || user.RefreshToken != r_token {
		books.LogError(slog.LevelError, "SERVICE", fmt.Errorf("token nil %s or token not matched", r_token), fmt.Sprintf("token nil %s or token not matched", r_token))
		return "", fmt.Errorf("token nil %s or token not matched", r_token)
	}

	token, _, err := CreateTokens(user.Email, user.ID, user.UserType)
	if err != nil {
		books.LogError(slog.LevelError, "SERVICE", err, "error creating token")
		return "", err
	}
	books.LogInfo(slog.LevelInfo, "SERVICE", "got new access token")
	return token, err
}

func (s userService) GetAllUser(ctx context.Context, limit int64, offset int64) (*UserModelList, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.userRepository.GetUsers(ctx, int(limit), int(offset))
}

func mutationHelper(dbUser *UserModel, user *UserModel) (*UserModel, error) {
	hasChange := false
	if dbUser.Name != user.Name {
		dbUser.Name = user.Name
		hasChange = true
	}
	if dbUser.Password != user.Password {
		dbUser.Password = user.Password
		hasChange = true
	}
	if dbUser.Phone != user.Phone {
		dbUser.Phone = user.Phone
		hasChange = true
	}

	if !hasChange {
		return nil, fmt.Errorf("value not changed")
	}
	return dbUser, nil
}
