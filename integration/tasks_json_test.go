package integration_test

import "testing"

func TestTasksJSONIncludesDetails(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "tasks", "--db", dbPath, "--json", "--recursive")
	requireSuccess(t, code)
	assertContains(t, out, `"notes":"Some notes"`)
	assertContains(t, out, `"tags":["urgent"]`)
	assertContains(t, out, `"checklist":[`)
	assertContains(t, out, `"title":"Check Item"`)
}

func TestTasksJSONTodayIndexReferenceDate(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "tasks", "--db", dbPath, "--format", "json", "--select", "uuid,today-index-reference-date")
	requireSuccess(t, code)
	assertContains(t, out, `"today_index_reference_date":135004288`)
	assertContains(t, out, `"today_index_reference_date":null`)
}
