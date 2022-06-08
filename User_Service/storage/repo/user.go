package repo

import (
	pb "Project/template-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
	GetByIdUser(id string) (*pb.User, error)
	GetAllUsers() (*pb.GetAllUser, error)
	DeleteById(id string) (*pb.Empty, error)
	UpdateById(*pb.UpdateByIdReq) (*pb.Empty, error)
	ListUser(*pb.ListUserReq) (*pb.ListUserResponse, error)
}
