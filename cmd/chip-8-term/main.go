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
	cmd := newFoobarCommand()
	return cmd.Execute()
}

func newFoobarCommand() *cobra.Command {
	var (
		fps uint8
	)

	cmd := &cobra.Command{
		Use:  "chip-8-term",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Printf("%v\n", args)
			//
			return nil
		},
	}
	cmd.PersistentFlags().Uint8Var(&fps, "keyboard-hz", 60, "reciprocal of duration of key pressed (default: 60Hz)")
	return cmd
}
