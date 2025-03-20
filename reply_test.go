package main

// 유닛테스트를 진행하는 파일의 이름은 반드시 ~~_test.go 형태여야 한다

import (
	"fmt"
	// "reflect" // .DeepEqual(x, y) 를 제공하는 모듈 (slice끼리 비교하기 위해 사용)
	"testing" // 유닛테스트에 사용되는 모듈
)

func TestCleanInput(t *testing.T) { // 테스트 함수의 이름은 반드시 Test~~~~ 형태여야 한다
	type testCase struct {
		input    string
		expected []string
	}

	cases := []testCase{
		{" test case  ", []string{"test", "case"}},
		{"Charmander Bulbasaur PIKACHU", []string{"charmander", "bulbasaur", "pikachu"}},
		{"", []string{}}, // @@@@ reflect.DeepEqual([], [])는 false가 나온다 https://github.com/golang/go/issues/42265
	}

	passCount := 0
	failCount := 0

	for _, c := range cases { // case는 예약어라 사용하면 안됨
		actual := cleanInput(c.input) // func cleanInput(text string) []string

		// if !reflect.DeepEqual(actual, c.expected) {
		// @@@@ [""]과 []를 비교해 길이가 달라 false가 나온 것이었음 이제 DeepEqual도 정상 작동함 <=== XXX reflect.DeepEqual([], [])는 false가 나온다 https://github.com/golang/go/issues/42265
		if !Equal(actual, c.expected) {
			failCount++
			t.Errorf(
				`---------------------------------
Inputs:     "%v"
Expecting:  %v
Actual:     %v
Fail
`,
				c.input, c.expected, actual,
			)
		} else {
			passCount++
			fmt.Printf(
				`---------------------------------
Inputs:     "%v"
Expecting:  %v
Actual:     %v
Pass
`,
				c.input, c.expected, actual,
			)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

// https://yourbasic.org/golang/compare-slices/
// 두 슬라이스 비교하는 함수
func Equal[T comparable](a, b []T) bool {
	// @@@ T any로 두면 T any 에는 ==나 !=가 사용 불가능한 임의의 타입들을 다 포함하기 때문에 밑의 if a[i] != b[i] 라인에서 invalid operation: a[i] != b[i] (incomparable types in type set) 에러 발생
	// @@@ ===> predeclared idnetifier인 comparable 사용
	// @@@ (comparable is an interface that is implemented by all comparable types (booleans, numbers, strings, pointers, channels, arrays of comparable types, structs whose fields are all comparable types))
	// @@@ https://stackoverflow.com/questions/68053957/go-with-generics-type-parameter-t-is-not-comparable-with
	// fmt.Printf("len(a): %d, len(b) %d", len(a), len(b))
	if len(a) != len(b) {
		return false
	}

	// for i := 0; i < len(a); i++ {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
