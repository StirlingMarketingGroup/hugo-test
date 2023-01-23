package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	specFile, err := os.ReadFile("tscc.spec.json")
	if err != nil {
		log.Fatal(err)
	}

	var spec struct {
		External      map[string]string
		CompilerFlags map[string]any
	}
	if err = json.Unmarshal(specFile, &spec); err != nil {
		log.Fatal(err)
	}

	var tsccCmd []string

	tmpDir, err := os.MkdirTemp("", "compile-ts")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		log.Fatal(err)
	}

	tsccCmd = append(tsccCmd,
		"./node_modules/@tscc/tscc/dist/main.js",
		"--prefix", fmt.Sprintf("%s/", tmpDir),
		"--module", "app:ts/main.ts",
	)

	for pkg, variable := range spec.External {
		tsccCmd = append(tsccCmd,
			"--external", fmt.Sprintf("%v:%v", pkg, variable),
		)
	}

	tsccCmd = append(tsccCmd,
		"--", "--",
	)

	for k, v := range spec.CompilerFlags {
		tsccCmd = append(tsccCmd,
			fmt.Sprintf("--%s", k), fmt.Sprintf("%v", v),
		)
	}

	cmd := exec.Command(tsccCmd[0], tsccCmd[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	if err = os.RemoveAll("js"); err != nil {
		log.Fatal(err)
	}

	if err = os.RemoveAll("../static/js"); err != nil {
		log.Fatal(err)
	}

	err = filepath.WalkDir(tmpDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		file := strings.TrimPrefix(path, tmpDir)
		var dst string
		switch filepath.Ext(path) {
		case ".js":
			dst = "js"
		case ".map":
			dst = "../static/js"
		default:
			return nil
		}

		newPath := dst + file

		if err = os.MkdirAll(filepath.Dir(newPath), os.ModePerm); err != nil {
			return err
		}

		if err = os.Rename(path, newPath); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
