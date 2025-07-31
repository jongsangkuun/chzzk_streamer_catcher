package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// LiveDataDB 데이터베이스 저장용 구조체
type LiveDataDB struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	LiveID                int       `gorm:"uniqueIndex" json:"liveId"`
	LiveTitle             string    `json:"liveTitle"`
	LiveThumbnailImageURL string    `json:"liveThumbnailImageUrl"`
	ConcurrentUserCount   int       `json:"concurrentUserCount"`
	OpenDate              time.Time `json:"openDate"`
	Adult                 bool      `json:"adult"`
	Tags                  Tags      `gorm:"type:json" json:"tags"`
	CategoryType          string    `json:"categoryType"`
	LiveCategory          string    `json:"liveCategory"`
	LiveCategoryValue     string    `json:"liveCategoryValue"`
	ChannelID             string    `json:"channelId"`
	ChannelName           string    `json:"channelName"`
	ChannelImageURL       string    `json:"channelImageUrl"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

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
