package things

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

func TestBuildUpdateURLErrorMissingAuthToken(t *testing.T) {
	_, err := BuildUpdateURL(UpdateOptions{ID: "123"}, "")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Error: Missing Things auth token. Run `things auth` for setup, set THINGS_AUTH_TOKEN, or pass --auth-token=TOKEN (Things > Settings > General > Things URLs)." {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildUpdateURLErrorMissingID(t *testing.T) {
	_, err := BuildUpdateURL(UpdateOptions{AuthToken: "tok"}, "")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Error: Must specify --id=id" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildUpdateURLCanceledOverridesCompleted(t *testing.T) {
	url, err := BuildUpdateURL(UpdateOptions{AuthToken: "tok", ID: "id", Completed: true, Canceled: true}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(url, "canceled=true") {
		t.Fatalf("expected canceled in %q", url)
	}
	if contains(url, "completed=true") {
		t.Fatalf("did not expect completed in %q", url)
	}
}

func TestBuildUpdateURLListIDPrecedence(t *testing.T) {
	url, err := BuildUpdateURL(UpdateOptions{AuthToken: "tok", ID: "id", List: "Work", ListID: "123"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(url, "list-id=123") {
		t.Fatalf("expected list-id in %q", url)
	}
	if contains(url, "list=Work") {
		t.Fatalf("did not expect list in %q", url)
	}
}

func TestBuildUpdateURLNotesFromInputOverrideFlag(t *testing.T) {
	opts := UpdateOptions{AuthToken: "tok", ID: "id", Notes: "FromFlag"}
	url, err := BuildUpdateURL(opts, "Title\n\nFromInput")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(url, "title=Title") {
		t.Fatalf("expected title in %q", url)
	}
	if !contains(url, "notes=FromInput") {
		t.Fatalf("expected notes from input in %q", url)
	}
	if contains(url, "notes=FromFlag") {
		t.Fatalf("did not expect notes from flag in %q", url)
	}
}

func TestBuildUpdateURLChecklistJoin(t *testing.T) {
	opts := UpdateOptions{AuthToken: "tok", ID: "id", ChecklistItems: []string{"One", "Two"}}
	url, err := BuildUpdateURL(opts, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(url, "checklist-items=One%0ATwo") {
		t.Fatalf("expected checklist-items in %q", url)
	}
}

func TestBuildUpdateURLLaterSetsEvening(t *testing.T) {
	opts := UpdateOptions{AuthToken: "tok", ID: "id", Later: true}
	url, err := BuildUpdateURL(opts, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(url, "when=evening") {
		t.Fatalf("expected when=evening in %q", url)
	}
}

func TestBuildUpdateURLWhenOverridesLater(t *testing.T) {
	opts := UpdateOptions{AuthToken: "tok", ID: "id", Later: true, When: "today"}
	url, err := BuildUpdateURL(opts, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(url, "when=today") {
		t.Fatalf("expected when=today in %q", url)
	}
	if contains(url, "when=evening") {
		t.Fatalf("did not expect when=evening in %q", url)
	}
}

func TestBuildUpdateURLTrailingAmpersand(t *testing.T) {
	url, err := BuildUpdateURL(UpdateOptions{AuthToken: "tok", ID: "id", Notes: "Notes"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(url, "&") {
		t.Fatalf("expected trailing ampersand in %q", url)
	}
}

func TestBuildChecklistStatusUpdateURL(t *testing.T) {
	built, err := BuildChecklistStatusUpdateURL(UpdateOptions{
		AuthToken:                "tok",
		ID:                       "id",
		CompleteChecklistItems:   []string{"One"},
		IncompleteChecklistItems: []string{"Two"},
	}, []ChecklistItemState{
		{Title: "One", Status: 0},
		{Title: "Two", Status: 3},
		{Title: "Three", Status: 3},
		{Title: "Four", Status: 2},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(built, "things:///json?auth-token=tok&data=") {
		t.Fatalf("expected json url, got %q", built)
	}
	encoded := strings.TrimSuffix(strings.TrimPrefix(built, "things:///json?auth-token=tok&data="), "&")
	data, err := url.QueryUnescape(encoded)
	if err != nil {
		t.Fatalf("unescape data: %v", err)
	}
	var operations []map[string]any
	if err := json.Unmarshal([]byte(data), &operations); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if got := operations[0]["operation"]; got != "update" {
		t.Fatalf("expected update operation, got %#v", got)
	}
	attrs := operations[0]["attributes"].(map[string]any)
	items := attrs["checklist-items"].([]any)
	first := items[0].(map[string]any)["attributes"].(map[string]any)
	if first["completed"] != true {
		t.Fatalf("expected first item completed, got %#v", first)
	}
	second := items[1].(map[string]any)["attributes"].(map[string]any)
	if _, ok := second["completed"]; ok {
		t.Fatalf("expected second item incomplete, got %#v", second)
	}
	third := items[2].(map[string]any)["attributes"].(map[string]any)
	if third["completed"] != true {
		t.Fatalf("expected third item to preserve completed, got %#v", third)
	}
	fourth := items[3].(map[string]any)["attributes"].(map[string]any)
	if fourth["canceled"] != true {
		t.Fatalf("expected fourth item to preserve canceled, got %#v", fourth)
	}
}

func TestBuildChecklistStatusUpdateURLRejectsMissingChecklistItem(t *testing.T) {
	_, err := BuildChecklistStatusUpdateURL(UpdateOptions{
		AuthToken:              "tok",
		ID:                     "id",
		CompleteChecklistItems: []string{"Missing"},
	}, []ChecklistItemState{{Title: "One", Status: 0}})
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != `Error: checklist item "Missing" not found` {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildChecklistStatusUpdateURLRejectsAmbiguousChecklistItem(t *testing.T) {
	_, err := BuildChecklistStatusUpdateURL(UpdateOptions{
		AuthToken:              "tok",
		ID:                     "id",
		CompleteChecklistItems: []string{"One"},
	}, []ChecklistItemState{{Title: "One", Status: 0}, {Title: "One", Status: 3}})
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != `Error: checklist item "One" is ambiguous (2 matches)` {
		t.Fatalf("unexpected error: %v", err)
	}
}
