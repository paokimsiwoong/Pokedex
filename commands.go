package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// exit 명령어
func commandExit(configptr *config, _ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0) // 프로그램 종료
	return nil
}

// help 명령어
func commandHelp(configptr *config, _ ...string) error {
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
func commandMap(configptr *config, _ ...string) error {
	// 주어진 주소로 page 받아오기
	page, err := configptr.client.GetAreaPage(configptr.Next)
	if err != nil {
		return fmt.Errorf("error getting area list page: %w", err)
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
func commandMapBack(configptr *config, _ ...string) error {
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
		return fmt.Errorf("error getting area list page: %w", err)
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

// explore 명령어
func commandExplore(configptr *config, args ...string) error {
	// @@@ 해답처럼 명령어가 잘못 들어올 경우 예외 처리
	if len(args) != 1 {
		return errors.New("wrong command : try explore <areaName>")
	}

	areaName := args[0]

	areaData, err := configptr.client.GetAreaData(areaName)
	if err != nil {
		return fmt.Errorf("error getting area data: %w", err)
	}

	fmt.Printf("Pokemons in %s: \n", areaName)

	for _, encounter := range areaData.PokemonEncounters {
		fmt.Println(" - " + encounter.Pokemon.Name)
	}

	return nil
}

// catch 명령어
func commandCatch(configptr *config, args ...string) error {
	// @@@ 해답처럼 명령어가 잘못 들어올 경우 예외 처리
	if len(args) != 1 {
		return errors.New("wrong command : try catch <PokemonName>")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemonData, err := configptr.client.GetPokemonData(pokemonName)
	if err != nil {
		return fmt.Errorf("error getting pokemon data: %w", err)
	}

	baseEXP := pokemonData.BaseExperience
	fmt.Printf("baseEXP: %d\n", baseEXP)

	// seed를 매번 바꿔주지 않으면 결과가 항상 같다
	// rand.Seed(time.Now().UnixNano()) // https://zerostarting.tistory.com/53 :
	// rand.Seed is deprecated: As of Go 1.20 there is no reason to call Seed with
	// a random value. Programs that call Seed with a known value to get
	// a specific sequence of results should use New(NewSource(seed)) to
	// obtain a local random generator.
	rand.New(rand.NewSource(time.Now().UnixNano())) // https://ditto-dev.tistory.com/80
	randomInt := rand.Intn(700)

	if randomInt > baseEXP {
		fmt.Printf("%s was caught!\n", pokemonName)
		configptr.pokedex[pokemonName] = pokemonData
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

// inspect 명령어
func commandInspect(configptr *config, args ...string) error {
	// @@@ 해답처럼 명령어가 잘못 들어올 경우 예외 처리
	if len(args) != 1 {
		return errors.New("wrong command : try inspect <PokemonName>")
	}

	pokemonName := args[0]

	pokemonData, ok := configptr.pokedex[pokemonName]

	if !ok {
		fmt.Printf("You have not caught %s\n", pokemonName)
	} else {
		fmt.Printf("Name: %s\n", pokemonData.Name)
		fmt.Printf("Height: %d\n", pokemonData.Height)
		fmt.Printf("Weight: %d\n", pokemonData.Weight)

		fmt.Println("Stats:")
		for _, stat := range pokemonData.Stats {
			fmt.Printf(" - %v: %v\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Println("Types:")
		for _, t := range pokemonData.Types { // @@@ type은 예약어이므로 피하기
			fmt.Printf(" - %s\n", t.Type.Name)
		}
	}

	return nil
}

func commandPokeDEX(configptr *config, args ...string) error {
	fmt.Println("Your Pokedex:")

	count := 0

	for key := range configptr.pokedex {
		fmt.Printf(" - %v\n", key)
		count++
	}
	fmt.Printf("You've caught %d Pokemon so far.\n", count)

	return nil
}
