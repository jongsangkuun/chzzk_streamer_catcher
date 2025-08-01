package database

import (
	"fmt"
	"log"
)

// 연결 테스트
func TestConnection() error {
	if DB == nil {
		return fmt.Errorf("데이터베이스 연결이 없습니다")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("DB 인스턴스 획득 실패: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("데이터베이스 핑 실패: %v", err)
	}

	log.Println("데이터베이스 연결 테스트 성공")
	return nil
}
