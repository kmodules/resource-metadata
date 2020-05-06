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

var _adminYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x90\x4d\x6b\xe3\x40\x0c\x86\xef\xfe\x15\x22\xd7\x65\xed\xcd\x2d\xcc\x6d\xc9\x5e\x16\x4a\xa1\x29\xf4\x52\x7a\x50\xc7\xc2\x16\xf1\x7c\x44\x9a\x71\x3f\x7e\x7d\xb1\xe3\x49\x69\x29\xa1\x39\xf4\x26\x84\xde\xe7\x91\x84\x91\xef\x48\x94\x83\x37\xe0\x28\x61\x8d\x31\xaa\x0d\x2d\xd5\x36\xb8\x66\x5c\xe3\x10\x7b\x5c\x57\x7b\xf6\xad\x81\x1d\x69\xc8\x62\x69\x3b\xa0\x6a\x35\x8d\xb7\x98\xd0\x54\x00\x1e\x1d\x19\xf8\xdb\x3a\xf6\x95\x46\xb2\x53\xef\x89\xb8\xeb\x93\x81\xf5\x9f\x0a\x80\x6d\xf0\x3a\x75\x7f\x83\x8a\x35\xd0\xa7\x14\xd5\x34\x8d\x6d\xfd\x47\xe5\x7e\xa3\xcd\x3c\xdc\xc8\x62\xb3\x93\x8d\xb4\xc1\x89\x5e\xeb\xd8\x55\x00\x00\xe9\x25\x92\x01\x76\xd8\x51\xa3\x63\xf7\xeb\xd9\x0d\x15\x00\xf9\x24\x4c\x8b\xe8\xb8\xd4\x35\x3a\xd2\x88\x96\x74\xce\x09\x1d\x32\x0b\xb5\x06\x92\x64\x7a\x47\xcd\x15\x40\x27\x21\x47\x35\x70\xbf\x5a\x3d\x2c\xad\xb2\x87\x99\x81\x27\x54\xe1\x5f\xb1\xe3\x04\x3b\xf4\xdd\x62\xf8\x36\x6e\x98\x92\x52\x82\x85\x57\x9e\x0c\x37\x39\x24\xbc\x10\x59\xaa\x43\xc9\x16\xea\x76\xc8\x9a\x48\x60\x17\x86\x4b\x1f\x21\x8f\x68\x6b\xcc\xa9\x0f\xc2\xaf\x98\x38\xf8\x7a\xbf\xd1\x9a\xc3\x17\x7e\x7b\xd4\xc8\x62\x39\xdd\xf4\xb3\xd6\xcf\xba\xed\xed\x7f\xf8\x27\x3c\x92\x9c\x7b\x9f\xa6\x20\xd8\xd1\x99\x63\x94\xdb\x85\xf2\x16\x00\x00\xff\xff\x1a\x37\x68\xe6\x27\x03\x00\x00")

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

	info := bindataFileInfo{name: "admin.yaml", size: 807, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _configYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x8f\xbd\x4e\xc4\x30\x10\x84\x7b\x3f\xc5\x28\x2d\x22\xd6\x35\x08\x6d\x7b\x05\x15\x0d\x48\x34\x88\x62\xe5\x2c\x8e\x75\xf1\x0f\x5e\x27\xc0\xdb\xa3\xe4\x72\x02\xea\xeb\xac\xd1\x78\xbe\x6f\xb9\x84\x17\xa9\x1a\x72\x22\x44\x69\xdc\x73\x29\xea\xf2\x20\xbd\xcb\xd1\x2e\x07\x9e\xca\xc8\x07\x73\x0a\x69\x20\x3c\x89\xe6\xb9\x3a\x39\x4e\xac\x6a\xd6\xfa\xc0\x8d\xc9\x00\x89\xa3\x10\x8e\x39\xbd\x07\x6f\xb4\x88\x5b\x43\x2e\xe1\xa1\xe6\xb9\x10\xba\xce\x00\x9f\x12\xfc\xd8\x08\x77\x06\x08\x2e\x27\x5d\x3b\xb7\xd0\xea\x08\x63\x6b\x45\xc9\x5a\x37\xa4\xff\x06\xa7\x7b\xb5\x5b\xd9\xd6\x1d\xee\x56\xb8\xa8\x75\x1b\xac\xd7\xc5\x1b\x00\x68\xdf\x45\x08\x21\xb2\x17\xab\x8b\xbf\xf9\x8a\x93\x01\x24\xb5\x1a\x64\x27\xfd\x95\xc4\x23\x17\xdd\x3e\x56\xf9\x98\x43\x95\x81\xd0\xea\x2c\xbf\x5b\xdb\x0b\xf0\xeb\x05\x4a\x78\xed\xba\xb7\x3d\xba\x98\x10\xce\x0e\xf1\x3c\x75\x01\x3c\x8b\xab\xd2\xae\x1f\xd7\x7d\xe7\x27\x00\x00\xff\xff\x73\x2c\xed\xe3\xa3\x01\x00\x00")

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

	info := bindataFileInfo{name: "config.yaml", size: 419, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _datastoreYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x8f\x31\x4f\x03\x31\x0c\x85\xf7\xfc\x8a\xa7\xae\x88\x46\x95\x18\x90\x57\x40\xec\x0c\x2c\x88\xc1\x4d\xac\xbb\xa8\x77\x49\x88\x7d\x05\xfe\x3d\xca\x71\x55\x61\xb3\xec\xe7\xf7\xbe\xc7\x35\xbd\x4a\xd3\x54\x32\x61\x16\xe3\x3d\xd7\xaa\xa1\x44\xd9\x87\x32\xfb\xf3\x81\xa7\x3a\xf2\xc1\x9d\x52\x8e\x84\x17\xd1\xb2\xb4\x20\x0f\x13\xab\xba\x2e\x8f\x6c\x4c\x0e\xc8\x3c\x0b\xe1\x91\x8d\xd5\x4a\x13\xa7\x55\x42\xdf\x73\x4d\xcf\xad\x2c\x95\x70\x5a\x8e\x12\x8f\xdd\xd5\x01\x9f\x92\x86\xd1\x08\x77\x0e\x48\xa1\x64\xed\xda\x5b\x68\x0b\x84\xd1\xac\x2a\x79\x1f\x62\xfe\x0f\x73\xba\x57\xbf\x8a\x7d\xdb\x38\x42\xe7\x10\xf5\xf1\x92\xbb\xd7\xf3\xe0\x00\xc0\xbe\xab\x10\xd2\xcc\x83\x78\x3d\x0f\x37\x5f\xf3\xe4\x00\xc9\xd6\x92\x6c\x61\xbf\xc8\x4f\x13\xab\xa5\xa0\xc2\x2d\x8c\xeb\x6b\x93\x8f\x25\x35\x89\x04\x6b\x8b\x5c\xdd\xd6\x09\x18\x7a\x1d\x25\xbc\xed\xae\x8d\x76\xef\xdb\xf1\x02\x46\x90\xbf\xbe\xa2\xee\x27\x00\x00\xff\xff\x1a\x82\xb0\x84\x68\x01\x00\x00")

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

	info := bindataFileInfo{name: "datastore.yaml", size: 360, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _helm2Yaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x8e\xb1\x4e\xc4\x30\x10\x44\x7b\x7f\xc5\xf4\x88\xb3\xee\x68\xd0\xb6\x34\xd4\x57\xd0\xaf\x9c\x21\xb1\xce\x76\x8c\xd7\x17\xe0\xef\x51\x42\x74\x02\x89\x8e\x76\xf4\x76\xdf\xd3\x1a\x5f\xd8\x2c\xce\x45\x90\xd9\xf5\xa0\xb5\x5a\x98\x07\x1e\xc2\x9c\xfd\x72\xd4\x54\x27\x3d\xba\x4b\x2c\x83\xe0\x4c\x9b\xaf\x2d\xf0\x29\xa9\x99\x5b\xf1\x41\xbb\x8a\x03\x8a\x66\x0a\x9e\x99\x32\x4e\xce\x2a\xc3\x3a\xbe\x33\x8e\x53\x17\x3c\x38\x20\x86\xb9\xd8\x3a\xde\xc3\x5a\x10\x4c\xbd\x57\x13\xef\xc3\x50\x7e\x2b\x2f\x8f\xe6\x37\xd8\xb7\xdd\x16\x56\x1b\xcd\x4f\x4c\xf9\x60\xcb\xe8\x00\xa0\x7f\x56\x0a\x62\xd6\x91\xde\x96\xf1\xee\x23\x27\x07\xb0\xf4\x16\xb9\x7b\xbe\x9b\xce\x4c\x54\xa3\x6d\x57\x8d\x6f\xd7\xd8\x38\x08\x5e\x35\x19\xb7\xad\x6a\x9f\x04\xdb\x77\xbf\x9c\x7c\xfb\xc9\xdf\xaa\xff\xd3\xbd\x7f\xbc\xa5\xff\x1d\xff\x15\x00\x00\xff\xff\x89\x40\x4e\x66\x89\x01\x00\x00")

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

	info := bindataFileInfo{name: "helm2.yaml", size: 393, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _helm3Yaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x8e\x3d\x4f\xc4\x30\x10\x44\x7b\xff\x8a\xe9\x11\xb1\x4e\x69\x90\x5b\x1a\xea\x2b\xe8\x57\xce\x90\x58\xe7\x2f\xbc\xbe\x00\xff\x1e\x25\x44\x27\x90\xe8\xae\x1d\xbd\xdd\xf7\xa4\x86\x57\x36\x0d\x25\x3b\x24\x76\x19\xa4\x56\xf5\x65\xe2\xe0\x4b\xb2\xeb\x49\x62\x5d\xe4\x64\x2e\x21\x4f\x0e\x67\x6a\xb9\x36\xcf\xe7\x28\xaa\x66\xc3\x27\xe9\xe2\x0c\x90\x25\xd1\xe1\x85\x31\x61\x34\x5a\xe9\xb7\xf1\x83\x61\x5e\xba\xc3\x68\x80\xe0\x4b\xd6\x6d\x7c\x84\x36\xef\xb0\xf4\x5e\xd5\x59\xeb\xa7\xfc\x57\x79\x79\x52\xbb\xc3\xb6\x1d\x36\xbf\xd9\xa8\x76\x61\x4c\x83\xae\xb3\x01\x80\xfe\x55\xe9\x10\x92\xcc\xb4\xba\xce\x0f\x9f\x29\x1a\x80\xb9\xb7\xc0\xc3\xf3\xd3\x74\x66\xa4\x28\x75\xbf\x6a\x7c\xbf\x86\xc6\xc9\xe1\x4d\xa2\x72\xdf\xaa\xf4\xc5\x61\xff\x6e\xd7\xd1\xb6\xdf\xfc\xad\xfa\x9e\xee\xe3\xe3\x2d\xfd\xff\xf8\xef\x00\x00\x00\xff\xff\x80\x70\x8d\x5b\x89\x01\x00\x00")

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

	info := bindataFileInfo{name: "helm3.yaml", size: 393, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _kubernetesYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x8e\xb1\x4e\xec\x30\x10\x45\x7b\x7f\xc5\x68\xdb\xf7\x36\xd6\x76\xab\x29\xa1\x44\x34\x20\xd1\x20\x8a\x59\x67\x94\x58\x49\x6c\xe3\x19\x87\xe5\xef\x91\x43\x36\x08\xa4\xad\x90\xe8\xac\xeb\x7b\xe7\x1c\x4a\xfe\x89\xb3\xf8\x18\x10\x26\x56\x6a\x28\x25\x71\xb1\xe5\xc6\xc5\xc9\xce\x07\x1a\x53\x4f\x07\x33\xf8\xd0\x22\x3c\xb0\xc4\x92\x1d\xdf\x8e\x24\x62\x6a\xbd\x25\x25\x34\x00\x81\x26\x46\xb8\x2b\x27\xce\x81\x95\xc5\x48\x62\x57\x3f\xde\xd8\x77\xbd\x22\x1c\x0c\x80\x77\x31\x48\x0d\xf7\x20\xd9\x21\xf4\xaa\x49\xd0\x5a\xd7\x86\xef\xd8\xe1\x28\x76\x29\xdb\xbc\x12\x5d\x25\xb2\xd8\x61\x23\x34\x32\x77\x06\x00\x40\xdf\x13\x23\xf8\x89\x3a\xb6\x32\x77\xff\xce\xd3\x68\x00\x38\x68\xf6\xbc\xd2\x3e\xed\x6e\x48\xbc\x5b\x26\x99\x5f\x8b\xcf\xdc\x22\x68\x2e\xbc\x44\x89\xb4\x47\xb0\xcb\x7b\xf3\xfc\x8d\xe9\xa9\xd2\x36\xc9\x6b\x9a\x17\xb7\x7b\x72\xbd\x0f\x2c\xd7\xf4\x96\xf5\x7a\xa9\xcb\xb1\x24\x41\x78\xde\xb9\xb1\x88\x72\x6e\xce\xfb\xe1\x28\x8d\x8f\xbb\xff\xb0\x65\x6b\xf2\xb2\x8e\x2e\x7a\x08\xd3\x17\xea\x07\x1d\x1e\x59\xff\xd2\xa0\xd2\x3e\x02\x00\x00\xff\xff\x61\xc5\x80\x8f\x80\x02\x00\x00")

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

	info := bindataFileInfo{name: "kubernetes.yaml", size: 640, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _monitoringYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8f\x3f\x4b\x04\x31\x10\xc5\xfb\x7c\x8a\xc7\xb5\xe2\x86\xeb\x8e\x69\x2d\xac\x04\xb1\xb0\x11\x8b\x90\x1d\xb2\xe1\x36\x7f\x9c\xc9\xae\xfa\xed\x25\xeb\x7a\x2a\x5c\x37\x30\xef\xcf\xef\xb9\x1a\x9f\x59\x34\x96\x4c\x48\xdc\xdc\xe0\x6a\x55\x5f\x46\x1e\x7c\x49\x76\x3d\xba\xb9\x4e\xee\x68\xce\x31\x8f\x84\x27\xd6\xb2\x88\xe7\xbb\xd9\xa9\x9a\x2e\x1f\x5d\x73\x64\x80\xec\x12\x13\x1e\x4a\x8e\xad\x48\xcc\xc1\x68\x65\xdf\x1f\xae\xc6\x7b\x29\x4b\x25\xa4\xcb\x73\xf0\x45\xb8\x68\x6f\x30\xc0\x3b\xc7\x30\x35\xc2\xc9\x00\xd1\x97\xac\xdd\x76\x0b\x15\x4f\x98\x5a\xab\x4a\xd6\xfa\x31\xff\x07\x3b\x9f\xd4\x6e\x62\x2b\x3b\x93\xef\x4c\xac\xf6\x4f\x8d\xae\xc1\x00\x40\xfb\xac\x4c\x88\xc9\x05\xb6\xba\x86\x9b\x8f\x34\x1b\x80\x73\x93\xc8\x7b\xdb\x37\xff\xa3\x94\xc4\x6d\xe2\x45\x37\x9f\xf0\xdb\x12\x85\x47\x42\x93\x85\x7f\xa3\xb6\x0b\x08\x7d\x97\x12\x5e\x0e\x57\xa7\x1d\x5e\x77\xdd\x0f\x21\xa1\x5e\xf2\x59\xcd\x57\x00\x00\x00\xff\xff\x44\x39\xbd\x3e\x7a\x01\x00\x00")

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

	info := bindataFileInfo{name: "monitoring.yaml", size: 378, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _networkingYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x91\x3d\x4f\xc3\x30\x10\x86\x77\xff\x8a\x53\x06\x16\x20\x56\x07\xa4\xca\x2b\x2c\x2c\x08\x81\xd4\x05\x31\x58\xce\xc9\x3d\x25\xb1\xcd\x9d\x93\xb6\xff\x1e\xb9\x49\xd4\x56\x30\xf0\xb1\x59\xd6\x7b\xef\xf3\x9c\xce\x26\xda\x20\x0b\xc5\x60\xa0\xc7\x6c\x6b\x9b\x92\xb8\xd8\x60\xed\x62\xaf\xc7\x95\xed\xd2\xd6\xae\x54\x4b\xa1\x31\xf0\x82\x12\x07\x76\x78\xdf\x59\x11\x55\xe2\x8d\xcd\xd6\x28\x80\x60\x7b\x34\x50\xbd\x22\x8f\xe4\x10\xae\xe0\x81\xc4\xc5\x11\xf9\x50\x29\x49\xe8\x4a\x66\x87\xe4\xb7\xd9\xc0\x9d\x02\x20\x17\x83\x94\xcf\x5b\x10\x76\x06\xb6\x39\x27\x31\x5a\xbb\x26\x5c\x1a\xb4\x6b\xd1\xc7\xb0\xe6\x19\xee\x0a\x1c\x45\x07\xcc\xbb\xc8\x2d\x05\x5f\xcb\xe8\x15\x00\x40\x3e\x24\x34\x40\xbd\xf5\xa8\x65\xf4\xd7\xfb\xbe\x53\x00\x18\x32\x13\xce\xb4\x49\x74\xf6\x94\xe3\x14\xe3\xc7\x40\x8c\x8d\x81\xcc\x03\x9e\x8a\x8e\x2f\x00\xcf\x71\x48\x62\xe0\xad\xaa\xde\xe7\xaf\x45\xc5\x80\x9c\x8a\x96\xee\xc7\xe0\x19\x8b\xe1\xef\xca\xcf\xd6\x69\xd7\x52\x53\xac\x6e\xa0\xc2\x7d\xc6\x50\x8e\x23\xdf\xb0\xe9\x0c\xb4\xc0\x37\xf1\x60\x3d\xf2\x5f\x25\xc6\x69\xfc\xe2\x04\x3f\x24\x3f\x4d\xfe\xf0\x1c\x3b\x72\xf4\xff\xed\xbf\x52\xe7\x4c\x5a\x00\x9f\x01\x00\x00\xff\xff\x2c\x2f\xed\xf4\xbb\x02\x00\x00")

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

	info := bindataFileInfo{name: "networking.yaml", size: 699, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _securityYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x90\xc1\x4e\xc3\x30\x10\x44\xef\xfe\x8a\x55\xaf\xa8\xb5\x7a\x03\xdf\x50\x7f\x00\x81\xc4\x05\x71\x58\x6d\x96\x74\xd5\xc4\x36\xbb\x4e\x20\x7f\x8f\xe2\xa6\x40\xa4\x5e\x7a\xb3\xc6\xb3\xf3\x46\x83\x59\x5e\x59\x4d\x52\x0c\xd0\x73\xc1\x1d\xe6\x6c\x94\x1a\xde\x51\xea\xfd\xb8\xc7\x2e\x1f\x71\xef\x4e\x12\x9b\x00\xcf\x6c\x69\x50\xe2\x43\x87\x66\x6e\xb6\x37\x58\x30\x38\x80\x88\x3d\x07\x78\x61\x1a\x54\xca\xe4\x2c\x33\xcd\xf2\x17\x4b\x7b\x2c\x01\x1e\x1c\x80\x50\x8a\x36\x8b\x5b\x30\xa5\x00\xc7\x52\xb2\x05\xef\xa9\x89\x6b\xe8\xe9\xde\x7c\x35\x7b\x5d\x78\x34\xf3\xd8\xbc\x2d\xf9\x3b\x1b\x5b\x07\x00\x50\xa6\xcc\x01\xa4\xc7\x96\xbd\x8d\xed\xdd\x77\xdf\x39\x00\x8e\x45\x85\x17\xd6\xb9\xd9\x53\xea\x84\xa6\x7a\xa3\xfc\x39\x88\x72\x13\xa0\xe8\xc0\x7f\x31\xf5\x05\xd0\x6a\x1a\xb2\x05\x78\xdb\xe4\x7a\xb3\x79\x5f\x3e\x2e\x65\x02\xe4\xd4\x5c\x9a\x54\x8f\xb0\xfd\x43\x1d\x58\x8b\x7c\x08\x61\xa9\xf2\x0d\xc0\x31\x4d\xd8\xb2\xae\xd6\xb8\x82\xa7\x35\x60\xfb\x3b\xbe\x8e\x42\x0c\x8f\x44\x69\x88\xe5\x46\xf6\x15\x8e\x9d\x03\xf1\x92\xf7\x13\x00\x00\xff\xff\x23\xdb\x8d\xb7\x2b\x02\x00\x00")

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

	info := bindataFileInfo{name: "security.yaml", size: 555, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _storageYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x91\xcb\x4a\x34\x31\x10\x85\xf7\x79\x8a\xc3\x6c\x7f\xfe\x0e\xb3\x72\xc8\x76\x5e\x40\x14\x66\x23\x2e\x8a\x74\x91\x09\x93\x9b\xa9\x74\xab\x6f\x2f\xdd\x93\xf6\x02\xba\x10\xdc\x85\xc3\xe1\xfb\x2a\x55\x54\xfc\x89\xab\xf8\x9c\x0c\x22\x37\x1a\xa8\x14\xb1\x79\xe4\xc1\xe6\xa8\xe7\x3d\x85\x72\xa6\xbd\xba\xf8\x34\x1a\xdc\xb1\xe4\xa9\x5a\x3e\x06\x12\x51\x4b\x7d\xa4\x46\x46\x01\x89\x22\x1b\xdc\xb7\x5c\xc9\xb1\x92\xc2\x76\x49\x9f\xd9\xbb\x73\x33\xb8\x51\x80\xb7\x39\xc9\x12\xfe\x87\x54\x6b\x70\x6e\xad\x88\xd1\xda\x8e\xe9\xab\xf3\x72\x10\xbd\x96\x75\xed\x3a\xbb\xe8\x58\xb4\x5c\xf1\x83\xcc\x4e\x01\x40\x7b\x2d\x6c\xe0\x23\x39\xd6\x32\xbb\x7f\x2f\x31\x28\x80\x53\xab\x9e\xbb\xea\x3a\xd7\xed\xf2\x43\x69\x9c\x1a\x4e\x39\x4c\x91\x71\x0c\xe4\xa3\xac\x94\xca\x4f\x93\xaf\x3c\x1a\xb4\x3a\xf1\x07\x78\x7d\x01\xae\xe6\xa9\x88\xc1\xc3\x6e\xf7\xd8\xa3\x6d\x2e\x83\xf2\x0e\x9e\x57\xae\xdd\xb0\x3f\x98\x43\xe4\xbf\x97\x7e\xf6\xf5\x0b\xe0\x78\x5d\xd9\xef\x5c\xdb\x7e\x2f\x07\x19\x7c\xfe\xc6\xdc\x0b\xfd\x1e\xea\x2d\x00\x00\xff\xff\x84\xbb\xd9\x85\x3c\x02\x00\x00")

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

	info := bindataFileInfo{name: "storage.yaml", size: 572, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _workloadsYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\xd1\xb1\x6e\xc2\x30\x10\x80\xe1\x3d\x4f\x71\x62\xad\x4a\x44\xa7\xca\x2b\x4c\x9d\x2a\x90\xda\xa1\xea\x70\xd8\x47\xe2\xe2\xf8\xdc\x3b\x07\xca\xdb\x57\x09\x89\x2a\x10\x1d\x68\xd9\xa2\x93\xf3\x7f\x67\x19\x93\x7f\x21\x51\xcf\xd1\x40\x43\x19\xa7\x98\x92\x5a\x76\x34\xb5\xdc\x94\xbb\x19\x86\x54\xe3\xac\xd8\xfa\xe8\x0c\x2c\x49\xb9\x15\x4b\xf3\x80\xaa\x45\x77\xdc\x61\x46\x53\x00\x44\x6c\xc8\xc0\x2b\xcb\x36\x30\x3a\x2d\x34\x91\xed\xe6\x7b\xf2\x55\x9d\x0d\x3c\x14\x00\xde\x72\xd4\x6e\x78\x0f\x2a\xd6\x40\x9d\x73\x52\x53\x96\xd6\xc5\x53\x75\xfb\xa8\x65\x7f\xb8\x94\x01\xb4\x1d\x48\x5a\xee\x47\x60\xaa\xbb\xaa\x00\x00\xc8\x87\x44\x06\x7c\x83\x15\x95\xba\xab\xee\xbe\x9a\x50\x00\x50\xcc\xe2\x69\xc0\x8e\xbb\x2d\x28\x05\x3e\x34\x14\xb3\xf6\x3f\x0a\x7d\xb6\x5e\xc8\x19\xc8\xd2\xd2\x4f\xab\xff\x02\xa8\x84\xdb\xa4\x06\xde\x26\xdd\x6a\x93\xf7\x61\x3c\x2e\x64\xc0\x9d\xf4\x46\x65\x49\x29\x78\x8b\xb0\xa2\x1b\x31\x72\x0c\x2a\x5d\x62\xb2\xe7\x08\x73\x8e\x59\x38\x04\x92\x2b\xc5\xdf\xb5\xae\x6b\x4f\xb2\x23\xbc\xca\x98\x69\xd3\x86\xdb\x5d\x50\x87\xe2\xd9\x0d\x17\x48\x0d\xc7\xdb\x31\xae\xef\x9d\x21\x4f\xbc\xbe\xb2\xbe\xc6\x6c\xeb\x0b\xf9\x8f\x63\x69\x0c\xcf\x85\xe3\xa5\xfa\x06\x83\xfe\x25\x6f\x85\xe3\x19\xf1\xcc\xee\xff\xef\x9d\xba\xc8\x77\x00\x00\x00\xff\xff\xec\x1a\xa7\x35\x03\x04\x00\x00")

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

	info := bindataFileInfo{name: "workloads.yaml", size: 1027, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
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
