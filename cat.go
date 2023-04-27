package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type cmdCat struct {
	*CmdRoot
}

func (c *cmdCat) New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cat",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println("Usage: cat [FILE]...")
				os.Exit(1)
			}

			for _, filename := range args {
				err := c.catFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error processing file '%s': %v\n", filename, err)
					os.Exit(1)
				}
			}
			return nil
		},
	}
	return cmd
}

func (c *cmdCat) catFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	return err
}
