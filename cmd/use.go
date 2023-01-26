/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dkartachov/nvm/pkg/ver"
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

		if !ver.IsProperFormat(version) {
			return fmt.Errorf("invalid argument: %s", version)
		}

		home, _ := os.UserHomeDir()
		versions, _ := os.ReadDir(filepath.Join(home, ".nvm/node_versions"))
		semanticVersion := strings.Count(version, ".")

		switch semanticVersion {
		case 0:
			for i := 0; i < len(versions); i++ {
				majorVersion := strings.Split(versions[i].Name(), ".")[0]

				if majorVersion == version {
					return nil
				}
			}
			break
		case 1:
			for i := 0; i < len(versions); i++ {
				majorMinorPatch := strings.Split(versions[i].Name(), ".")
				majorVersion := majorMinorPatch[0]
				minorVersion := majorMinorPatch[1]
				majorMinorVersion := majorVersion + "." + minorVersion

				if majorMinorVersion == version {
					return nil
				}
			}
			break
		case 2:
			if ver.Exists(version) {
				return nil
			}
			break
		}

		return fmt.Errorf("version does not exist: %s", version)
	},
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		semanticVersion := strings.Count(version, ".")
		home, _ := os.UserHomeDir()

		switch semanticVersion {
		case 0:
			version = findLatestMajorVersion(version)
			break
		case 1:
			version = findLatestMinorVersion(version)
			break
		}

		current, err := os.Create(filepath.Join(home, ".nvm/current.txt"))

		if err != nil {
			log.Fatal(err)
		}

		defer current.Close()

		current.WriteString(version)

		fmt.Println("Now using version " + version)
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

func getMinorPatchVersion(version string) string {
	majorMinorPatch := strings.Split(version, ".")
	// majorVersion := majorMinorPatch[0]
	minorVersion := majorMinorPatch[1]
	patchVersion := majorMinorPatch[2]

	return minorVersion + "." + patchVersion
}

func findLatestMajorVersion(version string) string {
	home, _ := os.UserHomeDir()
	versions, _ := os.ReadDir(filepath.Join(home, ".nvm/node_versions"))

	currentLatestMajorVersion := "0.0.0"

	for i := 0; i < len(versions); i++ {
		semanticVersions := strings.Split(versions[i].Name(), ".")
		majorVersion := semanticVersions[0]

		if majorVersion == version {
			currentMinorPatchVersion, _ := strconv.ParseFloat(getMinorPatchVersion(currentLatestMajorVersion), 32)
			minorPatchVersion, _ := strconv.ParseFloat(getMinorPatchVersion(versions[i].Name()), 32)

			if minorPatchVersion >= currentMinorPatchVersion {
				currentLatestMajorVersion = versions[i].Name()
			}
		}
	}

	return currentLatestMajorVersion
}

func findLatestMinorVersion(version string) string {
	home, _ := os.UserHomeDir()
	versions, _ := os.ReadDir(filepath.Join(home, ".nvm/node_versions"))

	currentLatestMinorVersion := "0.0.0"

	for i := 0; i < len(versions); i++ {
		semanticVersions := strings.Split(versions[i].Name(), ".")
		majorMinorVersion := semanticVersions[0] + "." + semanticVersions[1]

		if majorMinorVersion == version {
			currentPatchVersion, _ := strconv.ParseInt(strings.Split(currentLatestMinorVersion, ".")[2], 10, 32)
			patchVersion, _ := strconv.ParseInt(strings.Split(versions[i].Name(), ".")[2], 10, 32)

			if patchVersion >= currentPatchVersion {
				currentLatestMinorVersion = versions[i].Name()
			}
		}
	}

	return currentLatestMinorVersion
}
