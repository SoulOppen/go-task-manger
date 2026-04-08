package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"strconv"
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
	addDependsOn   string
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Gestionar tareas",
}

var taskAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Crear tarea",
	RunE: func(cmd *cobra.Command, args []string) error {
		if shouldPromptAdd(cmd) {
			var err error
			addName, addDescription, addRelevance, addDue, addDependsOn, err = promptAddFields(cmd.InOrStdin(), cmd.OutOrStdout())
			if err != nil {
				return err
			}
		}

		due, err := task.ParseDueDate(addDue)
		if err != nil {
			return err
		}
		t := task.NewTask(addName, addDescription, addRelevance, due)
		if s := strings.TrimSpace(addDependsOn); s != "" {
			t.DependsOnID = &s
		}
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
	Short: "Listar tareas",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			tasks, err := repo.ListOrdered(cmd.Context())
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNOMBRE\tREL\tENTREGA\tDEPENDE_DE\tDESCRIPCION\tCREADO")
			for _, t := range tasks {
				due := "-"
				if t.DueDate != nil {
					due = t.DueDate.Format(task.DateLayout)
				}
				depCol := "-"
				if t.DependsOnID != nil && *t.DependsOnID != "" {
					if t.DependsOnName != "" {
						depCol = t.DependsOnName
					} else {
						id := *t.DependsOnID
						if len(id) > 8 {
							depCol = id[:8] + "…"
						} else {
							depCol = id
						}
					}
				}
				desc := strings.ReplaceAll(t.Description, "\n", " ")
				if len(desc) > 40 {
					desc = desc[:37] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n", t.ID, t.Name, t.Relevance, due, depCol, desc, t.CreatedAt.UTC().Format(time.RFC3339))
			}
			return w.Flush()
		})
	},
}

var taskGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Ver detalle de tarea",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveTaskID(args, cmd.InOrStdin(), cmd.OutOrStdout())
		if err != nil {
			return err
		}
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			t, err := repo.GetByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			due := "-"
			if t.DueDate != nil {
				due = t.DueDate.Format(task.DateLayout)
			}
			depLine := "-"
			if t.DependsOnID != nil && *t.DependsOnID != "" {
				if t.DependsOnName != "" {
					depLine = fmt.Sprintf("%s (%s)", *t.DependsOnID, t.DependsOnName)
				} else {
					depLine = *t.DependsOnID
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "id: %s\nnombre: %s\ndescripcion: %s\nrelevancia: %d\nentrega: %s\ndepende_de: %s\ncreado: %s\n",
				t.ID, t.Name, t.Description, t.Relevance, due, depLine, t.CreatedAt.UTC().Format(time.RFC3339))
			return nil
		})
	},
}

var (
	updName           string
	updDescription    string
	updRelevance      int
	updDue            string
	updClearDue       bool
	updDependsOn      string
	updClearDependsOn bool
)

var taskUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Actualizar tarea",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveTaskID(args, cmd.InOrStdin(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			existing, err := repo.GetByID(cmd.Context(), id)
			if err != nil {
				return err
			}

			if shouldPromptUpdate(cmd) {
				if err := promptUpdateFields(existing, cmd.InOrStdin(), cmd.OutOrStdout()); err != nil {
					return err
				}
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
			if updClearDependsOn {
				existing.DependsOnID = nil
			} else if cmd.Flags().Changed("depends-on") {
				s := strings.TrimSpace(updDependsOn)
				if s == "" {
					existing.DependsOnID = nil
				} else {
					existing.DependsOnID = &s
				}
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
	Short: "Eliminar tarea",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveTaskID(args, cmd.InOrStdin(), cmd.OutOrStdout())
		if err != nil {
			return err
		}
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			return repo.Delete(cmd.Context(), id)
		})
	},
}

var taskPickCmd = &cobra.Command{
	Use:   "pick",
	Short: "Elegir una tarea al azar",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			t, err := repo.PickRandom(cmd.Context())
			if err != nil {
				return err
			}
			due := "-"
			if t.DueDate != nil {
				due = t.DueDate.Format(task.DateLayout)
			}
			depLine := "-"
			if t.DependsOnID != nil && *t.DependsOnID != "" {
				if t.DependsOnName != "" {
					depLine = fmt.Sprintf("%s (%s)", *t.DependsOnID, t.DependsOnName)
				} else {
					depLine = *t.DependsOnID
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "id: %s\nnombre: %s\ndescripcion: %s\nrelevancia: %d\nentrega: %s\ndepende_de: %s\ncreado: %s\n",
				t.ID, t.Name, t.Description, t.Relevance, due, depLine, t.CreatedAt.UTC().Format(time.RFC3339))
			return nil
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
	taskAddCmd.Flags().StringVar(&addDependsOn, "depends-on", "", "UUID de la tarea de la que depende (opcional)")

	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskGetCmd)
	taskCmd.AddCommand(taskUpdateCmd)
	taskUpdateCmd.Flags().StringVar(&updName, "name", "", "nuevo nombre")
	taskUpdateCmd.Flags().StringVar(&updDescription, "description", "", "nueva descripcion")
	taskUpdateCmd.Flags().IntVar(&updRelevance, "relevance", 0, "nueva relevancia 1-10")
	taskUpdateCmd.Flags().StringVar(&updDue, "due", "", "nueva fecha YYYY-MM-DD")
	taskUpdateCmd.Flags().BoolVar(&updClearDue, "clear-due", false, "quita fecha de entrega")
	taskUpdateCmd.Flags().StringVar(&updDependsOn, "depends-on", "", "UUID de la tarea de la que depende")
	taskUpdateCmd.Flags().BoolVar(&updClearDependsOn, "clear-depends-on", false, "quita dependencia de otra tarea")

	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskPickCmd)

	taskCmd.SilenceUsage = true
	for _, c := range []*cobra.Command{taskAddCmd, taskListCmd, taskGetCmd, taskUpdateCmd, taskDeleteCmd, taskPickCmd, taskAddPromptCmd} {
		c.SilenceUsage = true
	}
}

func shouldPromptAdd(cmd *cobra.Command) bool {
	return !cmd.Flags().Changed("name") &&
		!cmd.Flags().Changed("description") &&
		!cmd.Flags().Changed("relevance") &&
		!cmd.Flags().Changed("due") &&
		!cmd.Flags().Changed("depends-on")
}

func shouldPromptUpdate(cmd *cobra.Command) bool {
	return !cmd.Flags().Changed("name") &&
		!cmd.Flags().Changed("description") &&
		!cmd.Flags().Changed("relevance") &&
		!cmd.Flags().Changed("due") &&
		!cmd.Flags().Changed("clear-due") &&
		!cmd.Flags().Changed("depends-on") &&
		!cmd.Flags().Changed("clear-depends-on")
}

func resolveTaskID(args []string, in io.Reader, out io.Writer) (string, error) {
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		return strings.TrimSpace(args[0]), nil
	}
	id, err := promptLine(bufio.NewReader(in), out, "ID de la tarea")
	if err != nil {
		return "", err
	}
	if id == "" {
		return "", fmt.Errorf("el id es obligatorio")
	}
	return id, nil
}

func promptAddFields(in io.Reader, out io.Writer) (name, description string, relevance int, due, dependsOn string, err error) {
	reader := bufio.NewReader(in)
	name, err = promptLine(reader, out, "Nombre")
	if err != nil {
		return "", "", 0, "", "", err
	}
	description, err = promptLine(reader, out, "Descripcion")
	if err != nil {
		return "", "", 0, "", "", err
	}
	relevance, err = promptInt(reader, out, "Relevancia (1-10, default 5)", 5)
	if err != nil {
		return "", "", 0, "", "", err
	}
	due, err = promptLine(reader, out, "Fecha de entrega YYYY-MM-DD (opcional)")
	if err != nil {
		return "", "", 0, "", "", err
	}
	dependsOn, err = promptLine(reader, out, "UUID tarea de la que depende (opcional, Enter si ninguna)")
	if err != nil {
		return "", "", 0, "", "", err
	}
	return name, description, relevance, due, dependsOn, nil
}

func promptUpdateFields(current *task.Task, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)

	name, err := promptLine(reader, out, fmt.Sprintf("Nombre [%s]", current.Name))
	if err != nil {
		return err
	}
	if name != "" {
		current.Name = name
	}

	description, err := promptLine(reader, out, fmt.Sprintf("Descripcion [%s]", current.Description))
	if err != nil {
		return err
	}
	if description != "" {
		current.Description = description
	}

	relevance, err := promptInt(reader, out, fmt.Sprintf("Relevancia (1-10) [%d]", current.Relevance), current.Relevance)
	if err != nil {
		return err
	}
	current.Relevance = relevance

	currentDue := "-"
	if current.DueDate != nil {
		currentDue = current.DueDate.Format(task.DateLayout)
	}
	dueInput, err := promptLine(reader, out, fmt.Sprintf("Fecha entrega YYYY-MM-DD [%s] (use '-' para limpiar)", currentDue))
	if err != nil {
		return err
	}
	switch strings.TrimSpace(dueInput) {
	case "":
		// keep
	case "-":
		current.DueDate = nil
	default:
		due, err := task.ParseDueDate(dueInput)
		if err != nil {
			return err
		}
		current.DueDate = due
	}

	curDep := "-"
	if current.DependsOnID != nil && *current.DependsOnID != "" {
		curDep = *current.DependsOnID
	}
	depInput, err := promptLine(reader, out, fmt.Sprintf("UUID depende_de [%s] (Enter mantener, '-' quitar)", curDep))
	if err != nil {
		return err
	}
	switch strings.TrimSpace(depInput) {
	case "":
		// keep
	case "-":
		current.DependsOnID = nil
	default:
		s := strings.TrimSpace(depInput)
		current.DependsOnID = &s
	}
	return nil
}

func promptLine(reader *bufio.Reader, out io.Writer, label string) (string, error) {
	fmt.Fprintf(out, "%s: ", label)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func promptInt(reader *bufio.Reader, out io.Writer, label string, fallback int) (int, error) {
	value, err := promptLine(reader, out, label)
	if err != nil {
		return 0, err
	}
	if value == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("valor numerico invalido: %w", err)
	}
	return n, nil
}
