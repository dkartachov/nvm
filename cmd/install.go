package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	version := args[0]

	fmt.Println(version)

	if version == "node" {
		// installLatestNode()
		installLatest("latest")

		return
	}

	semanticVersion := strings.Count(version, ".")

	switch semanticVersion {
	case 0:
		installLatest("latest-v" + version + ".x")
		break
	case 1:
		// TODO add logic to fetch latest patch given major and minor
		installLatest("v" + version + ".0")
		break
	default:
		installVersion(version)
		break
	}
}

func installLatest(endpoint string) {
	dirUrl, err := url.JoinPath(baseNode, endpoint)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Get(dirUrl)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalln(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	latestFilename, err := getLatestFileFromHtml(bytes)

	if err != nil {
		log.Fatalln(err)
	}

	finalUrl, err := url.JoinPath(dirUrl, latestFilename)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err = http.Get(finalUrl)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	regex, err := regexp.Compile(`\d+(\.\d+)+`)

	if err != nil {
		log.Fatalln(err)
	}

	version := regex.FindString(latestFilename)
	home, _ := os.UserHomeDir()

	os.Chdir(filepath.Join(home, ".nvm/node_versions"))

	targz.Extract(resp.Body)

	os.Rename(strings.ReplaceAll(latestFilename, ".tar.gz", ""), version)

	setExecPermissions(version)
}

func installLatestNode() {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond)
	s.Prefix = "Fetching latest version..."

	s.Start()

	time.Sleep(500 * time.Millisecond)

	url := baseNode + "/latest"

	resp, err := http.Get(url)

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
	s.FinalMSG = "Installed version " + version + " ✔️\n"

	home, _ := os.UserHomeDir()

	os.Chdir(filepath.Join(home, ".nvm/node_versions"))

	targz.Extract(resp.Body)

	os.Rename(strings.ReplaceAll(latestFilename, ".tar.gz", ""), version)

	setExecPermissions(version)

	s.Stop()
}

func installLatestMajor(major string) {
	resp, err := http.Get(baseNode + "/latest-v" + major + ".x")

	if err != nil {
		fmt.Println("installLatestMajor: error getting latest directory for major release " + major)
		log.Fatalln(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("installLatestMajor: error reading response body")
		log.Fatalln(err)
	}

	resp.Body.Close()

	latestMajorFilename, err := getLatestFileFromHtml(bytes)

	if err != nil {
		fmt.Println("installLatestMajor: error getting latest filename from html")
		log.Fatalln(err)
	}

	resp, err = http.Get(baseNode + "/latest-v" + major + ".x/" + latestMajorFilename)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	regex, err := regexp.Compile(`\d+(\.\d+)+`)

	if err != nil {
		log.Fatalln(err)
	}

	version := regex.FindString(latestMajorFilename)

	fmt.Println(version)

	// s.Prefix = "Installing version " + major + "..."
	// s.FinalMSG = "Installed version " + major + " ✔️"

	home, _ := os.UserHomeDir()

	os.Chdir(filepath.Join(home, ".nvm/node_versions"))

	targz.Extract(resp.Body)

	os.Rename(strings.ReplaceAll(latestMajorFilename, ".tar.gz", ""), version)

	setExecPermissions(version)

	// s.Stop()
}

func installLatestMinor(version string) {
	fmt.Println("install latest major " + version)
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

	os.Chdir(filepath.Join(home, ".nvm/node_versions"))

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
			return "", fmt.Errorf("getLatestFileFromHtml: could not get file from html")
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

	os.Chdir(filepath.Join(home, ".nvm/node_versions", version, "bin"))

	files, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal("setExecPermissions: ", err)
	}

	for _, file := range files {
		os.Chmod(file.Name(), 0777)
	}
}
