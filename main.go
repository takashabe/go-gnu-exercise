package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := CmdRoot{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}
	if err := cmd.New().Execute(); err != nil {
		fmt.Fprintf(cmd.ErrStream, "%v\n", err)
		// TODO: exit codeちゃんとエラーハンドリングする
		os.Exit(1)
	}
	os.Exit(0)
}

type CmdRoot struct {
	OutStream io.Writer
	ErrStream io.Writer
}

func (r *CmdRoot) New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gn",
	}

	ls := &cmdLs{CmdRoot: r}
	cmd.AddCommand(ls.New())

	cat := &cmdCat{CmdRoot: r}
	cmd.AddCommand(cat.New())
	return cmd
}
