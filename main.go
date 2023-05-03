package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var goInstallDir string

func init() {
	goInstallDir = os.Getenv("GO_ROOT_PARENT_DIR")
}

func main() {
	if goInstallDir == "" {
		fmt.Println("GO_ROOT_PARENT_DIR env var empty.")
		return
	} else if sts, err := os.Stat(goInstallDir); err != nil {
		fmt.Println("access directory:", err)
		return
	} else if !sts.IsDir() {
		fmt.Println(goInstallDir, "is not a directory.")
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("Golang version switcher.")
		fmt.Println("Usage:")
		fmt.Printf("\t%s [version]\n", os.Args[0])
		listInstalled()
		return
	}

	goRoot := path.Join(goInstallDir, "go")
	if sts, err := os.Stat(goRoot); err != nil {
		fmt.Println("find old go root err:", err)
		return
	} else if !sts.IsDir() {
		fmt.Println("old go root not a dir:", goRoot)
		return
	}
	cmd := exec.Command("go", "version")
	oldVer := ""
	if out, err := cmd.Output(); err != nil {
		fmt.Println("run 'go version' err:", err.Error())
		return
	} else {
		// fmt.Println("version output:", string(out))
		ss := strings.SplitN(string(out), " ", 4)
		if len(ss) < 4 {
			fmt.Println("parse version fail:", string(out))
			return
		}
		oldVer = ss[2]
	}
	fmt.Println("old version is:", oldVer)
	renameOldTo := path.Join(goInstallDir, oldVer)

	targetVer := os.Args[1]
	targetRoot := path.Join(goInstallDir, fmt.Sprintf("go%s", targetVer))
	if sts, err := os.Stat(targetRoot); err != nil {
		fmt.Println("find target go root err:", err)
		return
	} else if !sts.IsDir() {
		fmt.Println("target go root not a dir:", targetRoot)
		return
	}

	if err := os.Rename(goRoot, renameOldTo); err != nil {
		fmt.Println("rename", goRoot, "to", renameOldTo, "fail:", err)
		return
	}
	if err := os.Rename(targetRoot, goRoot); err != nil {
		fmt.Println("rename", targetRoot, "to", goRoot, "fail:", err)
		return
	}
}

func listInstalled() {
	items, err := os.ReadDir(goInstallDir)
	if err != nil {
		fmt.Println("read dir error:", err)
		return
	}
	if out, err := exec.Command("go", "version").Output(); err != nil {
		fmt.Println("'go version' err:", err)
	} else {
		ss := strings.SplitN(string(out), " ", 4)
		if len(ss) < 4 {
			fmt.Println("parse version fail:", string(out))
			return
		}
		fmt.Println("Current:", ss[2])
	}
	fmt.Println("Install directory:", goInstallDir)
	fmt.Println("Go versions found: ")
	goDirRegex := regexp.MustCompile(`go\d\.\d{1,2}(\.\d{1,2})?`)
	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		if !goDirRegex.Match([]byte(item.Name())) {
			continue
		}
		if item.Name() == "go" {
			continue
		}
		fmt.Println("\t", item.Name())
	}
	fmt.Println("Enjoy it!")
}
