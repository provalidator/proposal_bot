package models_test

import (
	"testing"

	"github.com/vote_bot/models"
	"github.com/vote_bot/util"
)

func TestDB(t *testing.T) {
	util.Init()
	models.ConnectDatabase()
}
