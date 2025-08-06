package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// LiveDataDB 데이터베이스 저장용 구조체
type LiveDataDB struct {
	ID                  uint
	LiveID              int
	LiveTitle           string
	ConcurrentUserCount int
	OpenDate            time.Time
	Adult               bool
	Tags                Tags
	CategoryType        string
	LiveCategory        string
	LiveCategoryValue   string
	ChannelID           string
	ChannelName         string
	ChannelImageURL     string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type LiveDataDBList []*LiveDataDB

// Tags JSON 배열을 처리하기 위한 커스텀 타입
type Tags []string

// Value GORM을 위한 Value 인터페이스 구현
func (t Tags) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan GORM을 위한 Scan 인터페이스 구현
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into Tags", value)
	}

	return json.Unmarshal(bytes, t)
}

func (liveDataDB *LiveDataDB) Insert(db *sql.DB) error {
	query := `
		INSERT INTO live_data (
			live_id, live_title, concurrent_user_count, open_date, adult, 
			tags, category_type, live_category, live_category_value, 
			channel_id, channel_name, channel_image_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	tagsJSON, err := liveDataDB.Tags.Value()
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = db.Exec(query,
		liveDataDB.LiveID,
		liveDataDB.LiveTitle,
		liveDataDB.ConcurrentUserCount,
		liveDataDB.OpenDate,
		liveDataDB.Adult,
		tagsJSON,
		liveDataDB.CategoryType,
		liveDataDB.LiveCategory,
		liveDataDB.LiveCategoryValue,
		liveDataDB.ChannelID,
		liveDataDB.ChannelName,
		liveDataDB.ChannelImageURL,
		liveDataDB.CreatedAt,
		liveDataDB.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert live data: %w", err)
	}

	return nil
}

// BulkInsert 여러 개의 LiveDataDB를 한 번에 삽입하는 메서드
func BulkInsert(db *sql.DB, liveDataList LiveDataDBList) error {
	if len(liveDataList) == 0 {
		return nil
	}

	// VALUES 부분을 동적으로 생성
	valueStrings := make([]string, 0, len(liveDataList))
	valueArgs := make([]interface{}, 0, len(liveDataList)*14)

	for i, liveData := range liveDataList {
		// 각 레코드마다 ($1, $2, ..., $14) 형태로 placeholders 생성
		start := i * 14
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			start+1, start+2, start+3, start+4, start+5, start+6, start+7,
			start+8, start+9, start+10, start+11, start+12, start+13, start+14))

		// Tags를 JSON으로 변환
		tagsJSON, err := liveData.Tags.Value()
		if err != nil {
			return fmt.Errorf("failed to marshal tags for item %d: %w", i, err)
		}

		// 파라미터 추가
		valueArgs = append(valueArgs,
			liveData.LiveID,
			liveData.LiveTitle,
			liveData.ConcurrentUserCount,
			liveData.OpenDate,
			liveData.Adult,
			tagsJSON,
			liveData.CategoryType,
			liveData.LiveCategory,
			liveData.LiveCategoryValue,
			liveData.ChannelID,
			liveData.ChannelName,
			liveData.ChannelImageURL,
			liveData.CreatedAt,
			liveData.UpdatedAt,
		)
	}

	// 최종 쿼리 생성
	query := fmt.Sprintf(`
		INSERT INTO live_data (
			live_id, live_title, concurrent_user_count, open_date, adult, 
			tags, category_type, live_category, live_category_value, 
			channel_id, channel_name, channel_image_url, created_at, updated_at
		) VALUES %s`, strings.Join(valueStrings, ","))

	_, err := db.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert live data: %w", err)
	}

	return nil
}
