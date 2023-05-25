package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type cmdGrep struct {
	*CmdRoot
}

func (c *cmdGrep) New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "grep",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[1]
			pattern := args[0]

			f, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, pattern) {
					fmt.Println(line)
				}
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
