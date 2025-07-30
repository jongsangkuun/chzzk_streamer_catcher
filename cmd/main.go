package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	clientID := ""
	clientSecret := ""

	baseURL := "https://openapi.chzzk.naver.com/open/v1/channels"

	// 조회할 채널 ID 목록 (최대 20개)
	channelIDs := []string{"c847a58a1599988f6154446c75366523"}

	// 쿼리 파라미터 세팅
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	query := url.Values{}
	for _, id := range channelIDs {
		query.Add("channelIds", id)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// 헤더 세팅
	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Client-Secret", clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)

	// JSON을 예쁘게 포맷해서 출력
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		// JSON 파싱에 실패하면 원본 텍스트로 출력
		fmt.Printf("Response (raw):\n%s\n", string(body))
	} else {
		// JSON을 예쁘게 포맷해서 출력
		prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			fmt.Printf("Response (raw):\n%s\n", string(body))
		} else {
			fmt.Printf("Response (formatted):\n%s\n", string(prettyJSON))
		}
	}

}
