package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"github.com/appscode/go/sets"
	"sigs.k8s.io/yaml"
)

var (
	repoRoot     = "/home/tamal/go/src/kmodules.xyz/resource-metadata"
	dirResources = path.Join(repoRoot, "hub/resourcedescriptors")
	dirClasses   = path.Join(repoRoot, "hub/resourceclasses")
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
		var rd v1alpha1.ResourceDescriptor
		err = yaml.Unmarshal(data, &rd)
		if err != nil {
			return err
		}

		/*
			  - src: https://cdn.appscode.com/k8s/icons/acme.cert-manager.io/challenges.svg
			    type: image/svg+xml

				/home/tamal/go/src/kmodules.xyz/resource-metadata/icons/acme.cert-manager.io/challenges.png
				image/png
		*/
		rd.Spec.Icons = nil
		img = fmt.Sprintf("%s/%s.svg", groupDir(rd.Spec.Resource.Group), rd.Spec.Resource.Name)
		if !Exists(filepath.Join(dirIcons, img)) {
			missing = append(missing, img)
		} else {
			allIcons.Insert(img)
			rd.Spec.Icons = append(rd.Spec.Icons, v1alpha1.ImageSpec{
				Source: "https://cdn.appscode.com/k8s/icons/" + img,
				Type:   "image/svg+xml",
			})
		}
		img = fmt.Sprintf("%s/%s.png", groupDir(rd.Spec.Resource.Group), rd.Spec.Resource.Name)
		if !Exists(filepath.Join(dirIcons, img)) {
			missing = append(missing, img)
		} else {
			allIcons.Insert(img)
			rd.Spec.Icons = append(rd.Spec.Icons, v1alpha1.ImageSpec{
				Source: "https://cdn.appscode.com/k8s/icons/" + img,
				Type:   "image/png",
			})
		}

		data, err = yaml.Marshal(rd)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(path, data, 0644)
	})
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(dirClasses, func(path string, info os.FileInfo, err error) error {
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
		var rc v1alpha1.ResourceClass
		err = yaml.Unmarshal(data, &rc)
		if err != nil {
			return err
		}

		rc.Spec.Icons, missing = processIcons(rc.Spec.Icons, allIcons, missing)
		for i := range rc.Spec.Entries {
			if len(rc.Spec.Entries[i].Path) > 0 {
				rc.Spec.Entries[i].Icons, missing = processIcons(rc.Spec.Entries[i].Icons, allIcons, missing)
			}
		}

		data, err = yaml.Marshal(rc)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(path, data, 0644)
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

	if len(missing) > 0 {
		fmt.Println("MISSING icons:")
		for _, name := range sets.NewString(missing...).List() {
			fmt.Println(name)
		}
		fmt.Println("____________________________________________________")
	}

	if len(unusedIcons) > 0 {
		fmt.Println("UNUSED icons:")
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
		img := strings.TrimPrefix(u, "https://cdn.appscode.com/k8s/icons/")
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
				Source: u,
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
