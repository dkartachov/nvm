/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dkartachov/nvm/pkg/ver"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "remove a node version from the computer",
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

		if !ver.Exists(version) {
			return fmt.Errorf("Version " + version + " does not exist")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()

		if err != nil {
			log.Fatalln("Could not retrieve user home directory", err)
		}

		version := args[0]

		err = os.RemoveAll(filepath.Join(home, ".nvm/node_versions", version))

		if err != nil {
			log.Fatalln("Could not remove node version " + version)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
