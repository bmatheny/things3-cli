package things

import "strings"

// RenameProjectOptions defines options for rename-project.
type RenameProjectOptions struct {
	AuthToken string
	ID        string
	Title     string
}

// BuildRenameProjectURL builds a Things URL for the rename-project command.
func BuildRenameProjectURL(opts RenameProjectOptions) (string, error) {
	if opts.AuthToken == "" {
		return "", ErrMissingAuthToken
	}
	if opts.ID == "" {
		return "", errMissingID
	}
	if opts.Title == "" {
		return "", errMissingTitle
	}

	params := []string{
		"auth-token=" + URLEncode(opts.AuthToken),
		"id=" + URLEncode(opts.ID),
		"title=" + URLEncode(opts.Title),
	}

	return "things:///update-project?" + strings.Join(params, "&") + "&", nil
}
