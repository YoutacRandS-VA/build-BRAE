// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// The update-readmes.go tool creates or updates README.md files in
// the golang.org/x/build tree. It only updates files if they are
// missing or were previously generated by this tool. If the file
// contains a "<!-- End of auto-generated section -->" comment,
// the tool leaves content in the rest of the file unmodified.
//
// The auto-generated Markdown contains the package doc synopsis
// and a link to pkg.go.dev for the API reference.
package main

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root, err := build.Import("golang.org/x/build", "", build.FindOnly)
	if err != nil {
		log.Fatalf("failed to find golang.org/x/build root: %v", err)
	}
	err = filepath.Walk(root.Dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		rest := strings.TrimPrefix(strings.TrimPrefix(path, root.Dir), "/")
		switch rest {
		case "env", "version", "vendor":
			return filepath.SkipDir
		}
		pkgName := "golang.org/x/build/" + filepath.ToSlash(rest)

		bctx := build.Default
		bctx.Dir = path // Set Dir since some x/build packages are in nested modules.
		pkg, err := bctx.Import(pkgName, "", 0)
		if err != nil {
			// Skip.
			return nil
		}
		if pkg.Doc == "" {
			// There's no package comment, so don't create an empty README.
			return nil
		}
		if _, err := os.Stat(filepath.Join(pkg.Dir, "README")); err == nil {
			// Directory has exiting README; don't touch.
			return nil
		}
		readmePath := filepath.Join(pkg.Dir, "README.md")
		exist, err := os.ReadFile(readmePath)
		if err != nil && !os.IsNotExist(err) {
			// A real error.
			return err
		}
		const header = "Auto-generated by x/build/update-readmes.go"
		if len(exist) > 0 && !bytes.Contains(exist, []byte(header)) {
			return nil
		}
		var footer []byte
		if i := bytes.Index(exist, []byte("<!-- End of auto-generated section -->")); i != -1 {
			footer = exist[i:]
		}
		newContents := []byte(fmt.Sprintf(`<!-- %s -->

[![Go Reference](https://pkg.go.dev/badge/%s.svg)](https://pkg.go.dev/%s)

# %s

%s
%s`, header, pkgName, pkgName, pkgName, pkg.Doc, footer))
		if bytes.Equal(exist, newContents) {
			return nil
		}
		if err := os.WriteFile(readmePath, newContents, 0644); err != nil {
			return err
		}
		log.Printf("Wrote %s", readmePath)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}