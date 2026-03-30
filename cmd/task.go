package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/SoulOppen/task-manager-go/internal/db"
	"github.com/SoulOppen/task-manager-go/internal/task"
	"github.com/spf13/cobra"
)

// withTaskRepo abre MySQL, migra (tasks + users) y ejecuta fn.
func withTaskRepo(ctx context.Context, fn func(*task.Repository) error) error {
	return db.WithDB(ctx, func(d *sql.DB) error {
		return fn(task.NewRepository(d))
	})
}

var (
	addName        string
	addDescription string
	addRelevance   int
	addDue         string
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "CRUD de tareas (MySQL)",
}

var taskAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Crea una tarea",
	RunE: func(cmd *cobra.Command, args []string) error {
		due, err := task.ParseDueDate(addDue)
		if err != nil {
			return err
		}
		t := task.NewTask(addName, addDescription, addRelevance, due)
		if err := t.Validate(); err != nil {
			return err
		}
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			if err := repo.Create(cmd.Context(), t); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), t.ID)
			return nil
		})
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista tareas por relevancia y fecha de entrega",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			tasks, err := repo.ListOrdered(cmd.Context())
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNOMBRE\tREL\tENTREGA\tDESCRIPCION\tCREADO")
			for _, t := range tasks {
				due := "-"
				if t.DueDate != nil {
					due = t.DueDate.Format(task.DateLayout)
				}
				desc := strings.ReplaceAll(t.Description, "\n", " ")
				if len(desc) > 40 {
					desc = desc[:37] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\n", t.ID, t.Name, t.Relevance, due, desc, t.CreatedAt.UTC().Format(time.RFC3339))
			}
			return w.Flush()
		})
	},
}

var taskGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Muestra una tarea por id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			t, err := repo.GetByID(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			due := "-"
			if t.DueDate != nil {
				due = t.DueDate.Format(task.DateLayout)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "id: %s\nnombre: %s\ndescripcion: %s\nrelevancia: %d\nentrega: %s\ncreado: %s\n",
				t.ID, t.Name, t.Description, t.Relevance, due, t.CreatedAt.UTC().Format(time.RFC3339))
			return nil
		})
	},
}

var (
	updName         string
	updDescription  string
	updRelevance      int
	updDue            string
	updClearDue       bool
)

var taskUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Actualiza una tarea",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			existing, err := repo.GetByID(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("name") {
				existing.Name = strings.TrimSpace(updName)
			}
			if cmd.Flags().Changed("description") {
				existing.Description = strings.TrimSpace(updDescription)
			}
			if cmd.Flags().Changed("relevance") {
				existing.Relevance = updRelevance
			}
			if updClearDue {
				existing.DueDate = nil
			} else if cmd.Flags().Changed("due") {
				due, err := task.ParseDueDate(updDue)
				if err != nil {
					return err
				}
				existing.DueDate = due
			}
			if err := existing.Validate(); err != nil {
				return err
			}
			return repo.Update(cmd.Context(), existing)
		})
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Elimina una tarea",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			return repo.Delete(cmd.Context(), args[0])
		})
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskAddCmd)
	taskAddCmd.Flags().StringVar(&addName, "name", "", "nombre de la tarea")
	taskAddCmd.Flags().StringVar(&addDescription, "description", "", "descripcion")
	taskAddCmd.Flags().IntVar(&addRelevance, "relevance", 5, "relevancia 1-10")
	taskAddCmd.Flags().StringVar(&addDue, "due", "", "fecha de entrega YYYY-MM-DD (opcional)")
	_ = taskAddCmd.MarkFlagRequired("name")
	_ = taskAddCmd.MarkFlagRequired("description")

	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskGetCmd)
	taskCmd.AddCommand(taskUpdateCmd)
	taskUpdateCmd.Flags().StringVar(&updName, "name", "", "nuevo nombre")
	taskUpdateCmd.Flags().StringVar(&updDescription, "description", "", "nueva descripcion")
	taskUpdateCmd.Flags().IntVar(&updRelevance, "relevance", 0, "nueva relevancia 1-10")
	taskUpdateCmd.Flags().StringVar(&updDue, "due", "", "nueva fecha YYYY-MM-DD")
	taskUpdateCmd.Flags().BoolVar(&updClearDue, "clear-due", false, "quita fecha de entrega")

	taskCmd.AddCommand(taskDeleteCmd)

	taskCmd.SilenceUsage = true
	for _, c := range []*cobra.Command{taskAddCmd, taskListCmd, taskGetCmd, taskUpdateCmd, taskDeleteCmd} {
		c.SilenceUsage = true
	}
}
