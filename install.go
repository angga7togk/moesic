package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func getYtDlpPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".moesic", "bin", "yt-dlp")
}

func ytDlpURL() string {
	base := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/"
	switch runtime.GOOS {
	case "windows":
		return base + "yt-dlp.exe"
	case "darwin":
		return base + "yt-dlp_macos"
	case "linux":
		return base + "yt-dlp_linux"
	default:
		panic("unsupported OS	")
	}
}

func downloadYtDlp(dest string) error {
	fmt.Println("Downloading yt-dlp...")
	resp, err := http.Get(ytDlpURL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(dest, 0755)
		if err != nil {
			return err
		}
	}

	fmt.Println("yt-dlp downloaded successfully.")
	return nil
}

func ensureYtDlp() string {
	path := getYtDlpPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := downloadYtDlp(path)
		if err != nil {
			panic("failed to download yt-dlp: " + err.Error())
		}
	}
	return path
}

func InstallDependencies() {
	ensureYtDlp()
}
