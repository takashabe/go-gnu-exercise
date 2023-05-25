package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

type cmdGrep struct {
	*CmdRoot

	invert bool
}

func (c *cmdGrep) New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "grep",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := os.Stdin
			var fileName string

			if len(args) == 2 {
				fileName = args[1]
			}

			pattern := args[0]
			re, err := regexp.Compile(pattern)
			if err != nil {
				return err
			}

			if fileName != "" {
				f, err := os.Open(fileName)
				if err != nil {
					return err
				}
				defer f.Close()
				reader = f
			}

			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				line := scanner.Text()
				match := re.MatchString(line)
				if (match && !c.invert) || (!match && c.invert) {
					fmt.Println(line)
				}
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&c.invert, "v", "v", false, "")

	return cmd
}
