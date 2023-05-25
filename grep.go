package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type cmdGrep struct {
	*CmdRoot

	invert     bool
	ignoreCase bool
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
			if c.ignoreCase {
				pattern = "(?i)" + pattern
			}
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
				idx := re.FindAllStringIndex(line, -1)

				p := ""
				if len(idx) > 0 {
					lastIndex := 0
					for _, i := range idx {
						// 前回マッチしたところから今回マッチしたところまでを追加
						p = fmt.Sprintf("%s%s", p, line[lastIndex:i[0]])
						// マッチ文字列を色つける
						p = fmt.Sprintf("%s%s", p, color.GreenString(line[i[0]:i[1]]))
						lastIndex = i[1]
					}
					// 残った文字列を追加
					p = fmt.Sprintf("%s%s", p, line[lastIndex:])
					fmt.Println(p)
				}
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&c.invert, "v", "v", false, "")
	cmd.PersistentFlags().BoolVarP(&c.ignoreCase, "i", "i", false, "")

	return cmd
}
