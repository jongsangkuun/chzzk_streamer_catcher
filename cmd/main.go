package main

import (
	"encoding/json"
	"fmt"
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/common/env"
	costomLogger "github.com/jongsangkuun/chzzk_streamer_catcher/internal/log"
	"github.com/jongsangkuun/chzzk_streamer_catcher/internal/schema"
	"github.com/jongsangkuun/chzzk_streamer_catcher/pkg/chzzk"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Todo
// 수집 데이터 Postgres Bulk Insert로 수집 기능 추가
func main() {
	costomLogger.Init()

	envConfig, err := env.ParseEnv()
	if err != nil {
		costomLogger.Fatal("Failed to parse environment variables: ", err)
	}

	liveUrl := chzzk.LiveUrl

	clientID := envConfig.ChzzkClientId
	clientSecret := envConfig.ChzzkSecretId

	chzzkAPIResponse, NextPageToken, err := firstNextPageToken(chzzk.LiveUrl, clientID, clientSecret)

	if err != nil {
		costomLogger.Fatal("Failed to get first next page token: ", err)
	}

	for {
		if NextPageToken != "" {
			costomLogger.Info("다음 페이지 요청 시작...")

			// URL에 next 파라미터 추가
			nextU, err := url.Parse(liveUrl)
			if err != nil {
				costomLogger.Fatal("Failed to parse next URL: ", err)
			}

			q := nextU.Query()
			q.Set("next", NextPageToken)
			nextU.RawQuery = q.Encode()

			// 새로운 요청 생성
			nextReq, err := http.NewRequest("GET", nextU.String(), nil)
			if err != nil {
				costomLogger.Fatal("Failed to create next request: ", err)
			}

			// 헤더 세팅
			nextReq.Header.Set("Client-Id", clientID)
			nextReq.Header.Set("Client-Secret", clientSecret)

			client := &http.Client{}
			// 다음 페이지 요청 실행
			nextResp, err := client.Do(nextReq)
			if err != nil {
				costomLogger.Fatal("Next request failed: ", err)
			}
			defer nextResp.Body.Close()

			nextBody, err := ioutil.ReadAll(nextResp.Body)
			if err != nil {
				costomLogger.Fatal("Failed to read next response: %v", err)
			}

			costomLogger.Info("Next Response: ", nextBody)

			// 다음 페이지 응답 파싱
			nextApiResponse := schema.ChzzkAPIResponse{}
			err = json.Unmarshal(nextBody, &nextApiResponse)
			if err != nil {
				costomLogger.Fatal("Failed to unmarshal next response: %v", err)
			}

			NextPageToken = nextApiResponse.Content.Page.Next
			chzzkAPIResponse = append(chzzkAPIResponse, nextApiResponse.Content.Data...)

			costomLogger.Info("Next Live data count: ", len(nextApiResponse.Content.Data))
			costomLogger.Info("Next Next page token: ", nextApiResponse.Content.Page.Next)
			costomLogger.Info("Next page data: ", nextApiResponse.Content.Data)

		} else {
			costomLogger.Info("다음 페이지가 없습니다.")
			break
		}
	}
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("Live data count: ", len(chzzkAPIResponse))
}

func firstNextPageToken(liveUrl string, clientID string, clientSecret string) ([]schema.LiveData, string, error) {
	u, err := url.Parse(liveUrl)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to parse URL: ", err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to create request: ", err)
	}

	// 헤더 세팅
	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Client-Secret", clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("Request failed: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to read response: %v", err)
	}

	// ChzzkAPIResponse 구조체로 언마샬링
	chzzkAPIResponse := schema.ChzzkAPIResponse{}

	err = json.Unmarshal(body, &chzzkAPIResponse)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to unmarshal response: %v", err)
	}

	NextPageToken := chzzkAPIResponse.Content.Page.Next

	costomLogger.Info("Live data count: ", len(chzzkAPIResponse.Content.Data))
	costomLogger.Info("Next Page Token :", NextPageToken)
	costomLogger.Info(chzzkAPIResponse.Content.Data)

	return chzzkAPIResponse.Content.Data, NextPageToken, nil
}
