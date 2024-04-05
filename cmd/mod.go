/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/can3p/sackmesser/pkg/operations"
	"github.com/can3p/sackmesser/pkg/traverse/simplejson"
	"github.com/spf13/cobra"
)

func ModCommand() *cobra.Command {
	var deleteField bool
	var jsonValue bool

	var modCmd = &cobra.Command{
		Use:   "mod",
		Short: "modify input",
		Long:  "Parse incoming object and update (or delete) the specified field",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := io.ReadAll(os.Stdin)

			if err != nil {
				return err
			}

			// here we should dispatch against type to parse,
			// but we don't do that atm
			root, err := simplejson.Parse(input)

			if err != nil {
				return err
			}

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

			out, err := root.Serialize()

			if err != nil {
				return err
			}

			fmt.Print(string(out))

			return nil
		},
	}

	modCmd.Flags().BoolVarP(&deleteField, "delete", "d", false, "delete field from a given path")
	modCmd.Flags().BoolVarP(&jsonValue, "json", "j", false, "parse value as json")

	return modCmd
} // modCmd represents the mod command

func init() {
	rootCmd.AddCommand(ModCommand())
}
