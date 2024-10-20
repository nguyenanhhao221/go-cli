/*
Copyright © 2024 Hao Nguyen <hao@haonguyen.tech>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Todos",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		isActive, err := cmd.Flags().GetBool("active")
		if err != nil {
			return err
		}
		return listAction(os.Stdout, apiRoot, isActive)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("active", false, "Show only active task")
}

func listAction(out io.Writer, url string, isActive bool) error {
	items, err := getAll(url)
	if err != nil {
		return err
	}
	return printAll(out, items, isActive)
}

func printAll(out io.Writer, items []item, isActive bool) error {
	w := tabwriter.NewWriter(out, 3, 2, 0, ' ', 0)
	for k, v := range items {
		done := "-"
		if v.Done {
			done = "X"
		}
		if isActive && v.Done {
			continue
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t\n", done, k+1, v.Task)
	}

	return w.Flush()
}
