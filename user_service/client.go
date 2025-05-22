package userservice

import (
	"context"
	"log/slog"

	"github.com/shanto-323/Library/books"
	"github.com/shanto-323/Library/user_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AuthServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	s := pb.NewAuthServiceClient(conn)
	return &Client{
		conn:    conn,
		service: s,
	}, nil
}

func (c *Client) SignUp(ctx context.Context, name string, password string, email string, phone string, u_type string) (*UserModel, error) {
	resp, err := c.service.SignUp(
		ctx,
		&pb.SignUpRequest{
			User: &pb.User{
				Name:     name,
				Password: password,
				Phone:    phone,
				Email:    email,
				UserType: u_type,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	user := resp.UserModel.User
	return &UserModel{
		ID:           user.Id,
		Name:         user.Name,
		Password:     user.Password,
		Email:        user.Email,
		Phone:        user.Phone,
		UserType:     user.UserType,
		Token:        resp.UserModel.Token,
		RefreshToken: resp.UserModel.RefreshToken,
		CreatedAt:    resp.UserModel.CreatedAt.AsTime(),
		UpdatedAt:    resp.UserModel.UpdatedAt.AsTime(),
	}, nil
}

func (c *Client) SignIn(ctx context.Context, email string, password string) (*UserModel, error) {
	resp, err := c.service.SignIn(
		ctx,
		&pb.SignInRequest{
			Email:    email,
			Password: password,
		},
	)

	if err != nil {
		return nil, err
	}

	user := resp.UserModel.User
	return &UserModel{
		ID:           user.Id,
		Name:         user.Name,
		Password:     user.Password,
		Email:        user.Email,
		Phone:        user.Phone,
		UserType:     user.UserType,
		Token:        resp.UserModel.Token,
		RefreshToken: resp.UserModel.RefreshToken,
		CreatedAt:    resp.UserModel.CreatedAt.AsTime(),
		UpdatedAt:    resp.UserModel.UpdatedAt.AsTime(),
	}, nil
}

func (c *Client) Logout(ctx context.Context, id string) (*string, error) {
	resp, err := c.service.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &resp.Msg, nil

}

func (c *Client) UpdateUser(ctx context.Context, user *UserModel) (*UserModel, error) {
	resp, err := c.service.UpdateUser(ctx, &pb.UpdateUserRequest{
		User: &pb.User{
			Id:       user.ID,
			Name:     user.Name,
			Password: user.Password,
			Phone:    user.Phone,
			Email:    user.Email,
			UserType: user.UserType,
		},
	})
	if err != nil {
		return nil, err
	}

	u := resp.UserModel.User

	return &UserModel{
		ID:           u.Id,
		Name:         u.Name,
		Password:     u.Password,
		Email:        u.Email,
		Phone:        u.Phone,
		UserType:     u.UserType,
		Token:        resp.UserModel.Token,
		RefreshToken: resp.UserModel.RefreshToken,
		CreatedAt:    resp.UserModel.CreatedAt.AsTime(),
		UpdatedAt:    resp.UserModel.UpdatedAt.AsTime(),
	}, nil

}

func (c *Client) DeleteUser(ctx context.Context, id string) (*string, error) {
	resp, err := c.service.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &resp.Msg, nil
}

func (c *Client) GetUsers(ctx context.Context, limit int64, offset int64) (*UserModelList, error) {
	resp, err := c.service.GetUsers(ctx, &pb.GetAllUserRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	u := resp.UserModel
	users := []UserModel{}
	for _, user := range u {
		newUser := &UserModel{
			ID:           user.User.Id,
			Name:         user.User.Name,
			Password:     user.User.Password,
			Email:        user.User.Email,
			Phone:        user.User.Phone,
			UserType:     user.User.UserType,
			Token:        user.Token,
			RefreshToken: user.RefreshToken,
			CreatedAt:    user.CreatedAt.AsTime(),
			UpdatedAt:    user.UpdatedAt.AsTime(),
		}

		users = append(users, *newUser)
	}

	return &UserModelList{
		TotalPage: int(resp.TotalPages),
		TotalUser: int(resp.TotalUser),
		UserModel: users,
	}, nil
}

func (c *Client) GetAccessToken(ctx context.Context, id string, r_token string) (*string, error) {
	books.LogInfo(slog.LevelInfo, "client", id)
	resp, err := c.service.GetAccessToken(ctx, &pb.NewAccessTokenRequest{
		Id:           id,
		RefreshToken: r_token,
	})
	if err != nil {
		return nil, err
	}

	return &resp.Token, nil
}
