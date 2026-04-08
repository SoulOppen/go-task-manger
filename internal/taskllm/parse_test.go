package taskllm

import (
	"strings"
	"testing"
)

func TestParseTasksJSON_unaTarea(t *testing.T) {
	raw := `{"tasks":[{"name":"A","description":"B","relevance":7,"due":"","depends_on_id":""}]}`
	items, err := ParseTasksJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Name != "A" {
		t.Fatalf("%+v", items)
	}
}

func TestParseTasksJSON_variasTareas(t *testing.T) {
	raw := `{"tasks":[{"name":"A","description":"d1","relevance":5,"due":"","depends_on_id":""},{"name":"B","description":"d2","relevance":3,"due":"2026-12-01","depends_on_id":""}]}`
	items, err := ParseTasksJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[1].Due != "2026-12-01" {
		t.Fatalf("%+v", items)
	}
}

func TestParseTasksJSON_claveExtra(t *testing.T) {
	raw := `{"tasks":[{"name":"A","description":"B","relevance":5,"due":"","depends_on_id":"","extra":1}]}`
	_, err := ParseTasksJSON(raw)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseTasksJSON_conFence(t *testing.T) {
	raw := "```json\n{\"tasks\":[{\"name\":\"X\",\"description\":\"Y\",\"relevance\":5,\"due\":\"\",\"depends_on_id\":\"\"}]}\n```"
	items, err := ParseTasksJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Name != "X" {
		t.Fatal(items)
	}
}

func TestBuildTasks_relevanceFueraRango(t *testing.T) {
	items := []taskItem{{Name: "n", Description: "d", Relevance: 99, Due: "", DependsOnID: ""}}
	tasks, err := BuildTasks(items)
	if err != nil {
		t.Fatal(err)
	}
	if tasks[0].Relevance != 5 {
		t.Fatalf("got %d", tasks[0].Relevance)
	}
}

func TestBuildTasks_invalidDue(t *testing.T) {
	items := []taskItem{{Name: "n", Description: "d", Relevance: 5, Due: "no-fecha", DependsOnID: ""}}
	_, err := BuildTasks(items)
	if err == nil || !strings.Contains(err.Error(), "tarea") {
		t.Fatalf("expected due error, got %v", err)
	}
}
