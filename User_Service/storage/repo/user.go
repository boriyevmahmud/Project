package repo

import (
	pb "github.com/mahmud3253/Project/User_Service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
	GetByIdUser(id string) (*pb.User, error)
	GetAllUsers() (*pb.GetAllUser, error)
	DeleteById(id string) (*pb.Empty, error)
	UpdateById(*pb.UpdateByIdReq) (*pb.Empty, error)
	ListUser(*pb.ListUserReq) (*pb.ListUserResponse, error)
	CheckField(field, value string) (bool, error)
	RegisterUser(*pb.CreateUserAuthReqBody) (*pb.CreateUserAuthResBody, error)
}
