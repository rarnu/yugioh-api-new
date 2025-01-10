package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/server"
	"testing"
	"ygoapi/config"
	"ygoapi/database"
)

func TestCardDetail(t *testing.T) {
	server.InitServer()
	config.LoadDatabaseConfig()
	database.NewYgoName()
	cd, _ := database.YgoName.GetCardEffectDetail(14558127)
	fmt.Printf("data = %+v\n", cd)
}
