package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct { // cli 명령어들은 각각 cliCommand 구조체에 정보 저장
	name        string
	description string
	callback    func() error
}

// commandMap := map[string]cliCommand{ // @@@ :=는 함수 안에서만 사용 가능
var commandMap map[string]cliCommand = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	// "help": {
	// 	name:        "help",
	// 	description: "Displays a help message",
	// 	callback:    commandHelp,
	// }, commandHelp 안에 commanMap이 쓰이고 그 commandMap 안에 commandHelp가 있음 : initialization cycle for commandMap
}

func reply() {
	commandMap["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := cleanInput(input)
		command, ok := commandMap[cleaned[0]]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0) // 프로그램 종료
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range commandMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}
