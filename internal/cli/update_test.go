package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestUpdateCommandRequiresAuthToken(t *testing.T) {
	t.Setenv("THINGS_AUTH_TOKEN", "")
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--id", "123", "Title"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandWithAuthAndID(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--no-verify", "Title"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "auth-token=tok") {
		t.Fatalf("expected auth-token in url, got %q", url)
	}
	if !strings.Contains(url, "id=123") {
		t.Fatalf("expected id in url, got %q", url)
	}
}

func TestUpdateCommandLaterFlag(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--later", "--no-verify"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "when=evening") {
		t.Fatalf("expected when=evening in url, got %q", url)
	}
}

func TestUpdateCommandEmptyValuesClearFields(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--when=", "--deadline=", "--tags=", "--no-verify"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	for _, param := range []string{"when=&", "deadline=&", "tags=&"} {
		if !strings.Contains(url, param) {
			t.Fatalf("expected %q in url, got %q", param, url)
		}
	}
}

func TestUpdateCommandRejectsUnsafeTitle(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--no-verify", "tag=work"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandAllowsUnsafeTitleWithFlag(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--allow-unsafe-title", "--no-verify", "tag=work"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "title=tag%3Dwork") {
		t.Fatalf("expected title in url, got %q", url)
	}
}

func TestUpdateCommandBlocksEveningForNonToday(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "UP1", "--when=evening", "--no-verify"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandAllowsEveningForNonTodayWithFlag(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "UP1", "--when=evening", "--allow-non-today", "--no-verify"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if len(launcher.args) == 0 {
		t.Fatalf("expected open invocation")
	}
}

func TestUpdateCommandCompletesChecklistItemViaJSON(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "T1", "--complete-checklist-item", "Check Item"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	url := requireOpenURL(t, launcher)
	if !strings.HasPrefix(url, "things:///json?") {
		t.Fatalf("expected json url, got %q", url)
	}
	if !strings.Contains(url, "auth-token=tok") {
		t.Fatalf("expected auth token in url, got %q", url)
	}
	if !strings.Contains(url, "completed%22%3Atrue") {
		t.Fatalf("expected completed checklist payload, got %q", url)
	}
}

func TestUpdateCommandChecklistStatusRequiresID(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--complete-checklist-item", "Check Item"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandChecklistStatusRejectsOtherChanges(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "T1", "--complete-checklist-item", "Check Item", "--notes", "Nope"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}
