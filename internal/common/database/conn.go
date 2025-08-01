package database

import (
	"database/sql"
	"fmt"
	costomLogger "github.com/jongsangkuun/chzzk_streamer_catcher/internal/log"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
	Timezone     string

	// Connection Pool 설정
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

var DB *gorm.DB

// PostgreSQL 연결
func ConnectPostgreSQL(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Host,
		config.User,
		config.Password,
		config.DatabaseName,
		config.Port,
		config.SSLMode,
		config.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL 연결 실패: %v", err)
	}

	// Connection Pool 설정
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("DB 인스턴스 획득 실패: %v", err)
	}

	setupConnectionPool(sqlDB, config)

	DB = db
	costomLogger.Info("PostgreSQL 데이터베이스 연결 성공")
	return db, nil
}

// Connection Pool 설정
func setupConnectionPool(sqlDB *sql.DB, config Config) {
	// 기본값 설정
	maxOpenConns := config.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 25 // 기본값
	}

	maxIdleConns := config.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 5 // 기본값
	}

	connMaxLifetime := config.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour // 기본값: 1시간
	}

	connMaxIdleTime := config.ConnMaxIdleTime
	if connMaxIdleTime == 0 {
		connMaxIdleTime = time.Minute * 10 // 기본값: 10분
	}

	// Connection Pool 설정 적용
	sqlDB.SetMaxOpenConns(maxOpenConns)       // 최대 열린 연결 수
	sqlDB.SetMaxIdleConns(maxIdleConns)       // 최대 유휴 연결 수
	sqlDB.SetConnMaxLifetime(connMaxLifetime) // 연결 최대 생명주기
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime) // 유휴 연결 최대 시간

	costomLogger.Info("Connection Pool 설정: MaxOpen=", maxOpenConns, "MaxIdle=", maxOpenConns, "MaxLifetime=", connMaxLifetime, "MaxIdleTime=", connMaxIdleTime)
}

// 연결 닫기
func CloseConnection() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("DB 인스턴스 획득 실패: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("데이터베이스 연결 닫기 실패: %v", err)
	}

	costomLogger.Info("데이터베이스 연결 종료")
	return nil
}

// Connection Pool 상태 조회
func GetConnectionStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{"error": "데이터베이스 연결이 없습니다"}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
