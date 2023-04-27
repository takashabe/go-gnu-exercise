package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

type cmdLs struct {
	*CmdRoot

	long bool
}

func (c *cmdLs) New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			files, err := os.ReadDir(path)
			if err != nil {
				return fmt.Errorf("Error reading directory: %w", err)
			}

			for _, f := range files {
				info, err := f.Info()
				if err != nil {
					return fmt.Errorf("Error reading file info: %w", err)
				}
				if err := c.printFileInfo(info); err != nil {
					return fmt.Errorf("Error printing file info: %w", err)
				}
			}

			return nil
		},
	}
	cmd.PersistentFlags().BoolVarP(&c.long, "l", "l", false, "")
	return cmd
}

func (c *cmdLs) printFileInfo(info os.FileInfo) error {
	if c.long {
		mode := fileModeToString(info.Mode())
		size := info.Size()
		modTime := info.ModTime().Format("Jan _2 15:04")

		uid := strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Uid))
		gid := strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Gid))

		u, err := user.LookupId(uid)
		if err != nil {
			u.Username = uid
		}
		group, err := user.LookupGroupId(gid)
		if err != nil {
			group.Name = gid
		}

		fmt.Fprintf(c.OutStream, "%s %4s %4s %6d %s %s\n", mode, u.Username, group.Name, size, modTime, info.Name())
		return nil
	}
	fmt.Println(info.Name())
	return nil
}

func fileModeToString(mode os.FileMode) string {
	return mode.String()
}
