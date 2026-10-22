package model

import (
	"encoding/json"
	"testing"
)

// plainUser — те же поля, но без методов: новый именованный тип их не наследует,
// поэтому json.Marshal пойдёт через reflect, а не через наш MarshalJSON.
type plainUser User

func TestGeneratedMatchesStdlib(t *testing.T) {
	u := User{ID: 1, Login: "rob", Name: "Роб \"Pike\"", Email: "rob@go.dev", Age: 70, Admin: true, Password: "x"}

	got, err := u.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	want, _ := json.Marshal(plainUser(u))
	if string(got) != string(want) {
		t.Fatalf("\n got: %s\nwant: %s", got, want)
	}

	// encoding/json сам находит сгенерированный метод через интерфейс json.Marshaler.
	viaStd, _ := json.Marshal(u)
	if string(viaStd) != string(want) {
		t.Fatalf("json.Marshal(u) = %s", viaStd)
	}
	t.Logf("%s", got)
}

// Пустой Email с omitempty в JSON не попадает — как у encoding/json.
func TestOmitEmpty(t *testing.T) {
	u := User{ID: 2, Login: "ken", Age: 83}

	got, err := u.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	want, _ := json.Marshal(plainUser(u))
	if string(got) != string(want) {
		t.Fatalf("\n got: %s\nwant: %s", got, want)
	}
	t.Logf("%s", got)
}
