package postgres

import (
	"fmt"

	pb "github.com/mahmud3253/Project/User_Service/genproto"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *pb.User) (*pb.User, error) {
	var (
		rUser = pb.User{}
	)
	id1, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	fmt.Println(id1)
	err = r.db.QueryRow("INSERT INTO users (id, first_name, last_name) VALUES($1, $2, $3) RETURNING id, first_name, last_name", id1, user.FirstName, user.LastName).Scan(
		&rUser.Id,
		&rUser.FirstName,
		&rUser.LastName,
	)
	if err != nil {
		return &pb.User{}, err
	}

	return &rUser, nil
}

func (r *userRepo) GetByIdUser(ID string) (*pb.User, error) {
	var (
		rUser = pb.User{}
	)

	err := r.db.QueryRow("SELECT id, first_name, last_name from users WHERE id = $1", ID).Scan(
		&rUser.Id,
		&rUser.FirstName,
		&rUser.LastName,
	)
	if err != nil {
		return nil, err
	}

	return &rUser, err
}

func (r *userRepo) GetAllUsers() (*pb.GetAllUser, error) {
	rUser := pb.GetAllUser{}
	rows, err := r.db.Query("select id,first_name,last_name from users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var data pb.User
		err = rows.Scan(&data.Id, &data.FirstName, &data.LastName)
		if err != nil {
			return nil, err
		}

		rUser.Users = append(rUser.Users, &data)
	}
	return &rUser, nil
}

func (r *userRepo) DeleteById(ID string) (*pb.Empty, error) {
	_, err := r.db.Exec("delete from users where id=$1", ID)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (r *userRepo) UpdateById(req *pb.UpdateByIdReq) (*pb.Empty, error) {
	_, err := r.db.Exec("UPDATE users SET first_name=$1,last_name=$2 where id=$3", req.Users.FirstName, req.Users.LastName, req.UserId)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (r *userRepo) ListUser(req *pb.ListUserReq) (*pb.ListUserResponse, error) {
	rUser := pb.GetAllUser{}

	offset := (req.Page - 1) * req.Limit

	rows, err := r.db.Query("select id,first_name,last_name from users  OFFSET $1 LIMIT $2", offset, req.Limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var data pb.User
		err = rows.Scan(&data.Id, &data.FirstName, &data.LastName)
		if err != nil {
			return nil, err
		}
		rUser.Users = append(rUser.Users, &data)
	}
	count := 0
	countQuery := `SELECT count(*)from users`
	err = r.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}
	return &pb.ListUserResponse{
		Users: rUser.Users,
		Count: int64(count),
	}, nil
}
