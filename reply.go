package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/paokimsiwoong/Pokedex/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
	client   pokeapi.Client
	pokedex  map[string]pokeapi.PokemonData
}

type cliCommand struct { // cli 명령어들은 각각 cliCommand 구조체에 정보 저장
	name        string
	description string
	callback    func(*config, ...string) error
}

func reply(configptr *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := cleanInput(input)

		// @@@ 해답처럼 길이가 0인 edgy case 처리
		if len(cleaned) == 0 {
			continue
		}

		command, ok := getCommands()[cleaned[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue // continue로 바로 다음 입력으로 넘어가게 하기 (if block 밖에 새로운 코드가 추가되도 실행되지 않고 다시 입력 단계로 가도록)
		} else {
			if err := command.callback(configptr, cleaned[1:]...); err != nil {
				fmt.Println(err)
			}
			continue
		}
	}
}

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	} // 이부분 없으면 text = ""일 때 split len이 1이 된다

	// fmt.Printf("len(text): %d", len(text))
	trimmed := strings.Trim(text, " ")
	// fmt.Printf("len(trimmed): %d", len(trimmed))
	lowered := strings.ToLower(trimmed)
	// fmt.Printf("lowered: %v, len(lowered): %d", lowered, len(lowered))
	split := strings.Split(lowered, " ")
	// fmt.Printf("split: %v, split[0]: %v, len(split): %d", split, split[0], len(split))
	// @@@ .Split(string, sep)함수의 string에 ""이 들어가면 결과 slice가 [""]가 되어 len이 0이 아니라 1이 된다 ==> print시에는 ""이 생략되어 비어있는 것처럼 []로만 표시됨
	// @@@ https://pkg.go.dev/strings#Split

	return split
}

// commandMap := map[string]cliCommand{ // @@@ :=는 함수 안에서만 사용 가능
// var commandMap map[string]cliCommand = map[string]cliCommand{
// 	"exit": {
// 		name:        "exit",
// 		description: "Exit the Pokedex",
// 		callback:    commandExit,
// 	},
// 	"help": {
// 		name:        "help",
// 		description: "Displays a help message",
// 		callback:    commandHelp,
// 	}, // commandHelp 안에 commanMap이 쓰이고 그 commandMap 안에 commandHelp가 있음 : initialization cycle for commandMap
// }

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location areas in the Pokemon world",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of all the Pokemons in the area. It takes the name of a location as an argument",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon. It takes the name of a Pokemon as an argument",
			callback:    commandCatch,
		},
		// "inspect": {
		// 	name:        "inspect",
		// 	description: "Prints the name, height, weight, stats, and type(s) of the Pokemon. It takes the name of a Pokemon as an argument",
		// 	callback:    commandInspect,
		// },
	}
}

// @@@ 함수 정의 자체는 반환하는 map을 초기화하지 않는다
// @@@ ==> getCommands함수가 정의될 때 help의 callback commandHelp의 정의가 필요하지 않으므로 문제 해결
