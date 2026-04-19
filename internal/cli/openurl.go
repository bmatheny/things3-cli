package cli

import (
	"fmt"
	"os"
	"strings"
)

// redactAuthToken replaces auth-token values in a URL string with "***"
// to prevent secrets from appearing in logs or terminal output.
func redactAuthToken(s string) string {
	const prefix = "auth-token="
	i := strings.Index(s, prefix)
	if i < 0 {
		return s
	}
	start := i + len(prefix)
	end := strings.IndexByte(s[start:], '&')
	if end < 0 {
		return s[:start] + "***"
	}
	return s[:start] + "***" + s[start+end:]
}

func openURL(app *App, url string) error {
	if app == nil {
		return fmt.Errorf("Error: internal app not set")
	}
	if app.DryRun {
		fmt.Fprintln(app.Out, redactAuthToken(url))
		return nil
	}

	args := []string{url}
	if !app.Foreground {
		args = append([]string{"-g"}, args...)
	}

	if app.Debug {
		cmd := os.Getenv("OPEN")
		if cmd == "" {
			cmd = "open"
		}
		fmt.Fprintf(app.Err, "+ %s %s\n", cmd, redactAuthToken(strings.Join(args, " ")))
	}
	return app.Launcher.Open(args...)
}
