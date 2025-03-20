package main

import (
	"fmt"
	"os"
)

// exit 명령어
func commandExit(configptr *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0) // 프로그램 종료
	return nil
}

// help 명령어
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

// map 명령어
func commandMap(configptr *config) error {
	// 주어진 주소로 page 받아오기
	page, err := configptr.client.GetAreaPage(configptr.Next)
	if err != nil {
		return fmt.Errorf("error getting area page: %w", err)
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

// mapb 명령어
func commandMapBack(configptr *config) error {
	// map 명령어가 한번도 사용된 적 없을 때
	if configptr.Previous == "prev" {
		fmt.Println("You must call map command at least once to call mapb command")
		return nil
	}
	// 첫페이지를 보고 있을 때
	if configptr.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	// 주어진 주소로 page 받아오기
	page, err := configptr.client.GetAreaPage(configptr.Previous)
	if err != nil {
		return fmt.Errorf("error getting area page: %w", err)
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
