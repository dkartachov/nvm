/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)

		if err != nil {
			return err
		}

		version := args[0]
		home, _ := os.UserHomeDir()

		_, err = os.Stat(filepath.Join(home, ".nvm/versions/node", version))

		if err == nil {
			return nil
		}

		if os.IsNotExist(err) {
			log.Fatalf("version does not exist: %s", version)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		home, _ := os.UserHomeDir()

		bashrc, err := os.OpenFile(filepath.Join(home, ".bashrc"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		defer bashrc.Close()

		if err != nil {
			log.Fatal("Can't open .bashrc: ", err)
		}

		path := os.Getenv("PATH")
		paths := strings.Split(path, ":")

		for i := 0; i < len(paths); i++ {
			if strings.Contains(paths[i], ".nvm/versions/node") {
				paths[i] = paths[len(paths)-1]
				paths = paths[:len(paths)-1]
				path = strings.Join(paths, ":")

				break
			}
		}

		current, err := os.Create(filepath.Join(home, ".nvm/current.txt"))

		if err != nil {
			log.Fatal(err)
		}

		defer current.Close()

		current.WriteString(version)

		fmt.Println("PATH=" + filepath.Join(home, ".nvm/versions/node", version, "bin:") + path)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
