package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"sigs.k8s.io/yaml"
)

func check(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var original map[string]interface{}
	err = yaml.Unmarshal(data, &original)
	if err != nil {
		return "", err
	}
	sorted, err := json.Marshal(&original)
	if err != nil {
		return "", err
	}

	var rd v1alpha1.ResourceDescriptor
	err = yaml.Unmarshal(data, &rd)
	if err != nil {
		return "", err
	}
	parsed, err := json.Marshal(rd)
	if err != nil {
		return "", err
	}

	// Then, Check them
	differ := diff.New()
	d, err := differ.Compare(sorted, parsed)
	if err != nil {
		fmt.Printf("Failed to unmarshal file: %s\n", err.Error())
		os.Exit(3)
	}

	if d.Modified() {
		config := formatter.AsciiFormatterConfig{
			ShowArrayIndex: true,
			Coloring:       true,
		}

		f := formatter.NewAsciiFormatter(original, config)
		result, err := f.Format(d)
		if err != nil {
			return "", err
		}
		return result, nil
	} else {
		data, err := yaml.Marshal(rd)
		if err != nil {
			return "", err
		}
		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func main() {
	err := filepath.Walk("./hub/resourcedescriptors/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext != ".yml" && ext != ".yaml" && ext != ".json" {
			return nil
		}

		d, err := check(path)
		if err != nil {
			return err
		}
		if d != "" {
			return fmt.Errorf("parsing error in file %s: %s", path, d)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
