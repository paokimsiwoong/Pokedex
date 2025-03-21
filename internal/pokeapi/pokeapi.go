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
	locationURL       = "https://pokeapi.co/api/v2/location-area/"
	locationFirstPage = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20" // @@@ 두번째 페이지의 "previous"는 locationURL가 아니라 변경한 locationFirstPage 주소 ==> 원래 주소로 하면 최초실행 map로 첫페이지 부를때 주소와 map->map->mapb로 첫 페이지 부를 때 주소가 달라 cache miss가 되므로 변경
	pokemonURL        = "https://pokeapi.co/api/v2/pokemon/"                         // pokemon endpoint. 뒤에 포켓몬 id나 이름을 붙이면 관련 정보를 얻을 수 있다.
)

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
		url = locationFirstPage
	}

	// cache를 확인해서 이미 있으면 바로 그 저장된 값을 반환
	var page LocationAreaPage //decode한 json을 담을 구조체
	data, err := c.getBytes(url)
	if err != nil {
		return LocationAreaPage{}, fmt.Errorf("error getting []byte data: %w", err)
	}

	// []byte를 LocationAreaPage 구조체로 변환
	if err := json.Unmarshal(data, &page); err != nil {
		return LocationAreaPage{}, fmt.Errorf("error unmarshalling data: %w", err)
	}

	// c.cache.SetCurrent(url) // 현재 페이지 url 저장

	return page, nil
}

// explore 명령어에서 필요한 에어리어 데이터를 받아오는 메소드
func (c *Client) GetAreaData(areaName string) (AreaData, error) {
	url := locationURL + areaName

	var areaData AreaData //decode한 json을 담을 구조체
	data, err := c.getBytes(url)
	if err != nil {
		return AreaData{}, fmt.Errorf("error getting []byte data: %w", err)
	}

	// []byte를 AreaData 구조체로 변환
	if err := json.Unmarshal(data, &areaData); err != nil {
		return AreaData{}, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return areaData, nil
}

// catch 명령어에서 필요한 포켓몬 데이터를 받아오는 메소드
func (c *Client) GetPokemonData(pokemonName string) (PokemonData, error) {
	url := pokemonURL + pokemonName

	var pokemonData PokemonData //decode한 json을 담을 구조체
	data, err := c.getBytes(url)
	if err != nil {
		return PokemonData{}, fmt.Errorf("error getting []byte data: %w", err)
	}

	// []byte를 AreaData 구조체로 변환
	if err := json.Unmarshal(data, &pokemonData); err != nil {
		return PokemonData{}, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return pokemonData, nil
}

func (c *Client) getBytes(url string) ([]byte, error) {
	data, ok := c.cache.Get(url)

	if !ok { // cache에 url로 불러온 데이터가 안들어있을 때 => 새로 request
		fmt.Println("Cache miss!")
		// GET request 생성
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []byte{}, fmt.Errorf("error creating request: %w", err)
		}

		// default client로 request 전송 후 response 받기
		res, err := c.httpClient.Do(req)
		if err != nil {
			return []byte{}, fmt.Errorf("error sending request: %w", err)
		}
		defer res.Body.Close() // defer 이용해 Close 안 빼먹게 하기

		// response에서 []byte 데이터 받아오기
		// data, err := io.ReadAll(res.Body) // @@@ data 뒤에 :=이 되면 이 data는 이 {} 블록안의 새로운 data가 되어 {} 밖의 data와 다른 값, {} 끝나면 삭제
		newData, err := io.ReadAll(res.Body) // @@@ 이름을 newData로 새롭게 한 후 {} 나가기 전에 data = newData로 입력하기
		if err != nil {
			return []byte{}, fmt.Errorf("error reading response body: %w", err)
		}

		// cache에 저장
		c.cache.Add(url, newData)

		// {} 안에서만 존재하는 newData를 바깥의 data에 입력 (@@@ 또는 data, err 모두 기존의 것을 사용하도록 := 대신 = 사용해도 된다)
		data = newData

	} else {
		fmt.Println("Cache hit!")
	}

	return data, nil
}

// // explore 명령어에서 필요한 에어리어 데이터를 받아오는 메소드
// func (c *Client) GetAreaData(areaName string) (AreaData, error) {
// 	currentPageURL := c.cache.GetCurrent()

// 	// cache에서 현재 페이지 데이터 얻기
// 	var page LocationAreaPage //decode한 json을 담을 구조체
// 	data, err := c.getBytes(currentPageURL)
// 	if err != nil {
// 		return AreaData{}, fmt.Errorf("error getting []byte data: %w", err)
// 	}

// 	if err := json.Unmarshal(data, &page); err != nil {
// 		return AreaData{}, fmt.Errorf("error unmarshalling data: %w", err)
// 	}

// 	for _, area := range page.Results {
// 		// 현재 페이지에 입력한 에어리어가 있으면 그 에어리어 정보 불러오기
// 		if area.Name == areaName {
// 			url := area.Url

// 			var areaData AreaData //decode한 json을 담을 구조체
// 			data, err := c.getBytes(url)
// 			if err != nil {
// 				return AreaData{}, fmt.Errorf("error getting []byte data: %w", err)
// 			}

// 			// []byte를 AreaData 구조체로 변환
// 			if err := json.Unmarshal(data, &areaData); err != nil {
// 				return AreaData{}, fmt.Errorf("error unmarshalling data: %w", err)
// 			}

// 			return areaData, nil
// 		}
// 	}

// 	return AreaData{}, fmt.Errorf("error getting area data: You must choose an area from the current map page")
// }
