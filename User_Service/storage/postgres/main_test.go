package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/mahmud3253/Project/User_Service/config"
	"github.com/mahmud3253/Project/User_Service/pkg/db"
	"github.com/mahmud3253/Project/User_Service/pkg/logger"
)

var repo *userRepo

func TestMain(m *testing.M) {
	cfg := config.Load()

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	repo = NewUserRepo(connDB)

	os.Exit(m.Run())
}
