package postgres

import (
	"testing"
	"github.com/mahmud3253/Project/user-service/config"
	pb "github.com/mahmud3253/Project/user-service/genproto"
	"github.com/mahmud3253/Project/user-service/pkg/db"
	"github.com/mahmud3253/Project/user-service/storage/repo"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	CleanupFunc func()
	Repository  repo.UserStorageI
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	pgPool, cleanup := db.ConnectDBForSuite(config.Load())

	suite.Repository = NewUserRepo(pgPool)
	suite.CleanupFunc = cleanup
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *UserRepositoryTestSuite) TestUserCRUD() {

	user1 := pb.User{
		FirstName: "lala",
		LastName:  "akak",
	}
	user, err := suite.Repository.CreateUser(user1)
	suite.Nil(err)

	getUser, err := suite.Repository.GetUserById(user.Id)
	suite.Nil(err)
	suite.NotNil(getUser.FirstName, "user must not be nil")

	getUser.FirstName = "mask2"
	updatedUser, err := suite.Repository.UpdateById(getUser)
	suite.Nil(err)
	suite.NotNil(updatedUser)

	getupdatedUser, err := suite.Repository.GetUserById(getUser.Id)
	suite.Nil(err)
	suite.NotNil(getupdatedUser, "user must not be nil")
	suite.Equal(updatedUser.FirstName, getUser.FirstName, "lala")

	getuserList, _, err := suite.Repository.ListUser(10, 1)
	suite.Nil(err)
	suite.NotNil(getuserList)

	// err=suite.Repository.Delete(getupdatedUser.Id)
	// suite.Nil(err)

}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	suite.CleanupFunc()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
