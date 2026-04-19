package integration_test

import "testing"

func TestShowRequiresTarget(t *testing.T) {
	dbPath := writeTestDB(t)
	_, errOut, code := runThings(t, "", "show", "--db", dbPath)
	requireFailure(t, code)
	assertContains(t, errOut, "Must specify --id=ID or query")
}

func TestShowAcceptsID(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "show", "--db", dbPath, "--id=A1")
	requireSuccess(t, code)
	assertContains(t, out, "Home")
}

func TestShowAcceptsQuery(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "show", "--db", dbPath, "Project One")
	requireSuccess(t, code)
	assertContains(t, out, "Project One")
}

func TestShowAcceptsQueryFromStdin(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "Project One", "show", "--db", dbPath, "-")
	requireSuccess(t, code)
	assertContains(t, out, "Project One")
}

func TestShowRecursiveIncludesChecklist(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "show", "--db", dbPath, "--id=T1", "--recursive")
	requireSuccess(t, code)
	assertContains(t, out, "CHECKLIST")
	assertContains(t, out, "Check Item")
}

func TestShowRecursiveIncludesChecklistJSON(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "show", "--db", dbPath, "--id=T1", "--recursive", "--json")
	requireSuccess(t, code)
	assertContains(t, out, `"checklist":[`)
	assertContains(t, out, `"title":"Check Item"`)
}
