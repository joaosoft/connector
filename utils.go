package webserver

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"

	"io"

	"fmt"
	"path"

	"github.com/joaosoft/errors"
)

func GetEnv() string {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	return env
}

func Exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(file string, obj interface{}) ([]byte, error) {
	var err error

	if !Exists(file) {
		return nil, errors.New("0", "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if obj != nil {
		if err := json.Unmarshal(data, obj); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func ReadFileLines(file string) ([]string, error) {
	lines := make([]string, 0)

	if !Exists(file) {
		return nil, errors.New("0", "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFile(file string, obj interface{}) error {
	if !Exists(file) {
		return errors.New("0", "file don't exist")
	}

	jsonBytes, _ := json.MarshalIndent(obj, "", "    ")
	if err := ioutil.WriteFile(file, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
