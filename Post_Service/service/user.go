package service

import (
	pb "Project/template-service/genproto"
	//user "Project/template-service/genproto"
	l "Project/template-service/pkg/logger"
	cl "Project/template-service/service/grpc_client"
	storage "Project/template-service/storage"
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//UserService ...
type PostService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *PostService) Create(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	id := uuid.New()
	req.Id = id.String()
	user, err := s.storage.Post().Create(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostService) GetById(ctx context.Context, req *pb.GetByUserIdRequest) (*pb.Post, error) {
	post, err := s.storage.Post().GetById(req)
	if err != nil {
		return nil, err
	}

	return post, err
}

func (s *PostService) GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error) {
	posts, err := s.storage.Post().GetAllUserPosts(req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserPosts{
		Posts: posts,
	}, err
}

func (s *PostService) DeleteByIdPost(ctx context.Context, req *pb.DeleteByIdPostreq) (*pb.Empty, error) {
	_, err := s.storage.Post().DeleteByIdPost(req)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, err
}

func (s *PostService) UpdateByIdPost(ctx context.Context, req *pb.UpdateByIdPostreq) (*pb.Empty, error) {
	_, err := s.storage.Post().UpdateByIdPost(req)
	if err != nil {
		s.logger.Error("error while updating datas", l.Error(err))
	}
	return &pb.Empty{}, err
}
