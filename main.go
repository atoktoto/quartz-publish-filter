package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	srcDirPtr := flag.String("input", "public", "source directory")
	targetDirPtr := flag.String("output", "public-filtered", "target directory")

	srcDir := *srcDirPtr
	targetDir := *targetDirPtr
	publicTag := "#public"

	err := filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failure accessing a path %q: %v\n", path, err)
			return err
		}

		targetPath := targetDir + strings.TrimPrefix(path, srcDir)

		if info.IsDir() {
			fmt.Printf("Creating directory: %q\n", targetPath)
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		if info.IsDir() == false && strings.HasSuffix(info.Name(), ".md") {
			hasPublicTag, _ := hasTag(path, publicTag)
			if hasPublicTag {
				_, err := simpleCopy(path, targetPath)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
	}
}

func simpleCopy(src, dst string) (int64, error) {
	fmt.Printf("Copying: %q -> %q\n", src, dst)

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func hasTag(path string, tag string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hasTag := strings.Contains(scanner.Text(), tag)
		if hasTag {
			return true, err
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
