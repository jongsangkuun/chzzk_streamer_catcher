package conf

import (
	"database/sql"
	"fmt"
	customLogger "github.com/jongsangkuun/chzzk_streamer_catcher/internal/log"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// PostgreSQL 연결
func ConnectPostgreSQL(config Env) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s",
		config.PsqlHost,
		config.PsqlUser,
		config.PsqlPassword,
		config.PsqlDb,
		config.PsqlPort,
		config.SSLMode,  // SSL 모드 설정
		config.Timezone, // 타임존 설정
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL 연결 실패: %v", err)
	}

	// Connection Pool 설정
	setupConnectionPool(db, config)

	// 연결 테스트
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("PostgreSQL 연결 테스트 실패: %v", err)
	}

	DB = db
	customLogger.Info("PostgreSQL 데이터베이스 연결 성공")
	return db, nil
}

// 연결 풀 설정을 위한 구조체
type ConnectionPoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// 기본값 상수들
const (
	DefaultMaxOpenConns    = 25
	DefaultMaxIdleConns    = 5
	DefaultConnMaxLifetime = 60 * time.Minute
	DefaultConnMaxIdleTime = 10 * time.Minute
)

// 환경 설정에서 연결 풀 설정 생성
func createConnectionPoolConfig(config Env) ConnectionPoolConfig {
	poolConfig := ConnectionPoolConfig{}

	// MaxOpenConns 설정
	if maxOpen, err := strconv.Atoi(config.PsqlMaxOpenConns); err == nil && maxOpen > 0 {
		poolConfig.MaxOpenConns = maxOpen
	} else {
		poolConfig.MaxOpenConns = DefaultMaxOpenConns
	}

	// MaxIdleConns 설정
	if maxIdle, err := strconv.Atoi(config.PsqlMaxIdleConns); err == nil && maxIdle > 0 {
		poolConfig.MaxIdleConns = maxIdle
	} else {
		poolConfig.MaxIdleConns = DefaultMaxIdleConns
	}

	// ConnMaxLifetime 설정
	if lifetime, err := time.ParseDuration(config.PsqlConnMaxLifetime); err == nil {
		poolConfig.ConnMaxLifetime = lifetime
	} else {
		poolConfig.ConnMaxLifetime = DefaultConnMaxLifetime
	}

	// ConnMaxIdleTime 설정
	if idleTime, err := time.ParseDuration(config.PsqlConnMaxIdleTime); err == nil {
		poolConfig.ConnMaxIdleTime = idleTime
	} else {
		poolConfig.ConnMaxIdleTime = DefaultConnMaxIdleTime
	}

	return poolConfig
}

// Connection Pool 설정
func setupConnectionPool(sqlDB *sql.DB, config Env) {
	poolConfig := createConnectionPoolConfig(config)

	// Connection Pool 설정 적용
	sqlDB.SetMaxOpenConns(poolConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(poolConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(poolConfig.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(poolConfig.ConnMaxIdleTime)

	customLogger.Info("Connection Pool 설정: MaxOpen=", poolConfig.MaxOpenConns,
		"MaxIdle=", poolConfig.MaxIdleConns,
		"MaxLifetime=", poolConfig.ConnMaxLifetime,
		"MaxIdleTime=", poolConfig.ConnMaxIdleTime)
}

// 연결 닫기
func CloseConnection() error {
	if DB == nil {
		return nil
	}

	if err := DB.Close(); err != nil {
		return fmt.Errorf("데이터베이스 연결 닫기 실패: %v", err)
	}

	customLogger.Info("데이터베이스 연결 종료")
	return nil
}

// Connection Pool 상태 조회
func GetConnectionStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{"error": "데이터베이스 연결이 없습니다"}
	}

	stats := DB.Stats()
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

// 트랜잭션 시작
func BeginTx() (*sql.Tx, error) {
	if DB == nil {
		return nil, fmt.Errorf("데이터베이스 연결이 없습니다")
	}
	return DB.Begin()
}

// 헬스 체크
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("데이터베이스 연결이 없습니다")
	}
	return DB.Ping()
}

// Prepared Statement 생성 헬퍼
func PrepareStmt(query string) (*sql.Stmt, error) {
	if DB == nil {
		return nil, fmt.Errorf("데이터베이스 연결이 없습니다")
	}
	return DB.Prepare(query)
}

// InitializeDatabase SQL 파일을 사용하여 데이터베이스를 초기화
func InitializeDatabase(db *sql.DB) error {
	log.Println("Initializing database from SQL file...")

	// SQL 파일 경로
	sqlPath := filepath.Join("internal", "common", "model", "init.sql")

	// SQL 파일 읽기
	sqlContent, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("failed to read init.sql from %s: %w", sqlPath, err)
	}

	// SQL 실행
	_, err = db.Exec(string(sqlContent))
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %w", err)
	}

	log.Println("Database initialization completed successfully")
	return nil
}
