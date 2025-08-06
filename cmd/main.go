package main

import (
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/common/conf"
	customLogger "github.com/jongsangkuun/chzzk_streamer_catcher/internal/log"
	"github.com/jongsangkuun/chzzk_streamer_catcher/pkg/service"
	"time"
)

// Todo
// 수집 데이터 Postgres Bulk Insert로 수집 기능 추가
// 리펙토링 필수....
func main() {
	customLogger.Init()

	envConfig, err := conf.ParseEnv()
	if err != nil {
		customLogger.Fatal("Failed to parse environment variables: ", err)
	}

	dbConn, err := conf.ConnectPostgreSQL(envConfig)
	if err != nil {
		customLogger.Fatal("Failed to connect to PostgreSQL: ", err)
	}
	defer conf.CloseConnection()

	err = conf.InitializeDatabase(dbConn)
	if err != nil {
		customLogger.Fatal("Failed to initialize database: ", err)
	}
	for {
		_, err = service.CatcherService(envConfig, dbConn)
		if err != nil {
			customLogger.Fatal("Failed to catch live list: ", err)
		}
		time.Sleep(300 * time.Second)
	}
}
