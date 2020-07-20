// Package resourceclasses Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// admin.yaml
// config.yaml
// datastore.yaml
// helm2.yaml
// helm3.yaml
// kubernetes.yaml
// monitoring.yaml
// networking.yaml
// security.yaml
// storage.yaml
// workloads.yaml
package resourceclasses

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _adminYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x92\x4f\x8b\xdb\x30\x10\xc5\xef\xfe\x14\x43\xae\x25\x76\x73\x0b\xba\x95\xf4\x52\x28\x85\xa6\xa5\xf7\xa9\x3c\xd8\x43\xf4\x2f\x33\x92\xfb\xe7\xd3\x17\x3b\x56\xca\x6e\x58\x96\x0d\xbb\xb7\x67\xf1\xf4\x7b\xa3\xf1\xc3\xc4\x3f\x48\x94\x63\x30\xe0\x29\x63\x8b\x29\xa9\x8d\x3d\xb5\x36\xfa\x6e\xda\xa1\x4b\x23\xee\x9a\x13\x87\xde\xc0\x91\x34\x16\xb1\x74\x70\xa8\xda\xcc\xf6\x1e\x33\x9a\x06\xc0\x0a\x61\xe6\x18\xbe\xb3\x27\xcd\xe8\x93\x81\x50\x9c\x6b\x00\x02\x7a\x32\xf0\xa1\xf7\x1c\x1a\x4d\x64\x67\x37\x85\x2c\x4c\x3a\xcb\xed\x6a\xf8\x82\x9e\x34\xa1\x25\x6d\x00\x00\x84\xce\x85\x85\x7a\x03\x59\x0a\x2d\x47\xf9\x4f\x22\xb3\x28\x80\x41\x62\x49\x5a\xbf\xb6\xb0\xd9\xac\x52\xd6\x09\xcd\x82\xbd\x02\x6b\xca\x67\xf6\x9c\xe1\x88\x61\x58\x73\x5e\x08\x75\xf3\x7d\xa9\xd7\x2b\xb5\xae\x05\xbe\x96\x98\xf1\x2e\x70\x55\xe7\x4a\xa8\xec\x83\x2b\x9a\x49\xe0\x18\xdd\x7d\xab\x91\x9f\x68\x5b\x2c\x79\x8c\xc2\x7f\x97\x7f\xd4\x9e\xf6\xda\x72\xbc\x99\xc1\x5e\xa2\x64\x4d\xba\xbe\xee\xed\x93\x1f\x47\x1e\xbe\x7d\x82\x8f\xc2\x13\xc9\xf3\xcb\xd4\x1c\x05\x07\x7a\xf2\x51\xca\xfd\x95\xc4\x36\x86\xb5\x75\x2a\xd6\xc0\x98\x73\x52\xd3\x75\xb6\x0f\x0f\x7b\x7f\xda\x6b\xb7\x98\xbb\x4a\xb2\x73\xe5\x49\x3b\x9c\x8b\xdc\xea\x34\xfc\x1f\x0c\xd8\xe3\x40\x9d\x4e\xc3\xbb\xdf\xde\xbd\x02\x3d\x85\x5b\xfa\xe5\xec\x17\xf1\x30\x66\x03\xbb\xf7\xcd\xbf\x00\x00\x00\xff\xff\xe4\x8c\x57\x40\xb9\x03\x00\x00")

func adminYamlBytes() ([]byte, error) {
	return bindataRead(
		_adminYaml,
		"admin.yaml",
	)
}

func adminYaml() (*asset, error) {
	bytes, err := adminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "admin.yaml", size: 953, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _configYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x90\xbd\x4e\xc3\x40\x10\x84\xfb\x7b\x8a\x55\x5a\x14\x9f\xd2\x20\x74\x6d\x6a\x1a\x40\xf4\xab\xf3\x72\x5e\xc5\xf7\xc3\xee\xda\xc0\xdb\x23\x3b\x36\x28\x7d\xba\xb9\xd3\xe8\x9b\xd9\xc1\xc6\xef\x24\xca\xb5\x04\xc8\x64\xd8\x61\x6b\x1a\x6b\x4f\x5d\xac\xd9\xcf\x27\x1c\xdb\x80\x27\x77\xe1\xd2\x07\x78\x21\xad\x93\x44\x3a\x8f\xa8\xea\x16\x7b\x8f\x86\xc1\x01\x44\x21\x34\xae\xe5\x8d\x33\xa9\x61\x6e\x01\xca\x34\x8e\x0e\xa0\x60\xa6\x00\xe7\x5a\x3e\x38\x39\x6d\x14\x17\x3b\x15\x13\x26\x5d\xe4\xf1\xc6\x01\xcf\xd8\xd4\x01\x00\x08\x7d\x4e\x2c\xd4\x07\x30\x99\x68\xfd\xb2\x9f\x46\x61\x55\x00\x49\xea\xd4\x74\x7f\x1d\xe1\x70\xd8\xa4\x6c\x1d\x03\xc4\x95\x98\xaf\xc0\x3d\xe6\x95\xa2\x90\xdd\x2b\x42\xff\x68\x1c\x6b\xd9\xee\x51\x89\x01\x06\xb3\xa6\xc1\xfb\xd8\x97\xdb\x49\x2f\x4f\xea\x57\xb3\xdf\x31\x71\x59\x93\xd4\x5f\xfb\x76\x3a\xa7\xff\x2a\xc0\x19\x13\x79\x9d\xd3\xc3\x77\x5e\xe6\xfc\x22\x4e\x83\x05\x78\x74\xbf\x01\x00\x00\xff\xff\xd5\x24\x64\x71\xba\x01\x00\x00")

func configYamlBytes() ([]byte, error) {
	return bindataRead(
		_configYaml,
		"config.yaml",
	)
}

func configYaml() (*asset, error) {
	bytes, err := configYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "config.yaml", size: 442, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _datastoreYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x90\x3d\x4f\x33\x31\x10\x84\x7b\xff\x8a\xe9\x5f\x25\x56\xa4\xb7\x40\x6e\x01\xd1\x23\x44\xbf\xb1\x57\x77\xab\x9c\x3f\xf0\xee\x05\xf8\xf7\xc8\xf9\x10\x44\xb4\x74\xf6\x6a\x66\x9e\xdd\xa1\x26\xaf\xdc\x55\x6a\x09\xc8\x6c\xb4\xa5\xd6\x34\xd6\xc4\xdb\x58\xb3\x3f\xee\x68\x69\x33\xed\xdc\x41\x4a\x0a\x78\x66\xad\x6b\x8f\x7c\xbf\x90\xaa\x1b\xf2\x44\x46\xc1\x01\xb1\x33\x99\xd4\xf2\x22\x99\xd5\x28\xb7\x80\xb2\x2e\x8b\x03\x0a\x65\x0e\x78\x20\x23\xb5\xda\xd9\x69\xe3\x38\x1c\xd4\xe4\xa9\xd7\xb5\x05\x1c\xd6\x3d\xa7\xfd\xe0\x39\x80\x8b\x75\x61\x1d\x8a\xcd\xc5\xfb\xb8\x90\x9a\x44\x65\xea\x71\x76\x00\xd0\xf9\x6d\x95\xce\x29\xc0\xfa\xca\xa7\x91\x7d\x36\x0e\xa7\x17\x30\x8d\x5c\xbd\xfe\x36\xb7\x80\xb3\xff\x7c\x47\x00\xff\xcc\x66\x75\x80\xc4\x5a\x2e\x78\xed\x31\x60\x36\x6b\x1a\xbc\x8f\xa9\xdc\x76\x73\xb8\x53\x7f\x12\xfb\x6b\x5c\x1c\xb5\xb0\xfa\x74\x3d\x76\xab\xc7\xe9\x7b\x3b\x48\xa6\x89\xbd\x1e\xa7\x7f\x1f\x79\xf9\x23\x42\x2b\xbf\x09\xe7\xd9\x3b\xcb\x34\x5b\xc0\x7f\xf7\x15\x00\x00\xff\xff\xa0\x4d\x25\xd8\xe4\x01\x00\x00")

func datastoreYamlBytes() ([]byte, error) {
	return bindataRead(
		_datastoreYaml,
		"datastore.yaml",
	)
}

func datastoreYaml() (*asset, error) {
	bytes, err := datastoreYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "datastore.yaml", size: 484, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _helm2Yaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x8f\xbd\x4e\x03\x31\x10\x84\x7b\x3f\xc5\xf6\x88\x58\x09\x0d\xda\x96\x86\x3a\x42\xf4\x2b\xdf\xea\x6c\xc5\x3f\x2b\xef\xe6\x80\xb7\x47\xbe\x70\x91\x10\xa9\x80\x6e\x7f\x46\x33\xf3\x91\xa4\x57\xee\x9a\x5a\x45\x28\x6c\xb4\x23\x11\x0d\x6d\xe2\x5d\x68\xc5\x2f\x7b\xca\x12\x69\xef\x4e\xa9\x4e\x08\x47\xd6\x76\xee\x81\x9f\x32\xa9\xba\x21\x9f\xc8\x08\x1d\x40\xe8\x4c\x96\x5a\x7d\x49\x85\xd5\xa8\x08\x42\x3d\xe7\xec\x00\x2a\x15\x46\x78\xe6\x5c\xe0\xe0\x54\x38\x0c\x39\x57\xeb\x89\x75\x8c\xf7\x90\x42\xab\xeb\x38\x16\xed\x01\x21\x9a\x89\xa2\xf7\x61\xaa\xdf\xfb\x9c\x1e\xd5\xaf\x72\xdf\xbf\xaa\x84\x51\x85\xc7\x9e\x99\x94\x77\xba\xcc\xab\x13\x80\x7d\x08\x23\xa4\x42\x33\x7b\x5d\xe6\xbb\xf7\x92\xff\x29\x43\xea\xad\x8c\xed\x7a\x01\x3e\x5e\xb4\xba\x9e\x84\x2c\x22\xf8\xc8\xb9\xf8\xe5\xb0\xf9\x8c\xdf\x95\xfd\xf7\xad\x86\xeb\x15\xfb\x36\xf4\x1f\xcd\x37\xb2\x9f\xb4\x6f\x9c\xe6\x68\x08\x0f\xee\x33\x00\x00\xff\xff\xda\xa8\x8c\x29\x48\x02\x00\x00")

func helm2YamlBytes() ([]byte, error) {
	return bindataRead(
		_helm2Yaml,
		"helm2.yaml",
	)
}

func helm2Yaml() (*asset, error) {
	bytes, err := helm2YamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "helm2.yaml", size: 584, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _helm3Yaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x8f\xbd\x4e\xc3\x40\x10\x84\xfb\x7b\x8a\xed\x11\x3e\x45\x6e\xd0\xb6\x34\xd4\x11\xa2\x5f\x9d\x57\xf6\x29\xf7\xb3\xba\xdd\x18\x78\x7b\x74\x0e\x8e\x84\x48\x05\xe9\xf6\x67\x34\x33\x1f\x49\x7c\xe3\xa6\xb1\x16\x84\xcc\x46\x03\x89\x68\xa8\x13\x0f\xa1\x66\xbf\x1e\x28\xc9\x42\x07\x77\x8a\x65\x42\x38\xb2\xd6\x73\x0b\xfc\x9c\x48\xd5\x75\xf9\x44\x46\xe8\x00\x42\x63\xb2\x58\xcb\x6b\xcc\xac\x46\x59\x10\xca\x39\x25\x07\x50\x28\x33\xc2\x0b\xa7\x0c\xa3\x53\xe1\xd0\xe5\x5c\xac\x45\xd6\x3e\x3e\x42\x0c\xb5\x6c\x63\x5f\xb4\x05\x84\xc5\x4c\x14\xbd\x0f\x53\xf9\xd9\xe7\xf4\xa4\x7e\x93\xfb\xf6\x5d\x25\xf4\x2a\xdc\xf7\xc4\xa4\x3c\xe8\x3a\x6f\x4e\x00\xf6\x29\x8c\x10\x33\xcd\xec\x75\x9d\x1f\x3e\x72\xba\x53\x86\x94\x5b\x19\xfb\xf5\x02\x7c\xbc\x68\x75\x3b\x09\xd9\x82\xe0\x17\x4e\xd9\xaf\xe3\xee\xd3\x7f\x57\xf6\xbf\xb7\xea\xae\x57\xec\xdb\xd0\xff\x34\xdf\xc9\x7e\xd3\xbe\x73\x9c\x17\x43\x18\xdd\x57\x00\x00\x00\xff\xff\xdd\x99\x03\x59\x48\x02\x00\x00")

func helm3YamlBytes() ([]byte, error) {
	return bindataRead(
		_helm3Yaml,
		"helm3.yaml",
	)
}

func helm3Yaml() (*asset, error) {
	bytes, err := helm3YamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "helm3.yaml", size: 584, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _kubernetesYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x92\xbd\x6e\xeb\x30\x0c\x85\x77\x3f\x05\xf7\x8b\x58\xc8\x16\x68\xbc\x1d\x8b\x2e\x6d\xd1\x9d\x91\x09\x9b\xb0\xf5\x53\x91\x4a\xd3\xb7\x2f\xac\xc4\x29\xd2\x78\x2a\x02\x74\x23\xe9\x63\x7d\x87\xe0\xc1\xc4\x6f\x94\x85\x63\xb0\xe0\x49\xb1\xc5\x94\xc4\xc5\x8e\x5a\x17\xbd\x39\x6c\x71\x4a\x03\x6e\x9b\x91\x43\x67\xe1\x99\x24\x96\xec\xe8\x61\x42\x91\x66\x96\x77\xa8\x68\x1b\x00\x97\x09\x95\x63\x78\x65\x4f\xa2\xe8\x93\x85\x50\xa6\xa9\x01\x08\xe8\xc9\xc2\x63\xd9\x53\x0e\xa4\x24\x8d\x24\x72\xf3\x2f\x14\x34\x33\xc9\x5c\x6e\x80\x5d\x0c\xb5\x9c\x1b\xc9\xce\xc2\xa0\x9a\xc4\x1a\xe3\xba\x70\xed\x69\xdc\x89\xa9\x72\x93\xcf\x76\xdc\x6c\x87\xc4\xec\x51\xd8\xb5\x72\xe8\xeb\x3b\x00\xfa\x99\xc8\x02\x7b\xec\xc9\xc8\xa1\xff\x77\xf4\xd3\x5d\x08\x29\xac\x11\x96\xe9\x69\xe1\xff\xb3\xb2\xf6\x09\x75\xb0\x60\x6a\x9d\xe9\xbd\x70\xa6\xce\x82\xe6\x42\x75\xf3\x93\xfc\x09\xdd\xc0\x81\x64\x5d\x75\x06\x9d\xa1\x7d\x8e\x25\xc9\xd2\x6d\xc0\x4d\x45\x94\x72\x7b\xdc\x8c\x3b\x69\x39\xde\x7c\xb8\x1a\x2f\x2b\x59\xf0\xdf\xcc\x1f\x36\xe0\x85\xf4\x4f\xac\x9c\xb0\x97\x30\xfc\xfe\x50\xe3\x25\x6f\x97\x3c\xac\xa7\xe1\x2e\x88\xe5\xf4\xb7\x71\xf8\x20\xee\x07\xb5\xb0\x6d\xbe\x02\x00\x00\xff\xff\xcc\x90\x2c\xbe\x69\x03\x00\x00")

func kubernetesYamlBytes() ([]byte, error) {
	return bindataRead(
		_kubernetesYaml,
		"kubernetes.yaml",
	)
}

func kubernetesYaml() (*asset, error) {
	bytes, err := kubernetesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "kubernetes.yaml", size: 873, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _monitoringYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x90\xbd\x4e\x2c\x31\x0c\x46\xfb\x3c\x85\xfb\xab\xdd\x68\xbb\x55\xda\x5b\x50\x21\x21\x84\xe8\xad\x8c\x35\x63\xed\x24\x31\xb6\xb3\xc0\xdb\xa3\xec\x1f\x20\x44\x47\x97\x38\x5f\xce\xb1\x8d\xc2\xcf\xa4\xc6\xad\x26\x28\xe4\xb8\x45\x11\xcb\x6d\xa2\x6d\x6e\x25\x1e\x77\xb8\xca\x82\xbb\x70\xe0\x3a\x25\x78\x24\x6b\x5d\x33\xfd\x5f\xd1\x2c\x8c\xf8\x84\x8e\x29\x00\x64\x25\x74\x6e\xf5\x89\x0b\x99\x63\x91\x04\xb5\xaf\x6b\x00\xa8\x58\x28\xc1\x7d\xab\xec\x4d\xb9\xce\xc1\x84\xf2\xf8\x82\xc2\x77\xda\xba\x24\x28\xb7\xc7\x6d\x6e\x4a\xcd\x86\x3b\x00\x50\x75\x65\xb2\x11\xde\x5c\x38\x0f\xda\x0a\xf9\x42\xdd\x02\x00\x80\xd2\x4b\x67\xa5\x29\x81\x6b\xa7\x53\xc9\xdf\x85\xd2\xe9\x04\x30\x0f\xbe\x5d\x6f\x9b\x5f\x45\x67\xd4\x79\xb6\x04\x72\x73\xd0\xb0\x70\x6e\xf5\xd2\x83\x69\x4e\xb0\xb8\x8b\xa5\x18\xf3\x54\xbf\x2f\xeb\xb0\xb7\x78\x0a\xc7\x2b\x2b\x8f\x3d\x91\xc5\x2f\x5e\x3b\xce\x9f\x6d\x02\x17\x9c\x29\xda\x71\xfe\xf7\x56\xd6\xbf\x52\x48\xfd\xa9\x38\xd7\x5e\x89\xe7\xc5\x13\xec\xc3\x47\x00\x00\x00\xff\xff\xd4\x1c\xc1\x23\xf7\x01\x00\x00")

func monitoringYamlBytes() ([]byte, error) {
	return bindataRead(
		_monitoringYaml,
		"monitoring.yaml",
	)
}

func monitoringYaml() (*asset, error) {
	bytes, err := monitoringYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "monitoring.yaml", size: 503, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _networkingYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\xbd\x6e\xeb\x30\x0c\x85\x77\x3d\x05\x91\xe1\x2e\x17\xb1\x90\xa1\x40\xa0\xb5\x5d\xba\x14\x45\x5b\x64\x17\x64\x42\x21\x6c\xfd\x94\x94\x9d\xe4\xed\x0b\x27\x76\x1d\x23\x4b\xd1\x66\xa3\xa8\x23\x7e\x47\x24\x6d\xa6\x1d\xb2\x50\x8a\x06\x02\x16\x5b\xd9\x9c\xc5\xa5\x1a\x2b\x97\x82\xee\x37\xb6\xcd\x7b\xbb\x51\x0d\xc5\xda\xc0\x1b\x4a\xea\xd8\xe1\x63\x6b\x45\xd4\x20\xaf\x6d\xb1\x46\x01\x38\x46\x5b\x28\xc5\x0f\x0a\x28\xc5\x86\x6c\x20\x76\x6d\xab\x00\xa2\x0d\x68\xe0\x1d\xb9\x27\x87\xf0\x0f\x9e\x48\x5c\xea\x91\x4f\x4a\x32\xba\xe1\x2d\xc6\xc2\x84\x32\x84\xeb\xa5\x5c\x14\x00\x00\xe3\x67\x47\x8c\xb5\x81\xc2\x1d\x9e\x53\xe5\x94\xd1\x9c\x23\x00\xcf\xa9\xcb\x32\x9d\xd6\xb0\x5a\x8d\x21\x8f\x6e\x0d\xc8\x5c\x6e\x22\x3c\x47\xcf\x28\xf2\x3b\x44\xc4\x72\x48\xdc\x50\xf4\x55\xb3\x95\x8a\xd2\xf7\x0d\x1e\x0b\xc6\xa1\x9b\x72\x63\x82\xae\x88\x93\x8b\x5d\x3a\x59\x8f\xfc\x37\x37\xfd\xa5\xc8\x62\x72\x3f\xa2\xbf\x5c\xbe\x01\xaf\xa9\x25\x47\xf7\x6d\xc5\x4c\x1e\x15\x79\x86\x90\x4b\x71\x1c\xb7\xb0\x33\xb0\x2f\x25\x8b\xd1\xda\xd5\x71\xb9\x7e\xcd\x56\xf4\x59\xac\xa7\x72\x6e\xd8\x3c\x14\x7d\x85\x95\xde\xcf\x2e\x81\x82\xf5\xa8\xa5\xf7\xff\x8f\xa1\xbd\x17\x22\xc7\x5b\xc4\x25\x77\x40\xf2\xfb\x62\xe0\x41\x7d\x05\x00\x00\xff\xff\x26\x14\x1b\xab\x49\x03\x00\x00")

func networkingYamlBytes() ([]byte, error) {
	return bindataRead(
		_networkingYaml,
		"networking.yaml",
	)
}

func networkingYaml() (*asset, error) {
	bytes, err := networkingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "networking.yaml", size: 841, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _securityYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x91\xcf\x4e\x33\x31\x0c\xc4\xef\x79\x0a\xab\xd7\x4f\x6d\xd4\xdb\x47\x6e\xa8\x2f\x80\x00\x71\xb7\xbc\x66\x6b\x75\xf3\x87\xd8\x59\xd8\xb7\x47\xbb\x6d\x29\x55\x4f\x20\x6e\x89\x35\x9e\xdf\x68\x8c\x45\x5e\xb8\xaa\xe4\x14\x20\xb2\xe1\x06\x4b\x51\xca\x1d\x6f\x28\x47\x3f\x6e\x71\x28\x7b\xdc\xba\x83\xa4\x2e\xc0\x23\x6b\x6e\x95\x78\x37\xa0\xaa\x9b\xe5\x1d\x1a\x06\x07\x40\x95\xd1\x24\xa7\x67\x89\xac\x86\xb1\x04\x48\x6d\x18\x1c\x40\xc2\xc8\x01\x9e\x98\x5a\x15\x9b\x9c\x16\xa6\x79\x81\x93\x55\x61\x9d\x9f\xeb\x93\xe6\x21\x0f\x42\x93\x03\x00\xa8\xfc\xd6\xa4\x72\x17\xc0\x6a\xe3\x65\x64\x53\xe1\xb0\xbc\x00\xfa\x9a\x5b\xd1\xf3\x6f\x0d\xe5\xb2\x39\xef\x1e\x43\x06\x28\xb9\xd3\x13\x77\x51\x08\xeb\x37\xdc\x8e\xab\xc9\xab\x10\xda\x32\xfe\x31\x74\xcc\x13\xf6\x5c\xaf\x0a\xbb\x89\x40\xd7\x90\xf5\x57\x1d\x75\x14\x62\xb8\x27\xca\x2d\xd9\xaf\xf8\xab\xd5\x0d\x4d\x8f\xb6\x78\x71\x15\xca\xe9\x54\xb2\x56\x0a\xb0\x37\x2b\x1a\xbc\xa7\x2e\x5d\x5f\xfa\xf0\x5f\xfd\x22\xf6\x67\x3b\x9a\x8f\xcc\xea\xcf\x15\x6e\x74\xec\x2f\xa1\x40\x22\xf6\xec\x75\xec\xff\x7d\xc4\xe1\x6f\x00\x25\xdd\x02\x8e\xb3\x77\x96\x7e\x6f\x01\xee\xdc\x67\x00\x00\x00\xff\xff\x19\x99\x30\xda\xb0\x02\x00\x00")

func securityYamlBytes() ([]byte, error) {
	return bindataRead(
		_securityYaml,
		"security.yaml",
	)
}

func securityYaml() (*asset, error) {
	bytes, err := securityYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "security.yaml", size: 688, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _storageYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x91\xcf\x4e\xc3\x30\x0c\xc6\xef\x79\x0a\x6b\x57\xb4\x46\x3b\x31\xe5\xba\x17\x40\x80\x76\xb7\x52\xab\xb3\x96\x7f\xc4\x6e\x81\xb7\x47\xed\xd2\x01\x9a\xb8\x20\xb8\xb9\xee\xa7\xdf\xcf\xfa\x82\x85\x8f\x54\x85\x73\x72\x10\x49\xb1\xc3\x52\xc4\xe7\x9e\x3a\x9f\xa3\x9d\x76\x18\xca\x09\x77\xe6\xcc\xa9\x77\xf0\x48\x92\xc7\xea\xe9\x10\x50\xc4\xcc\xf1\x1e\x15\x9d\x01\xf0\x95\x50\x39\xa7\x67\x8e\x24\x8a\xb1\x38\x48\x63\x08\x06\x20\x61\x24\x07\x4f\x9a\x2b\x0e\x64\xa4\x90\x9f\xf3\x94\xb4\x32\xc9\x3c\x6e\x5b\xe4\x61\x3e\x43\x94\x92\xc2\x31\x87\x31\x12\x1c\x02\x72\x14\x03\x00\x50\xe9\x65\xe4\x4a\xbd\x03\xad\x23\x2d\x2b\x7d\x2f\xe4\x96\x09\x60\xa8\x79\x2c\xb2\x7e\x6d\x61\xb3\x69\x63\x6d\x17\x3b\x28\x57\xfc\xb4\xd0\xfd\x0a\xff\xc1\x1f\x22\xfd\x97\xfa\xab\xb5\x15\x03\x4b\xa5\xbf\x33\xca\x05\xd1\x9d\xf7\xd2\x71\xbe\xb1\xb7\xdf\xfe\x2a\x60\x9f\x53\x6b\x5e\xaa\x77\x70\x52\x2d\xe2\xac\xf5\x7d\xfa\xfe\xfa\xe7\xbd\xd8\x25\x6c\x57\x5a\x83\xd8\x55\x29\xd3\xf0\x79\x1e\x70\xc4\x81\xac\x4c\xc3\xdd\x5b\x0c\x7f\xc2\x2f\xe9\x96\x7f\xd9\xbd\x12\x0f\x27\x75\x70\x6f\x3e\x02\x00\x00\xff\xff\xa5\x9e\x10\x09\xc2\x02\x00\x00")

func storageYamlBytes() ([]byte, error) {
	return bindataRead(
		_storageYaml,
		"storage.yaml",
	)
}

func storageYaml() (*asset, error) {
	bytes, err := storageYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "storage.yaml", size: 706, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _workloadsYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xbf\x0e\x1a\x31\x0c\x87\xf7\x7b\x0a\x8b\xb5\x82\x13\x9d\xaa\xac\x30\x75\xaa\xa0\x6a\x67\x93\xb8\x77\x29\x49\x9c\xda\x39\x28\x6f\x5f\xdd\x3f\x21\xa0\x1d\x8a\x6e\x4b\x2c\xdf\xf7\xd9\x17\xfd\x30\xfb\x6f\x24\xea\x39\x19\x88\x54\x70\x83\x39\xab\x65\x47\x1b\xcb\xb1\xbe\x6c\x31\xe4\x16\xb7\xd5\xd9\x27\x67\xe0\x40\xca\x9d\x58\xda\x05\x54\xad\xfa\x76\x87\x05\x4d\x05\x60\x85\xb0\x78\x4e\x5f\x7d\x24\x2d\x18\xb3\x81\xd4\x85\x50\x01\x24\x8c\x64\xe0\x3b\xcb\x39\x30\x3a\xad\x34\x93\xed\xbf\xa0\x54\xc4\x93\xf6\xc7\xf5\xd4\xb4\xa7\x1c\xf8\x16\x29\x15\xad\x00\x00\x84\x7e\x75\x5e\xc8\x19\x28\xd2\xd1\x50\x2a\xb7\x4c\x66\x38\x01\x34\xc2\x5d\xd6\xf9\xb6\x86\x7e\xf2\xe9\x22\xd3\xa0\x06\xdc\x03\x73\x36\x1d\x28\x07\x6f\x11\x8e\xb4\xa0\x4a\x46\xa8\xd2\xdf\x54\xfd\xcf\x81\x1d\xa7\x22\x1c\x02\xc9\x5b\xd6\xd5\xea\x5f\xce\x9e\x6e\x1f\xe0\xb3\xfe\x58\xb0\xd0\x8f\x2e\x2c\xbb\xaa\x4e\xd4\xa7\x5d\xf7\x48\x91\xd3\xb2\x2a\x37\x30\x9f\x44\x9f\xf9\xf4\x96\xe1\x84\xc5\xb6\x2f\x8a\x9f\x23\x6d\x86\xef\x84\xd3\xdd\xf0\xff\x38\x2b\x9c\x9e\x90\x5f\xd8\x2d\xf5\xe2\x79\x44\x79\xcb\x69\x0a\x8f\x8a\x35\xd0\x96\x92\xd5\xd4\xb5\x75\xe9\x31\xc2\xe7\x4f\x5a\x0f\xcd\xf5\xcc\xb0\x7d\x7a\x49\xeb\xeb\x9c\xc9\x8d\x5e\x9a\xfb\x28\xe0\x23\x36\x54\xeb\xa5\xf9\xf0\x3b\x86\x85\x0c\x39\xbd\x1a\xc6\xda\x95\x7c\xd3\x16\x03\x1f\xab\x3f\x01\x00\x00\xff\xff\xa3\xc0\x9b\x65\x8b\x04\x00\x00")

func workloadsYamlBytes() ([]byte, error) {
	return bindataRead(
		_workloadsYaml,
		"workloads.yaml",
	)
}

func workloadsYaml() (*asset, error) {
	bytes, err := workloadsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "workloads.yaml", size: 1163, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"admin.yaml":      adminYaml,
	"config.yaml":     configYaml,
	"datastore.yaml":  datastoreYaml,
	"helm2.yaml":      helm2Yaml,
	"helm3.yaml":      helm3Yaml,
	"kubernetes.yaml": kubernetesYaml,
	"monitoring.yaml": monitoringYaml,
	"networking.yaml": networkingYaml,
	"security.yaml":   securityYaml,
	"storage.yaml":    storageYaml,
	"workloads.yaml":  workloadsYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"admin.yaml":      {adminYaml, map[string]*bintree{}},
	"config.yaml":     {configYaml, map[string]*bintree{}},
	"datastore.yaml":  {datastoreYaml, map[string]*bintree{}},
	"helm2.yaml":      {helm2Yaml, map[string]*bintree{}},
	"helm3.yaml":      {helm3Yaml, map[string]*bintree{}},
	"kubernetes.yaml": {kubernetesYaml, map[string]*bintree{}},
	"monitoring.yaml": {monitoringYaml, map[string]*bintree{}},
	"networking.yaml": {networkingYaml, map[string]*bintree{}},
	"security.yaml":   {securityYaml, map[string]*bintree{}},
	"storage.yaml":    {storageYaml, map[string]*bintree{}},
	"workloads.yaml":  {workloadsYaml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
