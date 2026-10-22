package bench

import (
	"encoding/json"
	"testing"

	"github.com/mailru/easyjson"

	refljson "lecture-07/06-refljson"
	"lecture-07/08-codegen/model"
)

// plainUser — те же поля без методов, чтобы encoding/json честно шёл через reflect.
type plainUser model.User

func testUser() model.User {
	return model.User{ID: 42, Login: "rob", Name: "Rob Pike", Email: "rob@go.dev", Age: 70, Admin: true}
}

// Все способы должны давать один и тот же JSON — иначе сравнивать нечего.
// Проверяем и заполненный, и пустой Email: omitempty поддерживают все.
func TestSameOutput(t *testing.T) {
	var enc refljson.Encoder
	empty := testUser()
	empty.Email = ""
	for _, u := range []model.User{testUser(), empty} {
		want, _ := json.Marshal(plainUser(u))
		refl, _ := refljson.Marshal(&u)
		cached, _ := enc.Marshal(&u)
		gen, _ := u.MarshalJSON()
		easy, _ := easyjson.Marshal(EasyUser(u))
		for name, got := range map[string][]byte{"refljson": refl, "encoder": cached, "gen": gen, "easyjson": easy} {
			if string(got) != string(want) {
				t.Errorf("%s:\n got %s\nwant %s", name, got, want)
			}
		}
	}
}

func BenchmarkReflectNaive(b *testing.B) {
	u := testUser()
	for b.Loop() {
		_, _ = refljson.Marshal(&u)
	}
}

// Тот же reflect, но разбор типа закеширован — как в encoding/json.
func BenchmarkReflectCached(b *testing.B) {
	u := testUser()
	var enc refljson.Encoder
	for b.Loop() {
		_, _ = enc.Marshal(&u)
	}
}

func BenchmarkStdlibJSON(b *testing.B) {
	p := plainUser(testUser())
	for b.Loop() {
		_, _ = json.Marshal(&p)
	}
}

func BenchmarkGenerated(b *testing.B) {
	u := testUser()
	for b.Loop() {
		_, _ = u.MarshalJSON()
	}
}

func BenchmarkEasyJSON(b *testing.B) {
	e := EasyUser(testUser())
	for b.Loop() {
		_, _ = easyjson.Marshal(e)
	}
}

// encoding/json + сгенерированный метод: стандартная библиотека находит
// json.Marshaler и зовёт его, но потом ещё валидирует и компактирует вывод.
func BenchmarkStdlibCallsGenerated(b *testing.B) {
	u := testUser()
	for b.Loop() {
		_, _ = json.Marshal(&u)
	}
}
