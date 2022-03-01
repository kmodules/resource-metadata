/*
Copyright AppsCode Inc. and Contributors

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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
)

const iconURLPrefix = "https://cdn.appscode.com/k8s/icons/"

var (
	repoRoot     = os.ExpandEnv("$HOME/go/src/kmodules.xyz/resource-metadata")
	dirResources = path.Join(repoRoot, "hub/resourceeditors")
	dirMenus     = path.Join(repoRoot, "hub/menuoutlines")
	dirIcons     = path.Join(repoRoot, "icons")
)

func main() {
	var missing []string
	var img string
	allIcons := sets.NewString()

	err := filepath.Walk(dirResources, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var rd v1alpha1.ResourceEditor
		err = yaml.Unmarshal(data, &rd)
		if err != nil {
			return err
		}

		/*
			  - src: https://cdn.appscode.com/k8s/icons/acme.cert-manager.io/challenges.svg
			    type: image/svg+xml

				$HOME/go/src/kmodules.xyz/resource-metadata/icons/acme.cert-manager.io/challenges.png
				image/png
		*/
		rd.Spec.Icons = nil
		img = fmt.Sprintf("%s/%s.svg", groupDir(rd.Spec.Resource.Group), rd.Spec.Resource.Name)
		if !Exists(filepath.Join(dirIcons, img)) {
			missing = append(missing, img)
		} else {
			allIcons.Insert(img)
			rd.Spec.Icons = append(rd.Spec.Icons, v1alpha1.ImageSpec{
				Source: iconURLPrefix + img,
				Type:   "image/svg+xml",
			})
		}
		img = fmt.Sprintf("%s/%s.png", groupDir(rd.Spec.Resource.Group), rd.Spec.Resource.Name)
		if !Exists(filepath.Join(dirIcons, img)) {
			missing = append(missing, img)
		} else {
			allIcons.Insert(img)
			rd.Spec.Icons = append(rd.Spec.Icons, v1alpha1.ImageSpec{
				Source: iconURLPrefix + img,
				Type:   "image/png",
			})
		}

		data, err = yaml.Marshal(rd)
		if err != nil {
			return err
		}

		data, err = v1alpha1.FormatMetadata(data)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(path, data, 0o644)
	})
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(dirMenus, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var rc v1alpha1.MenuOutline
		err = yaml.Unmarshal(data, &rc)
		if err != nil {
			return err
		}

		if rc.Spec.Home != nil {
			rc.Spec.Home.Icons, missing = processIcons(rc.Spec.Home.Icons, allIcons, missing)
		}
		for i := range rc.Spec.Sections {
			rc.Spec.Sections[i].Icons, missing = processIcons(rc.Spec.Sections[i].Icons, allIcons, missing)
			for j := range rc.Spec.Sections[i].Items {
				rc.Spec.Sections[i].Items[j].Icons, missing = processIcons(rc.Spec.Sections[i].Items[j].Icons, allIcons, missing)
			}
		}

		data, err = yaml.Marshal(rc)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(path, data, 0o644)
	})
	if err != nil {
		panic(err)
	}

	var unusedIcons []string
	err = filepath.Walk(dirIcons, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext != ".svg" && ext != ".png" {
			return nil
		}

		rel, err := filepath.Rel(dirIcons, path)
		if err != nil {
			return err
		}
		if !allIcons.Has(rel) {
			unusedIcons = append(unusedIcons, rel)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	if len(missing) > 0 {
		fmt.Println("MISSING icons:")
		fmt.Println("**************")
		for _, name := range sets.NewString(missing...).List() {
			fmt.Println(name)
		}
		fmt.Println("____________________________________________________")
	}

	if len(unusedIcons) > 0 {
		fmt.Println("UNUSED icons:")
		fmt.Println("*************")
		for _, name := range sets.NewString(unusedIcons...).List() {
			fmt.Println(name)
		}
		fmt.Println("____________________________________________________")
	}
}

func groupDir(group string) string {
	if group == "" {
		return "core"
	}
	return group
}

func processIcons(icons []v1alpha1.ImageSpec, allIcons sets.String, missing []string) ([]v1alpha1.ImageSpec, []string) {
	m := map[string]string{} // mime -> url
	for _, entry := range icons {
		m[entry.Type] = entry.Source
	}

	icons = nil
	if u, ok := m["image/svg+xml"]; ok {
		img := strings.TrimPrefix(u, iconURLPrefix)
		if !Exists(filepath.Join(dirIcons, img)) {
			missing = append(missing, img)
		} else {
			allIcons.Insert(img)
			icons = append(icons, v1alpha1.ImageSpec{
				Source: u,
				Type:   "image/svg+xml",
			})
		}

		img = strings.TrimSuffix(img, ".svg") + ".png"
		if !Exists(filepath.Join(dirIcons, img)) {
			missing = append(missing, img)
		} else {
			allIcons.Insert(img)
			icons = append(icons, v1alpha1.ImageSpec{
				Source: iconURLPrefix + img,
				Type:   "image/png",
			})
		}
	}
	return icons, missing
}

// Exists reports whether the named file or directory Exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}
