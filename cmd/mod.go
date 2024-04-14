/*
Copyright © 2024 Dmitrii Petrov <dpetroff@gmail.com>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/can3p/sackmesser/pkg/cobrahelpers"
	"github.com/can3p/sackmesser/pkg/operations"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/can3p/sackmesser/pkg/traverse/simpleyaml"
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ModCommand() *cobra.Command {
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
				ops := []*operations.OpInstance{}
				parser := operations.NewParser()

				for _, arg := range args {
					op, err := parser.Parse(arg)

					if err != nil {
						return errors.Wrapf(err, "Invalid operation: [%s]", arg)
					}

					ops = append(ops, op)
				}

				for _, op := range ops {
					if err := op.Apply(root); err != nil {
						return errors.Wrapf(err, "failed to apply operation: %s", op.String())
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

			// nolint:forbidigo
			fmt.Print(string(out))

			return nil
		},
	}

	modCmd.Flags().Var(cobrahelpers.NewEnumFlag(&inputFormat, "json", "json", "yaml"), "input-format", `input format: json or yaml`)
	modCmd.Flags().Var(cobrahelpers.NewEnumFlag(&outputFormat, "json", "json", "yaml"), "output-format", `input format: json or yaml`)

	return modCmd
} // modCmd represents the mod command

func init() {
	rootCmd.AddCommand(ModCommand())
}
