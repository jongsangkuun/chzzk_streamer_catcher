package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/chzzk"
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/common/conf"
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/common/model"
	customLogger "github.com/jongsangkuun/chzzk_streamer_catcher/internal/log"
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/schema"
	"io/ioutil"
	"net/http"
	"net/url"
)

func CatcherService(env conf.Env, dbConn *sql.DB) (schema.LiveDataList, error) {
	client := &http.Client{}
	headers := map[string]string{
		"Client-Id":     env.ChzzkClientId,
		"Client-Secret": env.ChzzkSecretId,
	}

	// 첫 번째 페이지 요청
	initialResponse, err := makeRequest(client, chzzk.LiveUrl, "", headers)
	if err != nil {
		return nil, fmt.Errorf("첫 번째 페이지 요청 실패: %w", err)
	}

	var allData schema.LiveDataList
	allData = append(allData, initialResponse.Content.Data...)

	// 페이지네이션 처리
	nextPageToken := initialResponse.Content.Page.Next
	for nextPageToken != "" {
		customLogger.Info("다음 페이지 요청 시작...")

		nextResponse, err := makeRequest(client, chzzk.LiveUrl, nextPageToken, headers)
		if err != nil {
			customLogger.Fatal("다음 페이지 요청 실패: ", err)
		}

		allData = append(allData, nextResponse.Content.Data...)
		nextPageToken = nextResponse.Content.Page.Next
	}

	customLogger.Info("다음 페이지가 없습니다.")
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("Live data count: ", len(allData))

	liveDBList, err := schema.ConvertLiveListToLiveDataDBList(allData)
	if err != nil {
		return nil, err
	}

	err = model.BulkInsert(dbConn, liveDBList)
	if err != nil {
		return nil, err
	}
	return allData, nil
}

func makeRequest(client *http.Client, baseURL string, nextPageToken string, headers map[string]string) (*schema.ChzzkAPIResponse, error) {
	requestURL, err := buildRequestURL(baseURL, nextPageToken)
	if err != nil {
		return nil, fmt.Errorf("URL 생성 실패: %w", err)
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	// 헤더 설정
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("요청 실행 실패: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %w", err)
	}

	var apiResponse schema.ChzzkAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("JSON 파싱 실패: %w", err)
	}

	return &apiResponse, nil
}

func buildRequestURL(baseURL string, nextPageToken string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if nextPageToken != "" {
		q := u.Query()
		q.Set("next", nextPageToken)
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
