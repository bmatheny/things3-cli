package things

import (
	"encoding/json"
	"fmt"
	"strings"
)

// UpdateOptions defines options for update.
type UpdateOptions struct {
	AuthToken                string
	ID                       string
	Notes                    string
	NotesSet                 bool
	PrependNotes             string
	AppendNotes              string
	When                     string
	WhenSet                  bool
	Later                    bool
	Deadline                 string
	DeadlineSet              bool
	Tags                     string
	TagsSet                  bool
	AddTags                  string
	Completed                bool
	Canceled                 bool
	Reveal                   bool
	Duplicate                bool
	CompletionDate           string
	CreationDate             string
	Heading                  string
	List                     string
	ListID                   string
	ChecklistItems           []string
	PrependChecklistItems    []string
	AppendChecklistItems     []string
	CompleteChecklistItems   []string
	IncompleteChecklistItems []string
}

// ChecklistItemState describes the checklist fields needed for JSON updates.
type ChecklistItemState struct {
	Title  string
	Status int
}

// BuildUpdateURL builds a Things URL for the update command.
func BuildUpdateURL(opts UpdateOptions, rawInput string) (string, error) {
	if opts.AuthToken == "" {
		return "", ErrMissingAuthToken
	}
	if opts.ID == "" {
		return "", errMissingID
	}

	var title string
	notes := opts.Notes
	notesSet := opts.NotesSet

	if rawInput != "" {
		if HasMultipleLines(rawInput) {
			title = FindTitle(rawInput)
			notes = FindNotes(rawInput)
			notesSet = notes != ""
		} else {
			title = rawInput
		}
	}

	params := make([]string, 0, 24)
	params = append(params, "auth-token="+URLEncode(opts.AuthToken))
	params = append(params, "id="+URLEncode(opts.ID))

	if title != "" {
		params = append(params, "title="+URLEncode(title))
	}

	if len(opts.ChecklistItems) > 0 {
		encoded := make([]string, 0, len(opts.ChecklistItems))
		for _, item := range opts.ChecklistItems {
			encoded = append(encoded, URLEncode(item))
		}
		params = append(params, "checklist-items="+Join(encoded...))
	}

	if len(opts.PrependChecklistItems) > 0 {
		encoded := make([]string, 0, len(opts.PrependChecklistItems))
		for _, item := range opts.PrependChecklistItems {
			encoded = append(encoded, URLEncode(item))
		}
		params = append(params, "prepend-checklist-items="+Join(encoded...))
	}

	if len(opts.AppendChecklistItems) > 0 {
		encoded := make([]string, 0, len(opts.AppendChecklistItems))
		for _, item := range opts.AppendChecklistItems {
			encoded = append(encoded, URLEncode(item))
		}
		params = append(params, "append-checklist-items="+Join(encoded...))
	}

	if opts.PrependNotes != "" {
		params = append(params, "prepend-notes="+URLEncode(opts.PrependNotes))
	}

	if opts.AppendNotes != "" {
		params = append(params, "append-notes="+URLEncode(opts.AppendNotes))
	}

	if opts.Heading != "" {
		params = append(params, "heading="+URLEncode(opts.Heading))
	}

	if opts.Duplicate {
		params = append(params, "duplicate=true")
	}

	if opts.When != "" {
		params = append(params, "when="+URLEncode(opts.When))
	} else if opts.Later {
		params = append(params, "when=evening")
	} else if opts.WhenSet {
		params = append(params, "when=")
	}

	if opts.Deadline != "" || opts.DeadlineSet {
		params = append(params, "deadline="+URLEncode(opts.Deadline))
	}

	if opts.Reveal {
		params = append(params, "reveal=true")
	}

	if opts.Tags != "" || opts.TagsSet {
		params = append(params, "tags="+URLEncode(opts.Tags))
	}

	if opts.AddTags != "" {
		params = append(params, "add-tags="+URLEncode(opts.AddTags))
	}

	if notes != "" || notesSet {
		params = append(params, "notes="+URLEncode(notes))
	}

	if opts.CreationDate != "" {
		params = append(params, "creation-date="+URLEncode(opts.CreationDate))
	}

	if opts.CompletionDate != "" {
		params = append(params, "completion-date="+URLEncode(opts.CompletionDate))
	}

	if opts.Canceled {
		params = append(params, "canceled=true")
	} else if opts.Completed {
		params = append(params, "completed=true")
	}

	if opts.ListID != "" {
		params = append(params, "list-id="+URLEncode(opts.ListID))
	} else if opts.List != "" {
		params = append(params, "list="+URLEncode(opts.List))
	}

	return "things:///update?" + strings.Join(params, "&") + "&", nil
}

// BuildChecklistStatusUpdateURL builds a Things JSON URL that replaces the
// current checklist with the same titles and updated completion statuses.
func BuildChecklistStatusUpdateURL(opts UpdateOptions, checklist []ChecklistItemState) (string, error) {
	if opts.AuthToken == "" {
		return "", ErrMissingAuthToken
	}
	if opts.ID == "" {
		return "", errMissingID
	}
	updated, err := applyChecklistStatusChanges(checklist, opts.CompleteChecklistItems, opts.IncompleteChecklistItems)
	if err != nil {
		return "", err
	}

	operation := []map[string]any{
		{
			"type":      "to-do",
			"operation": "update",
			"id":        opts.ID,
			"attributes": map[string]any{
				"checklist-items": updated,
			},
		},
	}
	data, err := json.Marshal(operation)
	if err != nil {
		return "", err
	}
	return "things:///json?auth-token=" + URLEncode(opts.AuthToken) + "&data=" + URLEncode(string(data)) + "&", nil
}

func applyChecklistStatusChanges(checklist []ChecklistItemState, completeTitles, incompleteTitles []string) ([]map[string]any, error) {
	if len(checklist) == 0 {
		return nil, fmt.Errorf("Error: todo has no checklist items")
	}
	requested := make(map[string]bool, len(completeTitles)+len(incompleteTitles))
	for _, title := range completeTitles {
		if requested[title] {
			return nil, fmt.Errorf("Error: checklist item %q was requested more than once", title)
		}
		requested[title] = true
	}
	for _, title := range incompleteTitles {
		if requested[title] {
			return nil, fmt.Errorf("Error: checklist item %q cannot be both complete and incomplete", title)
		}
		requested[title] = true
	}

	counts := make(map[string]int, len(checklist))
	for _, item := range checklist {
		counts[item.Title]++
	}
	for title := range requested {
		switch counts[title] {
		case 0:
			return nil, fmt.Errorf("Error: checklist item %q not found", title)
		case 1:
		default:
			return nil, fmt.Errorf("Error: checklist item %q is ambiguous (%d matches)", title, counts[title])
		}
	}

	complete := make(map[string]bool, len(completeTitles))
	for _, title := range completeTitles {
		complete[title] = true
	}
	incomplete := make(map[string]bool, len(incompleteTitles))
	for _, title := range incompleteTitles {
		incomplete[title] = true
	}

	updated := make([]map[string]any, 0, len(checklist))
	for _, item := range checklist {
		attrs := map[string]any{"title": item.Title}
		switch {
		case complete[item.Title]:
			attrs["completed"] = true
		case incomplete[item.Title]:
		case item.Status == 3:
			attrs["completed"] = true
		case item.Status == 2:
			attrs["canceled"] = true
		}
		updated = append(updated, map[string]any{
			"type":       "checklist-item",
			"attributes": attrs,
		})
	}
	return updated, nil
}
