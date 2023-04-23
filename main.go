package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var goRootDir string

func init() {
	goRootDir = os.Getenv("GO_ROOT_PARENT_DIR")
}

func main() {
	if goRootDir == "" {
		fmt.Println("GO_ROOT_PARENT_DIR env var empty.")
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("target go version not privided.")
		return
	}

	goRoot := path.Join(goRootDir, "go")
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
	renameOldTo := path.Join(goRootDir, oldVer)

	targetVer := os.Args[1]
	targetRoot := path.Join(goRootDir, fmt.Sprintf("go%s", targetVer))
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
