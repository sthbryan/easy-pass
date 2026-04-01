package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const repoOwner = "sthbryan"
const repoName = "easy-pass"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for updates or update easypass",
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if a new version is available",
	RunE:  runCheckUpdate,
}

var runUpdateCmd = &cobra.Command{
	Use:   "run",
	Short: "Download and install the latest version",
	RunE:  runUpdate,
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func runCheckUpdate(cmd *cobra.Command, args []string) error {
	current := version

	latest, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if latest == "" {
		fmt.Println("Could not determine latest version")
		return nil
	}

	if latest != current {
		fmt.Printf("New version available: %s (current: %s)\n", latest, current)
		fmt.Println("Run 'ep update run' to update")
	} else {
		fmt.Printf("Already on latest version: %s\n", current)
	}

	return nil
}

func runUpdate(cmd *cobra.Command, args []string) error {
	current := version

	latest, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if latest == "" {
		return fmt.Errorf("could not determine latest version")
	}

	if latest == current {
		fmt.Printf("Already on latest version: %s\n", current)
		return nil
	}

	fmt.Printf("Updating from %s to %s...\n", current, latest)

	assetName := getAssetName()
	downloadURL := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		repoOwner, repoName, latest, assetName,
	)

	tmpFile, err := downloadFile(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(tmpFile)

	if err := makeExecutable(tmpFile); err != nil {
		return fmt.Errorf("failed to make executable: %w", err)
	}

	installPath, err := getInstallPath()
	if err != nil {
		return fmt.Errorf("failed to determine install path: %w", err)
	}

	if err := os.Rename(tmpFile, installPath); err != nil {
		return fmt.Errorf("failed to install: %w (try running with sudo)", err)
	}

	fmt.Printf("Updated to %s\n", latest)
	return nil
}

func getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", nil
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func getAssetName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	ext := ""

	if os == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf("ep-%s-%s%s", os, arch, ext)
}

func getInstallPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return exe, nil
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "easypass-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func makeExecutable(path string) error {
	return os.Chmod(path, 0755)
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(checkUpdateCmd)
	updateCmd.AddCommand(runUpdateCmd)
}
