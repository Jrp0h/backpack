package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jrp0h/backpack/cmd"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

func main() {
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/commands/" + strings.ToLower(base) + "/"
	}

	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		url := linkHandler(name)
		return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
	}

	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, "./docs/md", filePrepender, linkHandler)

	if err != nil {
		fmt.Printf("Failed to build markdown: %s\n", err.Error())
		os.Exit(1)
	}

	err = doc.GenManTree(cmd.RootCmd, &doc.GenManHeader{}, "./docs/man")

	fmt.Printf("Success\n")
}
