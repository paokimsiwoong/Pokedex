package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func commandExit(configptr *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0) // 프로그램 종료
	return nil
}

func commandHelp(configptr *config) error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println("")

	return nil
}

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

func commandMap(configptr *config) error {
	// GET request 생성
	req, err := http.NewRequest("GET", configptr.Next, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// default client로 request 전송 후 response 받기
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close() // defer 이용해 Close 안 빼먹게 하기

	// json 데이터를 LocationAreaPage 구조체에 저장
	var page LocationAreaPage //decode한 json을 담을 구조체
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}

	// 다음, 이전 내용을 불러올 url들을 cofig 구조체에 저장
	configptr.Next = page.Next
	configptr.Previous = page.Prev

	// 불러온 에어리어 이름들 출력
	for _, area := range page.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapBack(configptr *config) error {
	// map 명령어가 한번도 사용된 적 없을 때
	if configptr.Previous == "0" {
		fmt.Println("You must call map command at least once to call mapb command")
		return nil
	}
	// 첫페이지를 보고 있을 때
	if configptr.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	// GET request 생성
	req, err := http.NewRequest("GET", configptr.Previous, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// default client로 request 전송 후 response 받기
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close() // defer 이용해 Close 안 빼먹게 하기

	// json 데이터를 LocationAreaPage 구조체에 저장
	var page LocationAreaPage //decode한 json을 담을 구조체
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}

	// 다음, 이전 내용을 불러올 url들을 cofig 구조체에 저장
	configptr.Next = page.Next
	configptr.Previous = page.Prev

	// 불러온 에어리어 이름들 출력
	for _, area := range page.Results {
		fmt.Println(area.Name)
	}

	return nil
}
