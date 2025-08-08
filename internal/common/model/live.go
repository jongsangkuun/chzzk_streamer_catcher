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

// BulkInsert 여러 개의 LiveDataDB를 한 번에 삽입하는 메서드
func BulkInsert(db *sql.DB, liveDataList LiveDataDBList) error {
	if len(liveDataList) == 0 {
		return nil
	}

	// PostgreSQL 최대 매개변수 개수 (65535)를 고려하여 배치 크기 결정
	// 각 행마다 12개의 매개변수가 사용되므로, 65535 / 12 ≈ 5460
	// 안전하게 더 작은 배치 크기를 사용합니다. 예를 들어 1000개
	batchSize := 1000
	if len(liveDataList) < batchSize {
		batchSize = len(liveDataList)
	}

	for i := 0; i < len(liveDataList); i += batchSize {
		end := i + batchSize
		if end > len(liveDataList) {
			end = len(liveDataList)
		}
		batch := liveDataList[i:end]

		if err := insertBatch(db, batch); err != nil {
			return fmt.Errorf("batch insert failed from index %d to %d: %w", i, end, err)
		}
	}

	return nil
}

// insertBatch는 지정된 배치 데이터를 데이터베이스에 삽입합니다.
func insertBatch(db *sql.DB, liveDataList LiveDataDBList) error {
	if len(liveDataList) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(liveDataList))
	valueArgs := make([]interface{}, 0, len(liveDataList)*12) // LiveDataDB 필드 개수에 맞춰 조정

	for i, liveData := range liveDataList {
		start := i * 12 // LiveDataDB 필드 개수 (12개)
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			start+1, start+2, start+3, start+4, start+5, start+6, start+7,
			start+8, start+9, start+10, start+11, start+12))

		tagsJSON, err := liveData.Tags.Value()
		if err != nil {
			return fmt.Errorf("failed to marshal tags for item %d: %w", i, err)
		}

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
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO live_data (
			live_id, live_title, concurrent_user_count, open_date, adult,
			tags, category_type, live_category, live_category_value,
			channel_id, channel_name, channel_image_url
		) VALUES %s`, strings.Join(valueStrings, ","))

	_, err := db.Exec(query, valueArgs...)
	if err != nil {
		// 여기서 오류 발생 시, 어떤 배치에서 문제가 발생했는지 로그를 남기는 것이 도움이 될 수 있습니다.
		return fmt.Errorf("failed to insert batch: %w", err)
	}

	return nil
}
