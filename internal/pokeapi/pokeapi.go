package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	locationURL = "https://pokeapi.co/api/v2/location-area/"
)

// GET "https://pokeapi.co/api/v2/location-area/" 로 받은 response 데이터를 저장할 구조체
type LocationAreaPage struct {
	Count   int            `json:"count"`
	Next    string         `json:"next"`
	Prev    string         `json:"previous"`
	Results []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// 서버에 request를 보낼 client를 담고 있는 구조체 (http.Client를 바로 쓰는 대신 구조체에 집어넣어서 method를 활용)
type Client struct {
	httpClient http.Client
}

// 새로운 client를 생성해 반환하는 함수
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

// Client 구조체에 연결된 method
func (c *Client) GetAreaPage(url string) (LocationAreaPage, error) {
	// map 최초 실행시 url == "next"
	if url == "next" {
		url = locationURL
	}

	// GET request 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaPage{}, fmt.Errorf("error creating request: %w", err)
	}

	// default client로 request 전송 후 response 받기
	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaPage{}, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close() // defer 이용해 Close 안 빼먹게 하기

	// json 데이터를 LocationAreaPage 구조체에 저장
	var page LocationAreaPage //decode한 json을 담을 구조체
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return LocationAreaPage{}, fmt.Errorf("error decoding response body: %w", err)
	}
	return page, nil
}
