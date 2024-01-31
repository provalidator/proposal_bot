package modules_test

import (
	"testing"

	"github.com/vote_bot/models"
	"github.com/vote_bot/modules"
	"github.com/vote_bot/util"
)

func TestDB(t *testing.T) {
	util.Init()
	models.GetDBProposals()
	modules.Proposals()

}
