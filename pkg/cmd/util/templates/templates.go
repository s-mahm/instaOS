package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func GenerateTemplates(cmd *cobra.Command) {
	if cmd == nil {
		panic("nil root command")
	}
	cmd.SilenceUsage = true
	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	cmd.PersistentFlags().Lookup("help").Hidden = true
	cmd.SetUsageFunc(usageFunc())
	cmd.SetHelpFunc(helpFunc())
}

func usageFunc() func(*cobra.Command) error {
	return func(c *cobra.Command) error {
		t := template.New("usage")
		t.Funcs(template.FuncMap{
			"commandsAvailable": commandsAvailable,
			"flagsAvailable":    flagsAvailable,
		})
		template.Must(t.Parse(usageTemplate))
		return t.Execute(os.Stdout, c)
	}
}

func helpFunc() func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		t := template.New("help")
		template.Must(t.Parse(helpTemplate))
		if err := t.Execute(os.Stdout, c); err != nil {
			c.Println(err)
		}
	}
}

func ralign(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

func commandsAvailable(c *cobra.Command) string {
	cmds := []string{}
	for _, cmd := range c.Commands() {
		if cmd.IsAvailableCommand() && !cmd.HasSubCommands() {
			cmds = append(cmds, "  "+ralign(cmd.Name(), cmd.NamePadding())+"   "+cmd.Short)
		}
	}
	return strings.Join(cmds, "\n")
}

func flagsAvailable(c *cobra.Command) string {
	return c.Flags().FlagUsagesWrapped(120)
}

const usageTemplate = `Usage:

{{- if not .HasSubCommands}}  {{.UseLine}}{{end}}
{{- if .HasSubCommands}}  {{ .CommandPath}}{{- if .HasAvailableFlags}} [flags]{{end}} [commands]{{end}}

{{- if .HasSubCommands}}

Commands:
{{commandsAvailable .}}
{{- end}}


{{- if .HasAvailableFlags}}

Options:
{{flagsAvailable .}}
{{- end}}

Use "instaOS <command> --help" for more information about a given command.

For more help on how to use instaOS, head to https://github.com/s-mahm/instaOS
`

const helpTemplate = `{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
