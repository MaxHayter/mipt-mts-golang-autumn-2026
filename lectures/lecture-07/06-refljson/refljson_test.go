package refljson

import (
	"encoding/json"
	"sync"
	"testing"
)

type UserID int

type User struct {
	ID       UserID `json:"id"`
	Login    string `json:"login"`
	RealName string `json:"-"`
	Email    string `json:"email,omitempty"`
	Admin    bool   `json:"admin"`
	Bio      string
	password string
}

func testUsers() []User {
	return []User{
		{ID: 42, Login: "rob", RealName: "Rob", Email: "rob@go.dev", Admin: true, Bio: "Привет, \"мир\"\n\t", password: "x"},
		{ID: 7, Login: "ken"}, // пустой Email — omitempty его выкидывает
	}
}

// Результат должен совпадать с encoding/json байт в байт.
func TestMarshalMatchesStdlib(t *testing.T) {
	for _, u := range testUsers() {
		got, err := Marshal(&u)
		if err != nil {
			t.Fatal(err)
		}
		want, _ := json.Marshal(&u)
		if string(got) != string(want) {
			t.Fatalf("\n got: %s\nwant: %s", got, want)
		}
		t.Logf("%s", got)
	}
}

// Второй вызов идёт по плану из кеша — он тоже должен совпасть.
func TestEncoderMatchesStdlib(t *testing.T) {
	var e Encoder
	for _, u := range testUsers() {
		want, _ := json.Marshal(&u)
		for range 2 {
			got, err := e.Marshal(&u)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(want) {
				t.Fatalf("\n got: %s\nwant: %s", got, want)
			}
		}
	}
}

// Кеш общий для горутин: go test -race ./code/06-refljson.
func TestEncoderConcurrent(t *testing.T) {
	var e Encoder
	u := testUsers()[0]
	want, _ := json.Marshal(&u)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			got, err := e.Marshal(&u)
			if err != nil || string(got) != string(want) {
				t.Errorf("got %s, %v", got, err)
			}
		})
	}
	wg.Wait()
}

func TestMarshalRejectsNonStruct(t *testing.T) {
	var e Encoder
	if _, err := Marshal(42); err == nil {
		t.Fatal("Marshal: ожидали ошибку для int")
	}
	if _, err := e.Marshal(42); err == nil {
		t.Fatal("Encoder: ожидали ошибку для int")
	}
}

func TestEncoderRejectsUnsupported(t *testing.T) {
	var e Encoder
	type bad struct{ Tags []string }
	if _, err := e.Marshal(bad{}); err == nil {
		t.Fatal("ожидали ошибку для []string")
	}
	if _, err := e.Marshal(bad{}); err == nil {
		t.Fatal("ошибочный план не должен попасть в кеш")
	}
}
