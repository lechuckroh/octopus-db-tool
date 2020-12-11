package util

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// MkdirPackage creates nested package directories.
func MkdirPackage(basedir, pkgName string) (string, error) {
	if pkgName == "" {
		return "", nil
	}
	dir := path.Join(append([]string{basedir}, strings.Split(pkgName, ".")...)...)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	return dir, nil
}

func WriteStringToFile(filename string, s string) error {
	if err := ioutil.WriteFile(filename, []byte(s), 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", filename)
	return nil
}

func WriteBytesToFile(filename string, bytes []byte) error {
	if err := ioutil.WriteFile(filename, bytes, 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", filename)
	return nil
}

func WriteLinesToFile(filename string, lines []string) error {
	if err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", filename)
	return nil
}

func Mkdir(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	return dir, nil
}
