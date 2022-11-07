package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dkartachov/nvm/pkg/targz"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

const baseNode = "https://nodejs.org/dist"

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "install a version of node",
	Long:  `install a specific version of node`,
	Args: func(cmd *cobra.Command, args []string) error {
		return validate(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		install(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validate(cmd *cobra.Command, args []string) error {
	err := cobra.ExactArgs(1)(cmd, args)

	if err != nil {
		return err
	}

	if args[0] == "node" {
		return nil
	}

	re, err := regexp.Compile(`^(\d+\.)?(\d+\.)?(\*|\d+)$`)

	if err != nil {
		return err
	}

	if re.MatchString(args[0]) {
		return nil
	}

	return fmt.Errorf("invalid argument: %s", args[0])
}

func install(cmd *cobra.Command, args []string) {
	switch args[0] {
	case "node":
		installLatest()
		break
	default:
		installVersion(args[0])
		break
	}
}

func installLatest() {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond)
	s.Prefix = "Fetching latest version..."

	s.Start()

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get(baseNode + "/latest")

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("error: could not fetch latest directory")

		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	if err != nil {
		fmt.Println("error: could not read response body")
		fmt.Println(err)

		return
	}

	latestFilename, err := getLatestFileFromHtml(bytes)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err = http.Get(baseNode + "/latest/" + latestFilename)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	regex, err := regexp.Compile(`\d+(\.\d+)+`)

	if err != nil {
		log.Fatalln(err)
	}

	version := regex.FindString(latestFilename)

	s.Prefix = "Installing version " + version + "..."
	s.FinalMSG = "Installed version " + version + " ✔️"

	home, _ := os.UserHomeDir()

	os.Chdir(filepath.Join(home, ".nvm/versions/node"))

	targz.Extract(resp.Body)

	os.Rename(strings.ReplaceAll(latestFilename, ".tar.gz", ""), version)

	setExecPermissions(version)

	s.Stop()
}

func installVersion(version string) {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond)
	s.Prefix = "Fetching version " + version + "..."

	s.Start()

	time.Sleep(500 * time.Millisecond)

	file := getFileNameFromVersion(version) + ".tar.gz"

	resp, err := http.Get(baseNode + "/v" + version + "/" + file)

	if err != nil {
		fmt.Println("error: could not fetch version " + version)
		fmt.Println(err)

		return
	}

	defer resp.Body.Close()

	s.Prefix = "Installing version " + version + "..."
	s.FinalMSG = "Installed version " + version + " ✔️"

	home, _ := os.UserHomeDir()

	os.Chdir(filepath.Join(home, ".nvm/versions/node"))

	targz.Extract(resp.Body)

	os.Rename(getFileNameFromVersion(version), version)

	setExecPermissions(version)

	s.Stop()
}

func getLatestFileFromHtml(bytes []byte) (string, error) {
	htmlSrc := string(bytes)
	htmlTokens := html.NewTokenizer(strings.NewReader(htmlSrc))

	for {
		tokenType := htmlTokens.Next()

		switch tokenType {
		case html.ErrorToken:
			return "", fmt.Errorf("error: could not get file from html")
		case html.StartTagToken:
			token := htmlTokens.Token()

			if token.Data == "a" {
				file := token.Attr[len(token.Attr)-1].Val

				if strings.Contains(file, "linux-x64.tar.gz") {
					return file, nil
				}
			}
		}
	}
}

func getFileNameFromVersion(version string) string {
	return "node-v" + version + "-linux-x64"
}

func setExecPermissions(version string) {
	home, _ := os.UserHomeDir()

	os.Chdir(filepath.Join(home, ".nvm/versions/node", version, "bin"))

	files, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal("setExecPermissions: ", err)
	}

	for _, file := range files {
		os.Chmod(file.Name(), 0777)
	}
}
