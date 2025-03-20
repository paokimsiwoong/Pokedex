package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := cleanInput(input)
		fmt.Printf("Your command was: %s\n", cleaned[0])
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
