package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sentire/pkg/models"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var describeCmd = &cobra.Command{
	Use:   "describe [command-path...]",
	Short: "Describe available commands and their schemas",
	Long:  "Output machine-readable JSON describing commands, arguments, flags, and output fields. Use without args to list all commands, or specify a command path (e.g. 'events list-issues') to describe a specific command.",
	RunE:  runDescribe,
}

func init() {
	rootCmd.AddCommand(describeCmd)
}

type commandDescription struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Args         []argDescription  `json:"args,omitempty"`
	Flags        []flagDescription `json:"flags,omitempty"`
	OutputFields []string          `json:"output_fields,omitempty"`
}

type argDescription struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type flagDescription struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

type describeOutput struct {
	Commands []commandDescription `json:"commands"`
}

var commandModelRegistry = map[string]reflect.Type{
	"events list-project":    reflect.TypeOf(models.Event{}),
	"events list-issue":      reflect.TypeOf(models.Event{}),
	"events list-issues":     reflect.TypeOf(models.Issue{}),
	"events get-event":       reflect.TypeOf(models.Event{}),
	"events get-issue":       reflect.TypeOf(models.Issue{}),
	"events get-issue-event": reflect.TypeOf(models.Event{}),
	"org list-projects":      reflect.TypeOf(models.Project{}),
	"org stats":              reflect.TypeOf(models.OrganizationStats{}),
	"projects list":          reflect.TypeOf(models.Project{}),
	"projects get":           reflect.TypeOf(models.Project{}),
	"inspect":                reflect.TypeOf(models.Event{}),
}

func runDescribe(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		cmdPath := strings.Join(args, " ")
		desc := describeCommand(rootCmd, cmdPath)
		if desc == nil {
			return NewInvalidInputError(fmt.Sprintf("unknown command: %s", cmdPath))
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(desc)
	}

	output := describeOutput{
		Commands: collectCommands(rootCmd, ""),
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

func collectCommands(cmd *cobra.Command, prefix string) []commandDescription {
	var result []commandDescription
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "completion" || child.Name() == "describe" {
			continue
		}
		fullName := child.Name()
		if prefix != "" {
			fullName = prefix + " " + child.Name()
		}
		if child.HasSubCommands() {
			result = append(result, collectCommands(child, fullName)...)
		} else {
			result = append(result, buildDescription(child, fullName))
		}
	}
	return result
}

func describeCommand(root *cobra.Command, path string) *commandDescription {
	parts := strings.Fields(path)
	current := root
	for _, part := range parts {
		found := false
		for _, child := range current.Commands() {
			if child.Name() == part {
				current = child
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	if current == root {
		return nil
	}
	desc := buildDescription(current, path)
	return &desc
}

func buildDescription(cmd *cobra.Command, fullName string) commandDescription {
	desc := commandDescription{
		Name:        fullName,
		Description: cmd.Short,
	}

	desc.Args = parseArgs(cmd.Use)

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		fd := flagDescription{
			Name:        f.Name,
			Type:        f.Value.Type(),
			Description: f.Usage,
		}
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
			fd.Default = f.DefValue
		}
		desc.Flags = append(desc.Flags, fd)
	})

	if modelType, ok := commandModelRegistry[fullName]; ok {
		desc.OutputFields = extractJSONFields(modelType)
	}

	return desc
}

func parseArgs(use string) []argDescription {
	parts := strings.Fields(use)
	var args []argDescription
	for _, p := range parts[1:] {
		if strings.HasPrefix(p, "<") && strings.HasSuffix(p, ">") {
			args = append(args, argDescription{
				Name:     strings.Trim(p, "<>"),
				Required: true,
			})
		} else if strings.HasPrefix(p, "[") && strings.HasSuffix(p, "]") {
			args = append(args, argDescription{
				Name:     strings.Trim(p, "[]"),
				Required: false,
			})
		}
	}
	return args
}

func extractJSONFields(t reflect.Type) []string {
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name != "" {
			fields = append(fields, name)
		}
	}
	return fields
}
