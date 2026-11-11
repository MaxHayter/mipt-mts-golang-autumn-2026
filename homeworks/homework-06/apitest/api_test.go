//go:build integration

package apitest

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type snapshot struct {
	ID       string `json:"id"`
	Version  int    `json:"version"`
	Status   string `json:"status"`
	Comments []struct {
		Author    string `json:"author"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
	} `json:"comments"`
}

func baseURL() string {
	if s := os.Getenv("INCIDENT_API_URL"); s != "" {
		return strings.TrimRight(s, "/")
	}
	return "http://localhost:8080"
}

func call(method, path string, body any) (int, []byte, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		reader = bytes.NewReader(b)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, baseURL()+path, reader)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "hw06-check")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("запустите оба приложения: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return 0, nil, err
	}
	if resp.Header.Get("X-Request-ID") != "hw06-check" {
		return 0, nil, fmt.Errorf("X-Request-ID не сохранён")
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return 0, nil, fmt.Errorf("ожидался application/json, получен %q", resp.Header.Get("Content-Type"))
	}
	return resp.StatusCode, b, nil
}

func request(t *testing.T, method, path string, body any, want int) []byte {
	t.Helper()
	code, b, err := call(method, path, body)
	if err != nil {
		t.Fatal(err)
	}
	if code != want {
		t.Fatalf("%s %s: code=%d, want=%d, body=%s", method, path, code, want, b)
	}
	if code >= 400 {
		var e struct {
			Error     string `json:"error"`
			RequestID string `json:"request_id"`
		}
		if err := json.Unmarshal(b, &e); err != nil {
			t.Fatal(err)
		}
		if e.Error == "" || e.RequestID != "hw06-check" {
			t.Fatalf("неверный формат ошибки: %s", b)
		}
	}
	return b
}

func subset(t *testing.T, got, want any, path string) {
	t.Helper()
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			t.Fatalf("%s: ожидался объект, получено %v", path, got)
		}
		for k, v := range w {
			subset(t, g[k], v, path+"."+k)
		}
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			t.Fatalf("%s: неверный массив %v", path, got)
		}
		for i, v := range w {
			subset(t, g[i], v, fmt.Sprintf("%s[%d]", path, i))
		}
	default:
		if got != want {
			t.Fatalf("%s: got=%v, want=%v", path, got, want)
		}
	}
}

func decode(t *testing.T, b []byte) snapshot {
	t.Helper()
	var s snapshot
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatal(err)
	}
	if s.ID == "" || s.Version < 1 || s.Comments == nil {
		t.Fatalf("неверный снимок: %s", b)
	}
	if s.Status != "OPEN" && s.Status != "INVESTIGATING" && s.Status != "RESOLVED" {
		t.Fatalf("неверный статус: %s", b)
	}
	for _, c := range s.Comments {
		if _, err := time.Parse(time.RFC3339Nano, c.CreatedAt); err != nil {
			t.Fatalf("created_at: %v", err)
		}
	}
	return s
}

func create(t *testing.T) string {
	t.Helper()
	return decode(t, request(t, "POST", "/incidents", map[string]string{"title": "Сбой оплаты", "service": "payments"}, 201)).ID
}

func TestScenarios(t *testing.T) {
	b, err := os.ReadFile("testdata/scenarios.json")
	if err != nil {
		t.Fatal(err)
	}
	var scenarios []struct {
		Name  string `json:"name"`
		Steps []struct {
			Method   string `json:"method"`
			Path     string `json:"path"`
			Body     any    `json:"body"`
			Code     int    `json:"code"`
			Expected any    `json:"expected"`
			Count    *int   `json:"comments_count"`
		} `json:"steps"`
	}
	if err := json.Unmarshal(b, &scenarios); err != nil {
		t.Fatal(err)
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			id := ""
			for i, s := range scenario.Steps {
				path := strings.ReplaceAll(s.Path, "{id}", id)
				b := request(t, s.Method, path, s.Body, s.Code)
				if s.Code < 300 {
					state := decode(t, b)
					if s.Method == "POST" && s.Path == "/incidents" {
						id = state.ID
					} else if state.ID != id {
						t.Fatalf("step %d: изменился ID", i)
					}
					if s.Count != nil && len(state.Comments) != *s.Count {
						t.Fatalf("step %d: comments=%d, want=%d", i, len(state.Comments), *s.Count)
					}
					if s.Expected != nil {
						var got any
						if err := json.Unmarshal(b, &got); err != nil {
							t.Fatal(err)
						}
						subset(t, got, s.Expected, fmt.Sprintf("step %d", i))
					}
				}
			}
		})
	}
}

func TestMalformedJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL()+"/incidents", strings.NewReader(`{"title":`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("got=%d, want=400", resp.StatusCode)
	}
}

func TestConcurrentComments(t *testing.T) {
	id := create(t)
	const n = 20
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code, b, err := call("POST", "/incidents/"+id+"/comments", map[string]string{"author": "Nick", "text": fmt.Sprintf("comment-%d", i)})
			if err != nil {
				errs <- err
			} else if code != 200 {
				errs <- fmt.Errorf("code=%d body=%s", code, b)
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	if t.Failed() {
		return
	}
	s := decode(t, request(t, "GET", "/incidents/"+id, nil, 200))
	if len(s.Comments) != n || s.Version != n+1 {
		t.Fatalf("comments=%d version=%d", len(s.Comments), s.Version)
	}
	seen := make(map[string]bool)
	for _, c := range s.Comments {
		seen[c.Text] = true
	}
	for i := 0; i < n; i++ {
		if !seen[fmt.Sprintf("comment-%d", i)] {
			t.Errorf("потерян comment-%d", i)
		}
	}
}

type event struct {
	data []byte
	err  error
}

func watch(t *testing.T, id string) (<-chan event, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL()+"/incidents/"+id+"/events", nil)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	if resp.StatusCode != 200 || !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/event-stream") {
		resp.Body.Close()
		cancel()
		t.Fatalf("неверный SSE ответ: %s", resp.Status)
	}
	out := make(chan event, 32)
	go func() {
		defer close(out)
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 4096), 2<<20)
		kind := ""
		var data []string
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				if len(data) > 0 {
					e := event{data: []byte(strings.Join(data, "\n"))}
					if kind != "snapshot" {
						e.err = fmt.Errorf("неожиданное событие %q", kind)
					}
					select {
					case out <- e:
					case <-ctx.Done():
						return
					}
				}
				kind = ""
				data = nil
				continue
			}
			if strings.HasPrefix(line, "event:") {
				kind = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			}
			if strings.HasPrefix(line, "data:") {
				data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
			}
		}
		if err := scanner.Err(); err != nil {
			select {
			case out <- event{err: err}:
			case <-ctx.Done():
			}
		}
	}()
	t.Cleanup(cancel)
	return out, cancel
}
func next(t *testing.T, ch <-chan event) snapshot {
	t.Helper()
	select {
	case e, ok := <-ch:
		if !ok {
			t.Fatal("SSE завершился до снимка")
		}
		if e.err != nil {
			t.Fatal(e.err)
		}
		return decode(t, e.data)
	case <-time.After(5 * time.Second):
		t.Fatal("SSE: снимок не пришёл за 5 секунд")
	}
	return snapshot{}
}
func ended(t *testing.T, ch <-chan event) {
	t.Helper()
	select {
	case e, ok := <-ch:
		if ok {
			t.Fatalf("ожидалось завершение SSE, получено %s / %v", e.data, e.err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("SSE не завершился")
	}
}
func TestSubscriptions(t *testing.T) {
	id := create(t)
	a, cancelA := watch(t, id)
	b, cancelB := watch(t, id)
	defer cancelB()
	for _, ch := range []<-chan event{a, b} {
		s := next(t, ch)
		if s.ID != id || s.Version != 1 {
			t.Fatalf("неверный первый снимок: %+v", s)
		}
	}
	request(t, "PUT", "/incidents/"+id+"/assignee", map[string]string{"user": "Nick"}, 200)
	for _, ch := range []<-chan event{a, b} {
		if s := next(t, ch); s.Version != 2 {
			t.Fatalf("version=%d, want=2", s.Version)
		}
	}
	cancelA()
	request(t, "PUT", "/incidents/"+id+"/status", map[string]string{"status": "INVESTIGATING"}, 200)
	if s := next(t, b); s.Version != 3 || s.Status != "INVESTIGATING" {
		t.Fatalf("неверное обновление: %+v", s)
	}
	request(t, "PUT", "/incidents/"+id+"/status", map[string]string{"status": "RESOLVED"}, 200)
	if s := next(t, b); s.Version != 4 || s.Status != "RESOLVED" {
		t.Fatalf("неверный финальный снимок: %+v", s)
	}
	ended(t, b)
	c, cancelC := watch(t, id)
	defer cancelC()
	if s := next(t, c); s.Version != 4 || s.Status != "RESOLVED" {
		t.Fatalf("неверный закрытый инцидент: %+v", s)
	}
	ended(t, c)
}
