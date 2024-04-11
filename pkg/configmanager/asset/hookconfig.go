/*
Copyright 2024 KylinSoft  Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package asset

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	OneMB           = 1024 * 1024 // 1MB
	MaxHookFileSize = OneMB
	ShellFileType   = "shell"
	YAMLFileType    = "yaml"
	YAMLFileExt     = ".yaml"
	YMLFileExt      = ".yml"
)

//若传入的是一个目录，则会解析当前目录下的文件（注意：不会递归处理子目录下的文件）
//若传入的是一个文件而非目录，则会直接解析该文件并返回
func GetCmdHooks(conf *HookConf) error {
	if conf == nil {
		return errors.New("received nil pointer for HookConf parameter")
	}

	if conf.PreHookScript != "" {
		ShellFiles, err := getDirAndShells(conf.PreHookScript)
		if err != nil {
			return err
		}
		conf.ShellFiles = ShellFiles
	}

	if conf.PostHookYaml != "" {
		postHookFiles, err := getDirAndYamls(conf.PostHookYaml)
		if err != nil {
			return err
		}
		conf.PostHookFiles = postHookFiles
	}

	return nil
}

func getDirAndShells(p string) ([]ShellFile, error) {
	var (
		hookFiles     []ShellFile
		totalFileSize int64
	)
	fileInfo, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		hf, err := resolveFile(p, ShellFileType)
		if err != nil {
			return nil, err
		}
		hookFiles = append(hookFiles, hf)

		if fileInfo.Size() > MaxHookFileSize {
			return nil, fmt.Errorf("total size of shell script file exceeds the limit: %d bytes (max: %d bytes)", fileInfo.Size(), MaxHookFileSize)
		}
		return hookFiles, nil
	}

	files, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %s", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("the file list is empty")
	}
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}
		filePath := filepath.Join(p, file.Name())
		if err := checkHookFile(filePath, ShellFileType); err != nil {
			logrus.Debugf("failed to check hook file: %v", err)
			continue
		}
		hf, err := resolveFile(filePath, ShellFileType)
		if err != nil {
			logrus.Debugf("failed to resolve hook file %s: %v\n", filePath, err)
			continue
		}
		hookFiles = append(hookFiles, hf)

		// Get file info for size
		fileStat, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}
		totalFileSize += fileStat.Size()
	}
	if len(hookFiles) == 0 {
		return nil, fmt.Errorf("no valid hook files found in folder: %s", p)
	}
	if totalFileSize > MaxHookFileSize {
		return nil, fmt.Errorf("total size of shell script file in the directory exceeds the limit: %d bytes (max: %d bytes)", fileInfo.Size(), MaxHookFileSize)
	}

	return hookFiles, nil
}

func getDirAndYamls(path string) ([]string, error) {
	file, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if file.IsDir() {
		return resolvePostHookPath(path)
	}
	if err := checkHookFile(path, YAMLFileType); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func resolveFile(f string, fileType string) (ShellFile, error) {
	var hf ShellFile
	hf.Name = path.Base(f)

	if err := checkHookFile(f, fileType); err != nil {
		return hf, err
	}

	file, err := os.Open(f)
	if err != nil {
		return hf, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return hf, err
	}
	hf.Mode = int(fileInfo.Mode().Perm())

	content, err := io.ReadAll(file)
	if err != nil {
		return hf, err
	}
	hf.Content = content

	return hf, nil
}

func resolvePostHookPath(p string) ([]string, error) {
	var files []string

	rd, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %s", err)
	}

	if len(rd) == 0 {
		return nil, fmt.Errorf("empty directory: %s", p)
	}

	for _, fi := range rd {
		if err := checkHookFile(path.Join(p, fi.Name()), YAMLFileType); err == nil {
			files = append(files, path.Join(p, fi.Name()))
		} else {
			logrus.Debugf("failed to check hook file:%v", err)
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no valid hook files found in directory: %s", p)
	}
	return files, nil
}

func checkHookFile(fileName string, fileType string) error {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", fileInfo.Name())
	}

	switch fileType {
	case ShellFileType:
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			firstLine := scanner.Text()
			if !strings.HasPrefix(firstLine, "#!") {
				return fmt.Errorf("non-executable file: %s", fileName)
			}
		}
	case YAMLFileType:
		ext := filepath.Ext(fileName)
		if ext != YAMLFileExt && ext != YMLFileExt {
			return fmt.Errorf("%s is an invalid file extension", fileName)
		}
	default:
		logrus.Debugf("unknown file type")
	}
	return nil
}
