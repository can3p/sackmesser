/*
Copyright © 2024 Dmitrii Petrov <dpetroff@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/can3p/sackmesser/pkg/cobrahelpers"
	"github.com/can3p/sackmesser/pkg/operations"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/simpleyaml"
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/spf13/cobra"
)

func ModCommand() *cobra.Command {
	var deleteField bool
	var jsonValue bool
	var inputFormat string
	var outputFormat string

	var modCmd = &cobra.Command{
		Use:   "mod",
		Short: "modify input",
		Long:  "Parse incoming object and update (or delete) the specified field",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := io.ReadAll(os.Stdin)

			if err != nil {
				return err
			}

			// here we should dispatch against type to parse,
			// but we don't do that atm
			var root types.Node
			switch inputFormat {
			case "yaml":
				root, err = simpleyaml.Parse(input)
			case "json":
				root, err = simplejson.Parse(input)
			default:
				return fmt.Errorf("Unkonwn input format: %s", inputFormat)
			}

			if err != nil {
				return err
			}

			if len(args) > 0 {
				path := strings.Split(strings.TrimLeft(args[0], "."), ".")

				if err != nil {
					return err
				}

				if deleteField {
					if err := operations.Delete(root, path); err != nil {
						return err
					}
				} else {

					if len(args) < 2 {
						panic("set operation requires at least two arguments")
					}

					var val any = args[1]

					if jsonValue {
						if err := json.Unmarshal([]byte(args[1]), &val); err != nil {
							return err
						}
					}

					if err := operations.Set(root, path, val); err != nil {
						return err
					}
				}
			}

			var outputRoot types.Node

			switch outputFormat {
			case "yaml":
				outputRoot = simpleyaml.FromNode(root)
			case "json":
				outputRoot = simplejson.FromNode(root)
			default:
				return fmt.Errorf("Unkonwn ouput format: %s", inputFormat)
			}

			out, err := outputRoot.Serialize()

			if err != nil {
				return err
			}

			fmt.Print(string(out))

			return nil
		},
	}

	modCmd.Flags().BoolVarP(&deleteField, "delete", "d", false, "delete field from a given path")
	modCmd.Flags().BoolVarP(&jsonValue, "json", "j", false, "parse value as json")
	modCmd.Flags().Var(cobrahelpers.NewEnumFlag(&inputFormat, "json", "json", "yaml"), "input-format", `input format: json or yaml`)
	modCmd.Flags().Var(cobrahelpers.NewEnumFlag(&outputFormat, "json", "json", "yaml"), "output-format", `input format: json or yaml`)

	return modCmd
} // modCmd represents the mod command

func init() {
	rootCmd.AddCommand(ModCommand())
}
