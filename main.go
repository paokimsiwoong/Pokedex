package main

import (
	"time"

	"github.com/paokimsiwoong/Pokedex/internal/pokeapi"
)

func main() {
	newClient := pokeapi.NewClient(5*time.Second, 5*time.Second) // 첫번째 인자: 연결에 5초 이상 걸리면 timeout, 두번째 인자: cache 저장 시간
	configptr := &config{
		Next:     "next",
		Previous: "prev",
		client:   newClient,
	}
	reply(configptr)
} // @@@ 예시처럼 main.go 파일에서 복잡한 함수들 다른 파일로 옮기기
