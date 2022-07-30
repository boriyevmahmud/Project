package service

import (
	"context"
	"fmt"

	pb "github.com/mahmud3253/Project/Post_Service/genproto"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	l "github.com/mahmud3253/Project/Post_Service/pkg/logger"
	cl "github.com/mahmud3253/Project/Post_Service/service/grpc_client"
	storage "github.com/mahmud3253/Project/Post_Service/storage"

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

const (
	topic          = "user.user"
	broker1Address = "localhost:9092"
)

func (s *PostService) Consume(ctx context.Context, a *pb.Empty) (*pb.Empty, error) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages

	t := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})

	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := t.ReadMessage(ctx)
		if err != nil {
			return nil, err
		}
		var user pb.User

		err = user.Unmarshal(msg.Value)
		if err != nil {
			return nil, err
		}

		for _, post := range user.Posts {
			id, err := uuid.NewUUID()
			if err != nil {
				s.logger.Error("failed while generating uuid for new post", l.Error(err))
				return nil, status.Error(codes.Internal, "failed while generating uuid")
			}
			post.Id = id.String()
			_, err = s.storage.Post().Create(post)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))

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
