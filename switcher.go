package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

func runSwitcher(cx *cli.Context) error {
	if goRootDir == "" {
		return errors.New("go sdk parent dir not specified")
	} else if sts, err := os.Stat(goRootDir); err != nil {
		return fmt.Errorf("access dir err: %s", err.Error())
	} else if !sts.IsDir() {
		fmt.Println(goRootDir, "is not a directory.")
		return fmt.Errorf("'%s' is not a dir", goRootDir)
	}

	target := cx.Args().Get(0)
	if target == "" {
		listInstalled()
		return nil
	}
	return switch2Ver(target)
}

func switch2Ver(targetVer string) error {
	goRoot := path.Join(goRootDir, "go")
	if sts, err := os.Stat(goRoot); err != nil {
		return fmt.Errorf("find old go root err: %s", err.Error())
	} else if !sts.IsDir() {
		return fmt.Errorf("old go root not a dir: %s", goRoot)
	}
	oldVer, ok := fetchVersion(goRoot)
	if !ok {
		return fmt.Errorf("run 'go version' failed")
	}
	if oldVer == targetVer {
		fmt.Println("go version is already", oldVer)
		return nil
	}
	fmt.Println("old version is", oldVer)
	renameOldTo := path.Join(goRootDir, "go"+oldVer)

	targetRoot := path.Join(goRootDir, fmt.Sprintf("go%s", targetVer))
	if sts, err := os.Stat(targetRoot); err != nil {
		return fmt.Errorf("find target go root err: %s", err.Error())
	} else if !sts.IsDir() {
		return fmt.Errorf("target go root not a dir: %s", targetRoot)
	}

	fmt.Printf("changing to %s ...\n", targetVer)
	if err := os.Rename(goRoot, renameOldTo); err != nil {
		return fmt.Errorf("rename '%s' to '%s' failed: %s", goRoot, renameOldTo, err.Error())
	}
	if err := os.Rename(targetRoot, goRoot); err != nil {
		return fmt.Errorf("rename '%s' to '%s' failed: %s", targetRoot, goRoot, err.Error())
	}

	fmt.Println("\nEnjoy it!")
	return nil
}

func listInstalled() {
	fmt.Println("Golang version switcher.")
	fmt.Println("Usage:")
	fmt.Printf("    %s [version]\n", filepath.Base(os.Args[0]))

	items, err := os.ReadDir(goRootDir)
	if err != nil {
		fmt.Println("read dir error:", err)
		return
	}
	fmt.Println("Current:", currVer())
	fmt.Println("Dir:", goRootDir)
	fmt.Println("Found versions:")
	for _, item := range items {
		if !item.IsDir() || !strings.HasPrefix(item.Name(), "go") {
			continue
		}
		if item.Name() == "go" {
			continue
		}
		ver, ok := fetchVersion(path.Join(goRootDir, item.Name()))
		if !ok {
			continue
		}
		fmt.Printf("- %s\n", ver)
	}
	fmt.Println("\nEnjoy it!")
}

func currVer() string {
	goRoot := path.Join(goRootDir, "go")
	if sts, err := os.Stat(goRoot); err != nil {
		return "unknown"
	} else if !sts.IsDir() {
		return "unknown"
	}
	ver, ok := fetchVersion(goRoot)
	if !ok {
		return "unknown"
	}
	return ver
}

func currGoRoot() string {
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return ""
	}
	if runtime.GOOS == "windows" {
		out = bytes.ReplaceAll(out, []byte{'\\'}, []byte{'/'})
	}
	return path.Dir(string(out))
}

func fetchVersion(root string) (string, bool) {
	goExe := path.Join(root, "bin", goBinName)
	st, err := os.Stat(goExe)
	if err != nil {
		return "", false
	}
	if st.IsDir() {
		return "", false
	}
	out, err := exec.Command(goExe, "version").Output()
	if err != nil {
		return "", false
	}
	ss := bytes.Split(bytes.TrimSpace(out), []byte{' '})
	if len(ss) < 3 {
		return string(bytes.TrimSpace(out)), true
	}
	return string(bytes.TrimPrefix(ss[2], []byte{'g', 'o'})), true
}
