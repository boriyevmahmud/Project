package postgres

import (
	pb "github.com/mahmud3253/Project/Post_Service/genproto"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type postRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) Create(post *pb.Post) (*pb.Post, error) {
	rPost := pb.Post{}
	query := `INSERT INTO posts (id,name,description,user_id) VALUES($1,$2,$3,$4) returning id,name,description,user_id`
	err := r.db.QueryRow(query, post.Id, post.Name, post.Description, post.UserId).Scan(
		&rPost.Id,
		&rPost.Name,
		&rPost.Description,
		&rPost.UserId)
	if err != nil {
		return nil, err
	}
	for _, media := range post.Medias {
		id := uuid.New()
		_, err := r.db.Exec(`INSERT INTO post_medias(id,type,link,post_id) VALUES ($1,$2,$3,$4)`, id, media.Type, media.Link, rPost.Id)
		if err != nil {
			return nil, err
		}
	}
	return &rPost, err
}

func (r *postRepo) GetById(req *pb.GetByUserIdRequest) (*pb.Post, error) {
	rPost := pb.Post{}

	err := r.db.QueryRow("SELECT id,name,description,user_id from posts WHERE id=$1", req.UserId).Scan(
		&rPost.Id,
		&rPost.Name,
		&rPost.Description,
		&rPost.UserId,
	)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query("SELECT id, type, link from post_medias WHERE post_id = $1", rPost.Id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var media pb.Media
		err := rows.Scan(
			&media.Id,
			&media.Type,
			&media.Link,
		)
		if err != nil {
			return nil, err
		}
		rPost.Medias = append(rPost.Medias, &media)
	}
	return &rPost, nil
}

func (r *postRepo) GetAllUserPosts(userID string) ([]*pb.Post, error) {
	var (
		posts []*pb.Post
	)

	rows, err := r.db.Query("SELECT id, name, description, user_id from posts WHERE user_id = $1", userID)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post pb.Post
		err := rows.Scan(
			&post.Id,
			&post.Name,
			&post.Description,
			&post.UserId,
		)
		if err != nil {
			return nil, err
		}

		var medias []*pb.Media
		rows, err := r.db.Query("SELECT id, type, link from post_medias WHERE post_id = $1", post.Id)

		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var media pb.Media
			err := rows.Scan(
				&media.Id,
				&media.Type,
				&media.Link,
			)
			if err != nil {
				return nil, err
			}

			post.Medias = append(medias, &media)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *postRepo) DeleteByIdPost(req *pb.DeleteByIdPostreq) (*pb.Empty, error) {
	_, err := r.db.Exec("delete from posts WHERE user_id=$1", req.PostId)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (r *postRepo) UpdateByIdPost(req *pb.UpdateByIdPostreq) (*pb.Empty, error) {
	var id string
	for _, post1 := range req.Post {

		query := `update posts set name=$1,description=$2 where user_id=$3 returning id`
		err := r.db.QueryRow(query, post1.Name, post1.Description, req.PostId).Scan(&id)

		if err != nil {
			return nil, err
		}
		for _, media := range post1.Medias {
			_, err := r.db.Exec(`update post_medias set type=$1,link=$2 where post_id=$3`, media.Type, media.Link, id)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil

}
