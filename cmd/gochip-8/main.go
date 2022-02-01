package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cmd := newRootCommand()
	return cmd.Execute()
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "chip-8-term",
		Args: cobra.ExactArgs(0),
	}
	cmd.AddCommand(NewColorCmd(), NewStartCommand())
	return cmd
}
