package cli

import (
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewRenameProjectCommand builds the rename-project subcommand.
func NewRenameProjectCommand(app *App) *cobra.Command {
	opts := things.RenameProjectOptions{}

	cmd := &cobra.Command{
		Use:   "rename-project --id <UUID> --title <TITLE>",
		Short: "Rename an existing project",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := resolveAuthToken(app, opts.AuthToken)
			if err != nil {
				return err
			}
			opts.AuthToken = token

			url, err := things.BuildRenameProjectURL(opts)
			if err != nil {
				return err
			}
			return openURL(app, url)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.AuthToken, "auth-token", "", "Things URL scheme authorization token")
	flags.StringVar(&opts.ID, "id", "", "ID of the project to rename")
	flags.StringVar(&opts.Title, "title", "", "New title for the project")

	return cmd
}
