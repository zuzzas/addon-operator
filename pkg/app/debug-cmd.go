package app

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	sh_app "github.com/flant/shell-operator/pkg/app"
	sh_debug "github.com/flant/shell-operator/pkg/debug"
)

func DefineDebugCommands(kpApp *kingpin.Application) {
	globalCmd := sh_app.CommandWithDefaultUsageTemplate(kpApp, "global", "manage global values")

	globalValuesCmd := globalCmd.Command("values", "Dump current global values.").
		Action(func(c *kingpin.ParseContext) error {
		dump, err := Global(sh_debug.DefaultClient()).Values(sh_debug.OutputFormat)
			if err != nil {
				return err
			}
			fmt.Println(string(dump))
			return nil
		})
	// -o json|yaml and --debug-unix-socket <file>
	AddOutputJsonYamlFlag(globalValuesCmd)
	sh_app.DefineDebugUnixSocketFlag(globalValuesCmd)

	globalConfigCmd := globalCmd.Command("config", "Dump global config values.").
		Action(func(c *kingpin.ParseContext) error {
		dump, err := Global(sh_debug.DefaultClient()).Config(sh_debug.OutputFormat)
			if err != nil {
				return err
			}
			fmt.Println(string(dump))
			return nil
		})
	// -o json|yaml and --debug-unix-socket <file>
	AddOutputJsonYamlFlag(globalConfigCmd)
	sh_app.DefineDebugUnixSocketFlag(globalConfigCmd)

	moduleCmd := sh_app.CommandWithDefaultUsageTemplate(kpApp, "module", "manage modules ant their values")

	moduleListCmd := moduleCmd.Command("list", "List available modules and their enabled status.").
		Action(func(c *kingpin.ParseContext) error {
			modules, err := Module(sh_debug.DefaultClient()).List(sh_debug.OutputFormat)
			if err != nil {
				return err
			}
			fmt.Println(string(modules))
			return nil
		})
	// -o json|yaml|text and --debug-unix-socket <file>
	sh_debug.AddOutputJsonYamlTextFlag(moduleListCmd)
	sh_app.DefineDebugUnixSocketFlag(moduleListCmd)

	var moduleName string
	moduleValuesCmd := moduleCmd.Command("values", "Dump module values by name.").
		Action(func(c *kingpin.ParseContext) error {
		dump, err := Module(sh_debug.DefaultClient()).Name(moduleName).Values(sh_debug.OutputFormat)
			if err != nil {
				return err
			}
			fmt.Println(string(dump))
			return nil
		})
	moduleValuesCmd.Arg("module_name", "").Required().StringVar(&moduleName)
	// -o json|yaml and --debug-unix-socket <file>
	AddOutputJsonYamlFlag(moduleValuesCmd)
	sh_app.DefineDebugUnixSocketFlag(moduleValuesCmd)

	moduleConfigCmd := moduleCmd.Command("config", "Dump module config values by name.").
		Action(func(c *kingpin.ParseContext) error {
			dump, err := Module(sh_debug.DefaultClient()).Name(moduleName).Config(sh_debug.OutputFormat)
			if err != nil {
				return err
			}
			fmt.Println(string(dump))
			return nil
		})
	moduleConfigCmd.Arg("module_name", "").Required().StringVar(&moduleName)
	// -o json|yaml and --debug-unix-socket <file>
	AddOutputJsonYamlFlag(moduleConfigCmd)
	sh_app.DefineDebugUnixSocketFlag(moduleConfigCmd)
}


func AddOutputJsonYamlFlag(cmd *kingpin.CmdClause) {
	cmd.Flag("output", "Output format: json|yaml.").Short('o').
		Default("yaml").
		EnumVar(&sh_debug.OutputFormat, "json", "yaml")
}


type GlobalRequest struct {
	client *sh_debug.Client
}

func Global(client *sh_debug.Client) *GlobalRequest {
	return &GlobalRequest{client: client}
}

func (gr *GlobalRequest) Values(format string) ([]byte, error) {
	url := fmt.Sprintf("http://unix/global/values.%s", format)
	return gr.client.Get(url)
}

func (gr *GlobalRequest) Config(format string) ([]byte, error) {
	url := fmt.Sprintf("http://unix/global/config.%s", format)
	return gr.client.Get(url)
}


type ModuleRequest struct {
	client *sh_debug.Client
	name string
}

func Module(client *sh_debug.Client) *ModuleRequest {
	return &ModuleRequest{client: client}
}

func (r *ModuleRequest) List(format string) ([]byte, error) {
	url := fmt.Sprintf("http://unix/module/list.%s", format)
	return r.client.Get(url)
}

func (mr *ModuleRequest) Name(name string) *ModuleRequest {
	mr.name = name
	return mr
}

func (mr *ModuleRequest) Values(format string) ([]byte, error) {
	url := fmt.Sprintf("http://unix/module/%s/values.%s", mr.name, format)
	return mr.client.Get(url)
}

func (mr *ModuleRequest) Config(format string) ([]byte, error) {
	url := fmt.Sprintf("http://unix/module/%s/config.%s", mr.name, format)
	return mr.client.Get(url)
}