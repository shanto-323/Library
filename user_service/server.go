package userservice

import (
	"context"
	"net"

	"github.com/shanto-323/Library/user_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	service Service
	pb.UnimplementedAuthServiceServer
}

func ListenGRPC(s Service, port string) error {
	ls, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterAuthServiceServer(serv, &grpcServer{
		service: s,
	})

	return serv.Serve(ls)
}

func (sr *grpcServer) SignUp(ctx context.Context, r *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	user, err := sr.service.SignUpNewUser(ctx, r.User.Name, r.User.Password, r.User.Email, r.User.Phone, r.User.UserType)
	if err != nil {
		return nil, err
	}
	return &pb.SignUpResponse{
		UserModel: &pb.UserModel{
			User: &pb.User{
				Id:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Phone:    user.Phone,
				Password: user.Password,
				UserType: user.UserType,
			},
			Token:        user.Token,
			RefreshToken: user.RefreshToken,
			CreatedAt:    timestamppb.New(user.CreatedAt),
			UpdatedAt:    timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (sr *grpcServer) SignIn(ctx context.Context, r *pb.SignInRequest) (*pb.SignInResponse, error) {
	user, err := sr.service.LoginWithEmail(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}
	return &pb.SignInResponse{
		UserModel: &pb.UserModel{
			User: &pb.User{
				Id:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Phone:    user.Phone,
				Password: user.Password,
				UserType: user.UserType,
			},
			Token:        user.Token,
			RefreshToken: user.RefreshToken,
			CreatedAt:    timestamppb.New(user.CreatedAt),
			UpdatedAt:    timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (sr *grpcServer) SignOut(ctx context.Context, r *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	err := sr.service.Logout(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.SignOutResponse{
		Msg: "Sign out success",
	}, err
}

func (sr *grpcServer) UpdateUser(ctx context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	u := r.User
	user, err := sr.service.UpdateUser(ctx, &UserModel{
		ID:       u.Id,
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
		Phone:    u.Phone,
		UserType: u.UserType,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{
		UserModel: &pb.UserModel{
			User: &pb.User{
				Id:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Phone:    user.Phone,
				Password: user.Password,
				UserType: user.UserType,
			},
			Token:        user.Token,
			RefreshToken: user.RefreshToken,
			CreatedAt:    timestamppb.New(user.CreatedAt),
			UpdatedAt:    timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (sr *grpcServer) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := sr.service.DeleteUserById(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{
		Msg: "Delete user success",
	}, err
}

func (sr *grpcServer) GetUsers(ctx context.Context, r *pb.GetAllUserRequest) (*pb.GetAllUserResponse, error) {
	users, err := sr.service.GetAllUser(ctx, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}

	userModelList := []*pb.UserModel{}
	for _, user := range users.UserModel {
		userModel := &pb.UserModel{
			User: &pb.User{
				Id:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Phone:    user.Phone,
				Password: user.Password,
				UserType: user.UserType,
			},
			Token:        user.Token,
			RefreshToken: user.RefreshToken,
			CreatedAt:    timestamppb.New(user.CreatedAt),
			UpdatedAt:    timestamppb.New(user.UpdatedAt),
		}
		userModelList = append(userModelList, userModel)
	}

	return &pb.GetAllUserResponse{
		TotalPages: int64(users.TotalPage),
		TotalUser:  int64(users.TotalUser),
		UserModel:  userModelList,
	}, nil
}
