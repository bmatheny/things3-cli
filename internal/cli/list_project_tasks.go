package cli

import (
	"fmt"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

// NewListProjectTasksCommand builds the list-project-tasks subcommand.
func NewListProjectTasksCommand(app *App) *cobra.Command {
	var dbPath string
	var id string
	opts := TaskQueryOptions{
		Status: "incomplete",
		Limit:  200,
	}
	var format string
	var selectRaw string
	var asJSON bool
	var noHeader bool

	cmd := &cobra.Command{
		Use:   "list-project-tasks --id <UUID>",
		Short: "List todos belonging to a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("Error: --id is required")
			}

			store, _, err := db.OpenDefault(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			opts.Project = id
			opts.HasURLSet = cmd.Flags().Changed("has-url")

			outputOpts, err := resolveTaskOutputOptions(format, asJSON, selectRaw, noHeader)
			if err != nil {
				return err
			}
			tasks, err := fetchTasks(store, store.Tasks, opts, false, []int{db.TaskTypeTodo})
			if err != nil {
				return formatDBError(err)
			}
			return printTasks(app.Out, tasks, outputOpts)
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "UUID of the project whose tasks to list (required)")
	cmd.Flags().StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	cmd.Flags().StringVar(&dbPath, "database", "", "Alias for --db")
	addTaskQueryFlags(cmd, &opts, true, true)
	addTaskOutputFlags(cmd, &format, &selectRaw, &asJSON, &noHeader)

	return cmd
}
