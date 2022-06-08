package repo

import (
	pb "Project/template-service/genproto"
	user "Project/template-service/genproto"
)

//UserStorageI ...
type PostStorageI interface {
	Create(*pb.Post) (*pb.Post, error)
	GetById(*pb.GetByUserIdRequest) (*pb.Post, error)
	GetAllUserPosts(userID string) ([]*pb.Post, error)
	DeleteByIdPost(*pb.DeleteByIdPostreq) (*user.Empty, error)
	UpdateByIdPost(*pb.UpdateByIdPostreq)(*pb.Empty,error)
}
