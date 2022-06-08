package repo

import (
	pb "github.com/mahmud3253/Project/Post_Service/genproto"
	user "github.com/mahmud3253/Project/Post_Service/genproto"
)

//UserStorageI ...
type PostStorageI interface {
	Create(*pb.Post) (*pb.Post, error)
	GetById(*pb.GetByUserIdRequest) (*pb.Post, error)
	GetAllUserPosts(userID string) ([]*pb.Post, error)
	DeleteByIdPost(*pb.DeleteByIdPostreq) (*user.Empty, error)
	UpdateByIdPost(*pb.UpdateByIdPostreq) (*pb.Empty, error)
}
