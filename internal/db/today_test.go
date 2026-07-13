package db

import (
	"database/sql"
	"reflect"
	"testing"
	"time"
)

func TestTodayTasks(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if err := seedTodayDB(conn); err != nil {
		t.Fatalf("seed db: %v", err)
	}

	store := &Store{conn: conn, path: ":memory:"}
	status := StatusIncomplete
	tasks, err := store.TodayTasks(TaskFilter{Status: &status, ExcludeTrashedContext: true})
	if err != nil {
		t.Fatalf("today tasks: %v", err)
	}

	titles := make(map[string]bool)
	buckets := make(map[string]int)
	for _, task := range tasks {
		titles[task.Title] = true
		if task.StartBucket != nil {
			buckets[task.Title] = *task.StartBucket
		}
	}

	for _, want := range []string{"Anytime Today", "This Evening", "Someday Past", "Overdue"} {
		if !titles[want] {
			t.Fatalf("expected task %q in results", want)
		}
	}
	for _, unwanted := range []string{"Suppressed Deadline", "Someday Future", "Completed Today"} {
		if titles[unwanted] {
			t.Fatalf("did not expect task %q in results", unwanted)
		}
	}
	if buckets["Anytime Today"] != 0 {
		t.Fatalf("expected Anytime Today start bucket 0, got %d", buckets["Anytime Today"])
	}
	if buckets["This Evening"] != 1 {
		t.Fatalf("expected This Evening start bucket 1, got %d", buckets["This Evening"])
	}
}

func TestTodayTasksCompositeOrderingAndLimit(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()
	if err := seedTodayDB(conn); err != nil {
		t.Fatalf("seed db: %v", err)
	}
	if _, err := conn.Exec(`DELETE FROM TMTask`); err != nil {
		t.Fatalf("clear tasks: %v", err)
	}

	today := thingsDate(time.Now())
	rows := []struct {
		uuid                     string
		bucket, reference, index *int
	}{
		{"NIL_BUCKET", nil, intPtr(999), intPtr(99)},
		{"NEW_REF", intPtr(0), intPtr(300), intPtr(50)},
		{"NULL_INDEX", intPtr(0), intPtr(200), nil},
		{"TIE_A", intPtr(0), intPtr(200), intPtr(1)},
		{"TIE_B", intPtr(0), intPtr(200), intPtr(1)},
		{"OLD_REF", intPtr(0), intPtr(100), intPtr(0)},
		{"NULL_REF", intPtr(0), nil, intPtr(0)},
		{"EVENING", intPtr(1), intPtr(999), intPtr(0)},
	}
	for _, row := range rows {
		if _, err := conn.Exec(`INSERT INTO TMTask
			(uuid, type, status, trashed, title, start, startDate, startBucket, todayIndexReferenceDate, todayIndex)
			VALUES (?, ?, ?, 0, ?, 1, ?, ?, ?, ?)`, row.uuid, TaskTypeTodo, StatusIncomplete, row.uuid, today, row.bucket, row.reference, row.index); err != nil {
			t.Fatalf("insert %s: %v", row.uuid, err)
		}
	}

	store := &Store{conn: conn, path: ":memory:"}
	status := StatusIncomplete
	tasks, err := store.TodayTasks(TaskFilter{Status: &status})
	if err != nil {
		t.Fatalf("today tasks: %v", err)
	}
	got := make([]string, len(tasks))
	for i, task := range tasks {
		got[i] = task.UUID
	}
	want := []string{"NIL_BUCKET", "NEW_REF", "NULL_INDEX", "TIE_A", "TIE_B", "OLD_REF", "NULL_REF", "EVENING"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected composite order:\n got %v\nwant %v", got, want)
	}

	if _, err := conn.Exec(`INSERT INTO TMTask
		(uuid, type, status, trashed, title, start, startDate, startBucket, todayIndexReferenceDate, todayIndex)
		VALUES ('SCHEDULED_LIMIT', ?, ?, 0, 'SCHEDULED_LIMIT', 2, ?, 0, 500, 0)`, TaskTypeTodo, StatusIncomplete, today); err != nil {
		t.Fatalf("insert scheduled limit candidate: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTask
		(uuid, type, status, trashed, title, start, deadline, startBucket, todayIndexReferenceDate, todayIndex)
		VALUES ('OVERDUE_LIMIT', ?, ?, 0, 'OVERDUE_LIMIT', 1, ?, 0, 400, 0)`, TaskTypeTodo, StatusIncomplete, today); err != nil {
		t.Fatalf("insert overdue limit candidate: %v", err)
	}

	limited, err := store.TodayTasks(TaskFilter{Status: &status, Limit: 1})
	if err != nil {
		t.Fatalf("limited today tasks: %v", err)
	}
	got = got[:0]
	for _, task := range limited {
		got = append(got, task.UUID)
	}
	want = []string{"NIL_BUCKET", "SCHEDULED_LIMIT", "OVERDUE_LIMIT"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected branch-limited candidates:\n got %v\nwant %v", got, want)
	}
}

func seedTodayDB(conn *sql.DB) error {
	statements := []string{
		`CREATE TABLE TMTask (
			uuid TEXT PRIMARY KEY,
			type INTEGER,
			status INTEGER,
			trashed INTEGER,
			title TEXT,
			notes TEXT,
			area TEXT,
			project TEXT,
			heading TEXT,
			start INTEGER,
			startDate INTEGER,
			startBucket INTEGER,
			todayIndexReferenceDate INTEGER,
			deadline INTEGER,
			deadlineSuppressionDate INTEGER,
			creationDate REAL,
			userModificationDate REAL,
			stopDate REAL,
			"index" INTEGER,
			todayIndex INTEGER,
			rt1_repeatingTemplate TEXT,
			rt1_recurrenceRule TEXT,
			repeater BLOB
		);`,
		`CREATE TABLE TMArea (uuid TEXT PRIMARY KEY, title TEXT, visible INTEGER, "index" INTEGER);`,
		`CREATE TABLE TMTag (uuid TEXT PRIMARY KEY, title TEXT, shortcut TEXT, parent TEXT);`,
		`CREATE TABLE TMTaskTag (tasks TEXT NOT NULL, tags TEXT NOT NULL);`,
	}
	for _, stmt := range statements {
		if _, err := conn.Exec(stmt); err != nil {
			return err
		}
	}

	today := thingsDate(time.Now())
	yesterday := thingsDate(time.Now().AddDate(0, 0, -1))
	future := thingsDate(time.Now().AddDate(0, 0, 2))
	todayBucket := 0
	eveningBucket := 1

	inserts := []struct {
		uuid        string
		title       string
		start       int
		startDate   *int
		startBucket *int
		deadline    *int
		deadSuppr   *int
		status      int
		trashed     int
		index       *int
		recurrence  *string
	}{
		{"T1", "Anytime Today", 1, &today, &todayBucket, nil, nil, StatusIncomplete, 0, nil, nil},
		{"T2", "Someday Past", 2, &yesterday, nil, nil, nil, StatusIncomplete, 0, nil, nil},
		{"T3", "Overdue", 1, nil, nil, &yesterday, nil, StatusIncomplete, 0, nil, nil},
		{"T4", "Suppressed Deadline", 1, nil, nil, &yesterday, intPtr(1), StatusIncomplete, 0, nil, nil},
		{"T5", "Someday Future", 2, &future, nil, nil, nil, StatusIncomplete, 0, nil, nil},
		{"T6", "Completed Today", 1, &today, nil, nil, nil, StatusCompleted, 0, nil, nil},
		{"T7", "This Evening", 1, &today, &eveningBucket, nil, nil, StatusIncomplete, 0, nil, nil},
	}

	for _, item := range inserts {
		_, err := conn.Exec(
			`INSERT INTO TMTask (uuid, type, status, trashed, title, start, startDate, startBucket, deadline, deadlineSuppressionDate, todayIndex, rt1_recurrenceRule)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			item.uuid,
			TaskTypeTodo,
			item.status,
			item.trashed,
			item.title,
			item.start,
			item.startDate,
			item.startBucket,
			item.deadline,
			item.deadSuppr,
			item.index,
			item.recurrence,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func thingsDate(t time.Time) int {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return date.Year()<<16 | int(date.Month())<<12 | date.Day()<<7
}

func intPtr(v int) *int {
	return &v
}
