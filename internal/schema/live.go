package schema

import (
	"fmt"
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/common/model"
	"time"
)

// ChzzkAPIResponse 치지직 API 공통 응답 구조체
type ChzzkAPIResponse struct {
	Code    int              `json:"code"`
	Message *string          `json:"message"`
	Content LiveListResponse `json:"content"`
}

// LiveListResponse 라이브 목록 API 응답 내용
type LiveListResponse struct {
	Data []LiveData `json:"data"`
	Page PageInfo   `json:"page"`
}

// LiveData 개별 라이브 스트림 정보
type LiveData struct {
	LiveID                int      `json:"liveId"`
	LiveTitle             string   `json:"liveTitle"`
	LiveThumbnailImageURL string   `json:"liveThumbnailImageUrl"`
	ConcurrentUserCount   int      `json:"concurrentUserCount"`
	OpenDate              string   `json:"openDate"`
	Adult                 bool     `json:"adult"`
	Tags                  []string `json:"tags"`
	CategoryType          string   `json:"categoryType"`
	LiveCategory          string   `json:"liveCategory"`
	LiveCategoryValue     string   `json:"liveCategoryValue"`
	ChannelID             string   `json:"channelId"`
	ChannelName           string   `json:"channelName"`
	ChannelImageURL       string   `json:"channelImageUrl"`
}

type LiveDataList []LiveData

// PageInfo 페이징 정보
type PageInfo struct {
	Next string `json:"next"`
}

// GetOpenDateTime OpenDate 문자열을 time.Time으로 변환
func (l *LiveData) GetOpenDateTime() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", l.OpenDate)
}

// IsAdultContent 성인 콘텐츠 여부 확인
func (l *LiveData) IsAdultContent() bool {
	return l.Adult
}

// HasTags 태그가 있는지 확인
func (l *LiveData) HasTags() bool {
	return len(l.Tags) > 0
}

func (l *LiveData) ConvertLiveDataToLiveDataDB() (model.LiveDataDB, error) {
	openDate, err := l.GetOpenDateTime()
	if err != nil {
		return model.LiveDataDB{}, fmt.Errorf("failed to parse open date: %w", err)
	}

	return model.LiveDataDB{
		LiveID:              l.LiveID,
		LiveTitle:           l.LiveTitle,
		ConcurrentUserCount: l.ConcurrentUserCount,
		OpenDate:            openDate,
		Adult:               l.Adult,
		Tags:                model.Tags(l.Tags),
		CategoryType:        l.CategoryType,
		LiveCategory:        l.LiveCategory,
		LiveCategoryValue:   l.LiveCategoryValue,
		ChannelID:           l.ChannelID,
		ChannelName:         l.ChannelName,
		ChannelImageURL:     l.ChannelImageURL,
	}, nil
}

func ConvertLiveListToLiveDataDBList(liveDataList LiveDataList) (model.LiveDataDBList, error) {
	var liveDataDBList model.LiveDataDBList
	for _, liveData := range liveDataList {
		liveDataDB, err := liveData.ConvertLiveDataToLiveDataDB()
		if err != nil {
			return nil, fmt.Errorf("failed to convert live data: %w", err)
		}
		liveDataDBList = append(liveDataDBList, &liveDataDB)
	}
	return liveDataDBList, nil
}
