package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/paokimsiwoong/Pokedex/internal/pokecache"
)

const (
	// locationURL = "https://pokeapi.co/api/v2/location-area/" // @@@ 두번째 페이지의 "previous"는 이 주소가 아니라 변경한 아래 주소
	locationURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20" // @@@ ==> 원래 주소로 하면 최초실행 map로 첫페이지 부를때 주소와 map->map->mapb로 첫 페이지 부를 때 주소가 달라 cache miss가 되므로 변경
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
	cache      pokecache.Cache
}

// 새로운 client를 생성해 반환하는 함수
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

// Client 구조체에 연결된 method
func (c *Client) GetAreaPage(url string) (LocationAreaPage, error) {
	// map 최초 실행시 url == "next"
	if url == "next" {
		url = locationURL
	}

	// cache를 확인해서 이미 있으면 바로 그 저장된 값을 반환
	data, ok := c.cache.Get(url)
	var page LocationAreaPage //decode한 json을 담을 구조체

	if !ok { // cache에 url로 불러온 데이터가 안들어있을 때 => 새로 request
		fmt.Println("Cache miss!")
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
		// if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		// 	return LocationAreaPage{}, fmt.Errorf("error decoding response body: %w", err)
		// } // cache에 []byte 저장을 위해 decoding 두 단계로 나누어서 진행하기
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return LocationAreaPage{}, fmt.Errorf("error reading response body: %w", err)
		}

		// cache에 저장
		c.cache.Add(url, data)

		// []byte를 LocationAreaPage 구조체로 변환
		if err := json.Unmarshal(data, &page); err != nil {
			return LocationAreaPage{}, fmt.Errorf("error unmarshalling data: %w", err)
		}
	} else { // cache에 이미 데이터가 들어 있으면 바로 그 데이터를 []byte를 구조체로 변환
		fmt.Println("Cache hit!")
		if err := json.Unmarshal(data, &page); err != nil {
			return LocationAreaPage{}, fmt.Errorf("error unmarshalling data: %w", err)
		}
	}

	return page, nil
}
