package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := []string{"."}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	for _, path := range args {
		if path != "." {
			err := visit(path, "")
			if err != nil {
				log.Printf("error visiting path %s: %v\n", path, err)
			}
		}
		path, err := os.Getwd()
		if err != nil {
			log.Printf("could not get working directory %v", err)
		}
		err = visit(path, "")
		if err != nil {
			log.Printf("error visiting path %s: %v\n", path, err)
		}
	}
}

func visit(root string, indent string) error {
	fileInfo, err := os.Stat(root)
	if err != nil {
		return err
	}
	if fileInfo.Name()[0] == '.' {
		return nil
	}

	fmt.Println(fileInfo.Name())

	if !fileInfo.IsDir() {
		return nil
	}

	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	names := []string{}

	for _, fileInfo := range fileInfos {
		if fileInfo.Name()[0] == '.' {
			continue
		}
		names = append(names, fileInfo.Name())
	}

	for i, name := range names {
		var add string
		if i == len(names)-1 {
			fmt.Printf("%s %s", indent, "└────")
			add = "     "
		} else {
			fmt.Printf("%s %s", indent, "├────")
			add = " │   "
		}
		if err := visit(filepath.Join(root, name), indent+add); err != nil {
			return err
		}
	}
	return nil
}
