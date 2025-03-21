package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) { // 테스트 함수의 이름은 반드시 Test~~~~ 형태여야 한다
	type testCase struct {
		key   string
		value []byte
	}

	cases := []testCase{
		{"http://mysite.com", []byte("my site is so good")},
		{"http://myexsittee.com", []byte("example    ")},
	}

	passCount := 0
	failCount := 0

	const interval = 5 * time.Second // cache 저장시간

	for _, c := range cases { // case는 예약어라 사용하면 안됨
		cache := NewCache(interval)

		cache.Add(c.key, c.value)
		savedData, ok := cache.Get(c.key)

		if !ok {
			failCount++
			t.Errorf(
				`---------------------------------
Input key:     "%v"
Input value:     "%v"
fail to find key
Fail
`,
				c.key, c.value,
			)
		} else if string(c.value) != string(savedData) {
			failCount++
			t.Errorf(
				`---------------------------------
Input key:     "%v"
Input value(expected):     "%v"
saved data:		"%v"
Fail
`,
				c.key, string(c.value), string(savedData),
			)
		} else {
			passCount++
			fmt.Printf(
				`---------------------------------
Input key:     "%v"
Input value(expected):     "%v"
saved data:		"%v"
Pass
`,
				c.key, string(c.value), string(savedData),
			)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

func TestReapLoop(t *testing.T) {
	type testCase struct {
		key   string
		value []byte
	}

	cases := []testCase{
		{"http://mysite.com", []byte("my site is so good")},
	}

	passCount := 0
	failCount := 0

	const interval = 5 * time.Millisecond // cache 저장시간

	for _, c := range cases { // case는 예약어라 사용하면 안됨
		cache := NewCache(interval)

		cache.Add(c.key, c.value)
		savedData, ok := cache.Get(c.key)

		if !ok {
			failCount++
			t.Errorf(
				`---------------------------------
Input key:     "%v"
Input value:     "%v"
fail to find key
Fail
`,
				c.key, c.value,
			)
		} else {
			passCount++
			fmt.Printf(
				`---------------------------------
Input key:     "%v"
Input value(expected):     "%v"
saved data:		"%v"
Pass
`,
				c.key, string(c.value), string(savedData),
			)
		}

		time.Sleep(interval * 2)

		_, ok = cache.Get(c.key)

		if ok {
			failCount++
			t.Errorf(
				`---------------------------------
Input key:     "%v"
Input value:     "%v"
fail to erase key
Fail
`,
				c.key, c.value,
			)
		} else {
			passCount++
			fmt.Printf(
				`---------------------------------
Input key:     "%v"
Input value:     "%v"
fail to find key
Pass
`,
				c.key, c.value,
			)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

// @@@ 예시 코드
// func TestAddGet(t *testing.T) {
// 	const interval = 5 * time.Second
// 	cases := []struct {
// 		key string
// 		val []byte
// 	}{
// 		{
// 			key: "https://example.com",
// 			val: []byte("testdata"),
// 		},
// 		{
// 			key: "https://example.com/path",
// 			val: []byte("moretestdata"),
// 		},
// 	}

// 	for i, c := range cases {
// 		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
// 			cache := NewCache(interval)
// 			cache.Add(c.key, c.val)
// 			val, ok := cache.Get(c.key)
// 			if !ok {
// 				t.Errorf("expected to find key")
// 				return
// 			}
// 			if string(val) != string(c.val) {
// 				t.Errorf("expected to find value")
// 				return
// 			}
// 		})
// 	}
// }

// func TestReapLoop(t *testing.T) {
// 	const baseTime = 5 * time.Millisecond
// 	const waitTime = baseTime + 5*time.Millisecond
// 	cache := NewCache(baseTime)
// 	cache.Add("https://example.com", []byte("testdata"))

// 	_, ok := cache.Get("https://example.com")
// 	if !ok {
// 		t.Errorf("expected to find key")
// 		return
// 	}

// 	time.Sleep(waitTime)

// 	_, ok = cache.Get("https://example.com")
// 	if ok {
// 		t.Errorf("expected to not find key")
// 		return
// 	}
// }
