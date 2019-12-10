// Package crds Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// meta.appscode.com_resourceclasses.yaml
// meta.appscode.com_resourcedescriptors.yaml
package crds

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

var _metaAppscodeCom_resourceclassesYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x58\x5f\x6f\xdb\x38\x12\x7f\xf7\xa7\x18\xf4\x1e\xda\x02\xb1\x8c\xb4\x28\xee\xce\x6f\x85\xdb\x1e\x72\xdb\x76\x8b\x24\xed\xfb\x98\x1c\xcb\xb3\xa1\x48\x2d\xff\xb8\x49\x16\xfb\xdd\x17\x43\x49\xb6\x6c\xcb\x76\x6a\x6c\xb1\x7c\xd3\xcc\x70\xf8\xe3\x6f\xfe\x90\xd4\x68\x3c\x1e\x8f\xb0\xe6\x6f\xe4\x03\x3b\x3b\x05\xac\x99\xee\x23\x59\xf9\x0a\xc5\xdd\x7f\x42\xc1\x6e\xb2\xba\x9c\x53\xc4\xcb\xd1\x1d\x5b\x3d\x85\x59\x0a\xd1\x55\xd7\x14\x5c\xf2\x8a\xde\xd1\x82\x2d\x47\x76\x76\x54\x51\x44\x8d\x11\xa7\x23\x00\xe5\x09\x45\x78\xcb\x15\x85\x88\x55\x3d\x05\x9b\x8c\x19\x01\x58\xac\x68\x0a\xbe\x9d\xae\x0c\x86\x40\xa1\x90\xb9\x05\xd6\x75\x50\x4e\x53\xa1\x5c\x35\x0a\x35\x29\xf1\x84\x5a\x67\xf7\x68\xbe\x78\xb6\x91\xfc\xcc\x99\x54\xd9\x20\xba\x31\xfc\xff\xe6\xd7\xcf\x5f\x30\x2e\xa7\x50\x84\x88\x31\x85\xa2\x5e\x62\xa0\x11\x40\xb7\xd2\x4d\x16\x67\x41\x7c\xa8\x69\x0a\x21\x7a\xb6\xe5\xee\xec\x0e\x7c\xb1\x87\xbc\xe7\xeb\x6d\x49\x3d\x47\x1a\xa3\x7c\x96\xde\xa5\x7a\x0a\xfb\x3b\x68\x66\x65\xa0\x00\x0d\x77\x1d\x6b\x33\xd9\x76\x96\x1b\x0e\xf1\x97\x7d\xdd\x47\x0e\x31\xeb\x6b\x93\x3c\x9a\x3d\xc2\xb2\x2e\xb0\x2d\x93\x41\xbf\xa3\x1d\x01\xd4\x9e\x02\xf9\x15\x7d\xb5\x77\xd6\x7d\xb7\x1f\x98\x8c\x0e\x53\x58\xa0\xc9\xdc\x04\xe5\x64\x03\x9f\x05\x5e\x8d\x8a\xb4\xc8\xd2\xbc\xf3\xd2\x42\x6e\x08\x9d\xc2\x1f\x7f\x8e\x00\x56\x68\x58\x67\x62\x1a\xa5\xab\xc9\xbe\xfd\x72\xf5\xed\xf5\x8d\x5a\x52\x85\x8d\x50\x16\x76\x35\xf9\xc8\x9d\x0f\x19\xbd\xfc\x5a\xcb\x00\x34\x05\xe5\xb9\xce\x1e\xe1\xb9\xb8\x6a\x6c\x40\x4b\x46\x51\x80\xb8\x24\x58\x35\x32\xd2\x10\xf2\x32\xe0\x16\x10\x97\x1c\xc0\x53\xde\xa2\x8d\x19\x52\xcf\x2d\x88\x09\x5a\x70\xf3\xdf\x48\xc5\x02\x6e\x84\x06\x1f\x20\x2c\x5d\x32\x1a\x94\xb3\x2b\xf2\x11\x3c\x29\x57\x5a\x7e\x5c\x7b\x0e\x10\x5d\x5e\xd2\x60\xa4\x96\xfb\x6e\xe4\xb4\xb3\x68\x84\x84\x44\x17\x80\x56\x43\x85\x0f\xe0\x49\xd6\x80\x64\x7b\xde\xb2\x49\x28\xe0\x93\xf3\x04\x6c\x17\x6e\x0a\xcb\x18\xeb\x30\x9d\x4c\x4a\x8e\x5d\x45\x29\x57\x55\xc9\x72\x7c\x98\x28\x67\xa3\xe7\x79\x8a\xce\x87\x89\xa6\x15\x99\x49\xe0\x72\x8c\x5e\x2d\x39\x92\x8a\xc9\xd3\x04\x6b\x1e\x67\xe0\x36\xe6\xb2\xac\xf4\xbf\xd6\xa1\x7a\xde\x43\xba\x93\xdf\xcd\xc8\x89\x77\x90\x77\x49\x3d\xe0\x00\xd8\x4e\x6b\xf0\x6f\xe8\x15\x91\xb0\x72\xfd\xfe\xe6\x76\x9d\x65\x39\x04\xdb\x9c\x67\xb6\x37\xd3\xc2\x86\x78\x21\x8a\xed\x82\x7c\x13\xb8\x85\x77\x55\xf6\x48\x56\xd7\x8e\x6d\xcc\x1f\xca\x30\xd9\x6d\xd2\x43\x9a\x57\x1c\x25\xd2\xbf\x27\x0a\x51\xe2\x53\xc0\x0c\xad\x75\x11\xe6\x04\xa9\x96\xf2\xd3\x05\x5c\x59\x98\x61\x45\x66\x86\x81\x7e\x3a\xed\xc2\x70\x18\x0b\xa5\xa7\x89\xef\xb7\xc3\x6d\xc3\x86\xad\xb5\xb8\xeb\x74\xdd\x18\xaa\xa1\xb6\x8e\xfe\x97\x7b\xcd\x96\xf4\xc0\xea\x32\x48\x73\xc4\xb9\xa1\xe1\x09\x73\xe7\x0c\xe1\x76\xed\xb0\x72\x36\xec\x9a\x6f\x65\xcc\x95\x58\xe4\x8c\xb1\xe0\xea\xa6\x31\xe7\x26\x26\x75\x97\xa7\xc3\xc2\x79\xd1\x62\x5d\x1b\x56\xb9\x3e\x8b\x1d\x8f\x90\xdd\xe4\x38\xf9\x2a\x5b\x00\x5b\x65\x92\x6e\xcb\xbe\x49\xb3\x0b\x08\xfc\xd8\x95\x1b\x57\x94\x71\xef\xba\xe2\x48\xd5\x1e\xe2\x5d\xcc\x15\x96\x74\x53\x93\x92\xf2\x8f\xc8\xb2\x81\xde\xd2\x38\x77\x29\x0a\x62\x16\x3b\x48\x81\x34\x60\xd8\x73\x09\xd9\x44\x0d\xed\xe6\x50\xc4\x9a\x21\xbb\x18\x92\xef\xa0\x7c\xd1\xd1\xf9\x12\x6e\x85\x03\x7e\xa4\xa6\xdd\x51\x0b\x8c\x2d\xd4\x7c\x4f\x26\xc0\x0b\x2a\xca\xe2\x62\xd0\x25\xc0\xab\x37\xf7\xaf\xde\xbc\xdc\x07\x09\xc7\x52\xa5\x45\xea\xd5\x13\x80\xde\xae\x23\x94\x43\xdd\x80\x5b\x97\x7e\x26\x0f\x88\xe3\x92\x9a\x34\x98\x07\x67\x52\xa4\x03\x68\xbf\x5e\x7f\xec\x3a\x6f\xe3\x48\x92\x07\xde\x61\xc4\xac\x6a\x03\xd6\xb5\xa1\x6c\x52\xac\xd5\x43\x41\x92\x81\x9e\xda\x63\x44\x0b\x6b\xd7\x1f\x66\xf0\xea\xf5\x7f\xff\x7d\x16\x27\x59\xfd\xc3\xd1\xab\xd8\x36\xf9\xba\x1d\xc2\x26\x70\xf0\x2c\x7f\x4d\x6a\x5b\x3e\x3b\x27\x50\xd2\x12\xd9\x93\xde\x87\x35\x96\x10\xee\x49\x07\xbb\x4e\x5f\x85\xde\xe3\xc3\x96\xc6\xb0\xbd\x3b\xde\x08\x3e\x8a\x45\x66\x1a\xd7\xf5\xbf\x36\x58\x51\x0e\x4f\x3e\x39\xad\x26\x2d\x21\x9e\x37\xa5\xb5\x8f\xce\x41\x48\x7e\x81\x8a\x7a\x77\x3d\xd0\x4e\xa5\xaa\x3b\xdf\x2f\x40\x63\x58\xce\x1d\x7a\x1d\x2e\x80\xa2\x3a\xa7\x0b\x08\xe0\x13\x0d\xa0\x4d\xc6\x0e\xce\x36\x86\x81\x38\x9d\x40\x75\xaa\x31\xf4\xe1\x9d\xce\xb0\x77\x9b\x0f\xe9\xbf\xcb\x54\xa1\x05\x4f\xa8\xa5\xc5\xe7\x8d\x91\x8d\x40\xf7\xb5\x69\x0a\xe6\x40\x6d\x48\x32\xd6\xc9\xd7\x2e\xac\x73\x53\x82\x7d\x56\x6d\x24\x6f\x9e\x00\xfc\xab\x37\xe2\x88\x15\x1a\xf3\x00\xf9\xd0\x0f\x80\x11\x10\xbe\xd3\x3c\x70\xcc\x71\xf7\x14\xc2\x19\x18\xce\x48\xed\x0a\x39\xe7\x00\xf9\xe3\x09\xfe\x69\x63\x77\xe8\xbc\xeb\xb9\xea\xb8\x3c\x7a\xea\xe5\xce\xd0\x77\x6b\x9b\x4b\x51\xf6\xd7\x29\xb2\x9b\xcd\x29\x08\xf2\x9e\xb8\x68\xfa\x47\xd8\x4f\x42\x39\x1c\x6b\x54\x77\xd2\x5d\xa4\x1b\x9f\xc0\xf0\x94\x42\x99\x49\x8d\xa8\x98\x9b\xec\x89\x03\xd3\x6a\x5e\xb1\x4e\x68\x06\x22\xe7\x3c\x38\x5f\xa2\xe5\xc7\x03\x6c\x1c\x2f\x0e\xaa\x90\x9f\x92\x5d\xef\xc5\x4e\x02\x94\x2f\x95\xf9\xe3\xfc\x74\x6a\x1f\x7a\xa7\x57\x95\xa7\x53\xb7\x68\xbf\xeb\xc9\xf4\x9f\x5c\x4a\x2a\x3f\x64\x36\x05\x35\xa7\x7f\xaa\x92\x6a\x8c\xcb\xbd\xd8\x1d\xc8\xb0\x23\x6b\x1f\x72\xbf\xf3\x1c\x3d\xb9\xc4\x16\x53\xf9\xa6\xdc\x3e\x29\xbb\x97\x35\x24\x8b\xd5\x9c\xcb\xe4\x52\x30\x0f\xc0\x5a\x6e\xf6\x0b\x26\x79\xfe\x74\x8b\x15\x00\x57\xbb\xfb\xcf\xce\x1d\x05\xfb\x5c\xb2\xde\xd9\x87\xaa\x75\xd0\xdc\x59\xb7\xd6\x92\xf3\x03\x57\x8e\x35\x60\x8a\x4e\x0a\x46\x6e\x9d\xe4\x95\x94\xc0\x09\xdf\x29\x48\x28\x87\xbd\xa9\xfc\xd7\x05\x2a\xf4\x61\x89\xc6\x0c\x45\xf0\x78\x3d\x95\x43\x4f\x87\x93\xb1\xe9\xc7\xe1\xac\xc9\xab\xfd\x97\xff\x13\xe7\x1e\xbb\xe5\xe4\xcd\x0c\xc8\x3b\xa8\x03\xaa\x16\xc8\xdf\x91\xf7\x2b\x0e\xfc\x43\xcf\xaa\xef\xc4\xe5\x32\x0e\xdb\xcb\x15\xa9\x24\x3f\x3a\xbe\xf1\xf1\xfa\xf5\xb7\x25\xec\x5e\x78\xa3\x21\x12\xc2\x96\xb4\x05\xbd\x25\x6b\x70\x1d\x7f\x9f\xee\x88\xba\x78\xc2\xea\x12\x4d\xbd\xc4\xcb\x8d\xac\xfd\x25\xd7\xfc\x2a\xeb\xa9\x01\xf2\xbf\x28\x3d\x85\xe8\x13\xb5\xbf\x96\x9c\xc7\x92\x5a\xc9\x5f\x01\x00\x00\xff\xff\xba\x76\xbb\x6d\x86\x14\x00\x00")

func metaAppscodeCom_resourceclassesYamlBytes() ([]byte, error) {
	return bindataRead(
		_metaAppscodeCom_resourceclassesYaml,
		"meta.appscode.com_resourceclasses.yaml",
	)
}

func metaAppscodeCom_resourceclassesYaml() (*asset, error) {
	bytes, err := metaAppscodeCom_resourceclassesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "meta.appscode.com_resourceclasses.yaml", size: 5254, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _metaAppscodeCom_resourcedescriptorsYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x5c\x5f\x73\xdb\x36\x12\x7f\xf7\xa7\xd8\xf1\x3d\xa4\x9d\x91\xe8\x8b\x3b\x9d\xbb\xd3\x9b\xc7\x49\x6f\x7c\x97\xa6\x99\xd8\xcd\xcb\xb5\x0f\x10\xb9\x12\x51\x93\x00\x0b\x80\xb2\xd5\x4e\xbf\xfb\xcd\x2e\x08\x8a\x92\x08\x92\x52\x92\x49\x8b\x27\x0b\x04\x16\x8b\xc5\xfe\xf9\x01\x0b\xf8\x62\x3e\x9f\x5f\x88\x4a\x7e\x40\x63\xa5\x56\x0b\x10\x95\xc4\x67\x87\x8a\x7e\xd9\xe4\xf1\x9f\x36\x91\xfa\x6a\xf3\x72\x89\x4e\xbc\xbc\x78\x94\x2a\x5b\xc0\x6d\x6d\x9d\x2e\xdf\xa3\xd5\xb5\x49\xf1\x15\xae\xa4\x92\x4e\x6a\x75\x51\xa2\x13\x99\x70\x62\x71\x01\x90\x1a\x14\x54\xf9\x20\x4b\xb4\x4e\x94\xd5\x02\x54\x5d\x14\x17\x00\x4a\x94\xb8\x00\xd3\x74\xcf\xd0\xa6\x46\x56\x4e\x1b\x9b\x50\xff\x44\x54\x95\x4d\x75\x86\x49\xaa\xcb\x0b\x5b\x61\x4a\xd4\x44\x96\xf1\x10\xa2\x78\x67\xa4\x72\x68\x6e\x75\x51\x97\xca\xd2\xb7\x39\xfc\xe7\xfe\x87\xb7\xef\x84\xcb\x17\x90\x58\x27\x5c\x6d\x93\x2a\x17\x16\x2f\x00\xc2\x68\xf7\x5c\xcd\x15\x6e\x5b\xe1\x02\xac\x33\x52\xad\x0f\x7b\x87\x09\x24\x47\xdc\x77\x68\xdd\xac\xb1\x43\x28\x13\x8e\x7e\xae\x8d\xae\xab\x05\x1c\xcf\xc0\xf7\x62\x46\x01\xbc\xfc\x76\x92\x0b\x53\xe7\x8f\x85\xb4\xee\xbf\x91\x06\x6f\xa4\x75\xdc\xa8\x2a\x6a\x23\x8a\x5e\xf1\xf1\x77\x9b\x6b\xe3\xde\xee\x46\x9c\x83\xc9\xfc\x07\xa9\xd6\x75\x21\x4c\x5f\xd7\x0b\x80\xca\xa0\x45\xb3\xc1\x1f\xd5\xa3\xd2\x4f\xea\x3b\x89\x45\x66\x17\xb0\x12\x05\xcb\xd1\xa6\x9a\x26\xcb\x84\x2b\x91\x22\xd1\xb4\xf5\x32\x90\x6a\x06\xf3\xc2\x5f\xc0\xef\x7f\x5c\x00\x6c\x44\x21\x33\x16\xa2\xff\xa8\x2b\x54\x37\xef\xee\x3e\x7c\x73\x9f\xe6\x58\x0a\x5f\x49\x03\xeb\x0a\x8d\x93\x81\x06\x95\x8e\x3e\xb6\x75\x00\x81\x5f\x56\xd3\x17\x44\xca\xb7\x81\x8c\x34\x10\x2d\xb8\x1c\x61\xe3\xeb\x30\x03\xcb\xc3\x80\x5e\x81\xcb\xa5\x05\x83\x3c\x45\xe5\x98\xa5\x0e\x59\xa0\x26\x42\x81\x5e\xfe\x82\xa9\x4b\xe0\x9e\xc4\x60\x2c\x49\xb2\x2e\x32\x48\xb5\xda\xa0\x71\x60\x30\xd5\x6b\x25\x7f\x6b\x29\x5b\x70\x9a\x87\x2c\x84\xc3\x66\x75\x42\x61\x15\x55\xa2\x20\x21\xd4\x38\x03\xa1\x32\x28\xc5\x16\x0c\xd2\x18\x50\xab\x0e\x35\x6e\x62\x13\xf8\x5e\x1b\x04\xa9\x56\x7a\x01\xb9\x73\x95\x5d\x5c\x5d\xad\xa5\x0b\x16\x98\xea\xb2\xac\x95\x74\xdb\xab\x54\x2b\x67\xe4\xb2\xa6\x25\xbf\xca\x70\x83\xc5\x95\x95\xeb\xb9\x30\x69\x2e\x1d\xa6\xae\x36\x78\x25\x2a\x39\x67\xc6\x95\x63\x33\x2e\xb3\xbf\xb5\x4b\xf5\xa2\xc3\xe9\x81\x2d\xf8\xc2\x4a\x1a\x95\x3b\x69\x28\x48\x0b\xa2\xe9\xe6\xf9\xdf\x89\x97\xaa\x48\x2a\xef\x5f\xdf\x3f\xb4\xaa\xc6\x4b\xb0\x2f\x73\x96\xf6\xae\x9b\xdd\x09\x9e\x04\x25\xd5\x0a\x8d\x5f\xb8\x95\xd1\x25\x53\x44\x95\x55\x5a\x2a\xc7\x3f\xd2\x42\xa2\xda\x17\xba\xad\x97\xa5\x74\xb4\xd2\xbf\xd6\x68\x1d\xad\x4f\x02\xb7\x42\x29\xed\x60\x89\x50\x57\x64\xaa\x59\x02\x77\x0a\x6e\x45\x89\xc5\xad\xb0\xf8\xd9\xc5\x4e\x12\xb6\x73\x12\xe9\xb8\xe0\xbb\xee\x73\xbf\xa1\x97\x56\x5b\x1d\xbc\x62\x28\x7d\x36\x44\x25\xdd\x39\xc9\x3d\xed\x74\x58\x1e\x55\x1e\xac\x73\x70\x41\xde\xcf\xee\x7c\x3c\x8f\x2d\x57\x12\x49\x03\x3c\x7d\x58\x69\x03\xec\x3c\xcc\x11\x4d\xf2\x3b\x19\x42\x45\x4e\x5b\xaa\x75\x72\xd4\x20\xc6\xfa\x31\x4b\x3d\x9f\x0f\x78\xee\xfc\xf0\x1a\x9a\xd7\xa5\x50\x60\x50\x64\x62\x59\xe0\xde\xf7\xc6\x2d\xf4\x12\x0d\x82\x3b\xe6\x16\x62\x4b\xb7\x2b\x2b\x6d\x4a\xe1\x26\x70\xeb\x1b\x32\xa3\x0a\x74\xe5\xa3\x1b\xfc\xe0\xbd\x24\x8f\xe2\x1d\x9b\x17\xfb\x4a\xf7\x09\x97\xf9\x21\x23\x69\x18\x86\x87\x1c\xe1\x05\x05\x9c\x17\x5d\xfa\x55\x55\x48\xcc\x82\xbf\xaa\x8c\x2c\x85\xd9\x82\xcc\x48\x4f\x57\xb2\x77\xd9\x76\x52\xa0\x6e\xc2\x5a\x69\x1d\x48\xd5\xd8\x9d\x0d\x9d\xb7\x64\xee\x4d\x3b\xe9\x3d\x70\xb0\xf9\x08\x51\x62\x8e\x6c\x1d\xbb\xe6\x96\xd7\x4b\x0a\x96\x57\x3f\xdc\xdc\x5d\x35\x12\x98\xdf\x7b\x45\x4b\xd9\x5f\x5f\x2d\x0b\xbd\xbc\x2a\x85\x75\x68\xae\x1a\x1f\x6f\xaf\xae\x93\xbf\x93\x8d\x91\xcd\xcc\xd9\xd0\x22\x63\x92\x82\x96\xda\xe0\x59\x0b\xfa\x8b\xd5\x8a\xe1\xc1\xf8\x92\x06\x24\xd1\xf8\x47\x59\x56\x05\x72\x25\x54\xc2\xe5\x33\x90\x09\x26\xf0\x24\x5d\x0e\xc2\x18\xb1\x05\xa5\x7d\x34\x3a\x8b\x2f\x06\x23\xe3\x3c\x51\xb3\x5e\x6b\xe0\x0f\x24\x19\x76\xa7\xe7\x2b\x7c\x65\xa4\x36\xd2\x6d\x27\xf0\x12\x9a\x36\x4a\x4f\x51\x72\x8d\xa6\x51\xf3\x26\x6e\x18\x2c\x84\x93\x1b\x04\x59\x56\xda\x38\xa1\xa2\x9a\x14\x22\x7b\xa3\x7e\xa9\x2e\x2b\x61\xbc\x9a\x6b\x97\xa3\xb1\x09\xbc\xd1\x4f\x68\x40\xd5\xe5\x92\x02\x8b\x30\x34\x51\x45\x1e\xc9\x60\x16\xa1\x9a\xcb\x75\x8e\xa6\x65\x35\x81\x06\x69\x82\xcb\x85\xe3\xd8\xb4\x44\xd0\xa5\x74\x0e\x33\x32\x88\x42\x96\x92\xfe\x64\x58\x14\xa1\x69\x53\x54\xc2\x48\xdd\x82\x8a\x25\xc2\x5a\x6e\x50\xd1\xaa\x1c\x8c\xd7\x4b\xa2\x71\x2a\x24\xb1\x6f\xae\x07\x16\xa9\x91\x68\x4f\x0b\xfe\x3e\xbe\x42\xec\x76\xfc\xea\x0c\x78\xa2\x3d\x9f\x13\x99\xf3\x5f\xc6\xc2\x09\x2e\x48\x83\xd9\xb1\x74\xe6\xad\xed\xf7\x7c\x22\xf3\xe9\xa9\x0e\xeb\xd8\xf3\x89\xb8\x38\xaa\xee\x8d\xef\xdd\x4f\xec\x29\x0e\xa2\xba\x52\x98\x32\xc2\x98\x18\xd9\x87\x83\x6c\x41\x60\xa6\x5f\x37\x26\xb8\xa0\x07\x2c\x2b\x42\xc0\x67\x13\x60\xc3\x89\x3b\xd8\x11\x0a\x06\x57\x68\x50\xa5\xfd\x53\x3b\xd0\xee\xcb\xf7\x6d\x6b\x76\x06\x82\x17\x98\x7c\xb3\xb7\x6e\x83\xae\x36\x8a\x7c\xe5\xed\xfd\x87\xc6\xe8\x5c\xd4\x53\xec\x86\xe6\xc8\x2a\xcc\x1a\x5d\x1b\xfd\x2c\xfc\xa4\xe0\x6e\x05\x28\xd2\x1c\x8c\x7e\x82\x5c\xf8\x98\xa0\xd6\x45\x70\xb7\xb3\x08\x61\xc9\x51\xbb\x21\xe8\xe3\xe5\xc3\x3e\x75\xfa\xae\xb4\x9a\xb7\xf2\xcb\x40\x1b\xa8\x6d\xd4\x3e\xc8\xb3\x5a\x72\xf8\x6d\x0f\x10\x3e\x5c\x37\x04\x03\xe5\x04\x5e\x3f\x0b\x8a\x5b\x0b\x50\x2f\x41\x5d\xc3\x4f\x2a\x42\xf2\x70\x72\xee\x49\x07\xc0\x39\xdb\x9f\xc2\xff\x68\xd0\x59\x3b\xf2\xcf\x31\x97\xd1\x19\x79\xa6\x2c\x8d\x3e\x53\xf6\xba\x4f\x90\x2e\x37\x88\xdd\xd1\x62\x92\x8c\xf1\x30\x23\x44\xfe\x73\x72\x38\xe4\xec\x31\x8c\x3a\x7b\xec\xf7\xb3\xd0\xc7\xce\x4a\xd7\x66\xda\xdc\x79\xdc\x99\xa8\xe4\xbf\x8d\xae\xab\xa9\x92\x98\x3d\xbe\x6c\xfb\x74\x18\x6c\xeb\xae\x2f\x7b\xe9\x44\xbc\x81\x2f\x23\x76\x15\x77\x3f\xbe\x58\x2c\x30\x75\xda\x4c\xb0\xba\x1b\x28\xc4\x12\x8b\xb6\x8b\x47\x23\xbe\xee\xd7\x1a\xcd\x16\xf4\x06\x0d\x19\x07\x3a\x0a\xea\xad\x09\xc5\xa4\xf3\xe0\x51\x66\x5d\x70\xf3\x52\xb8\x34\x7f\x43\xd4\x6c\xb3\xcf\x76\x69\xfe\xfa\x99\xf6\x95\x1c\x43\xd8\xd2\x6f\xde\xbe\xa2\xad\xdf\x4d\x4c\x99\xb1\xac\xdc\xf6\x90\x4f\xa6\x44\xae\xa2\x28\x1a\x0f\x6d\x13\xb8\xe1\x63\xac\x83\xa6\x11\xaa\x81\x80\xd2\x6d\xff\xde\x96\xc3\xfe\xb9\xa5\xd4\x99\x54\xac\xdd\x81\xe8\x8f\x64\xe1\x45\x4f\x58\x5e\xaf\xa6\xcd\x01\x76\x21\xb2\xf4\xfb\x74\x2f\xfe\x5d\x4d\x47\xc0\x51\x1a\x83\xaa\x78\xc4\xf6\x91\xc6\x74\x86\x6b\xb0\xf5\x38\xd3\xe0\x7d\x3a\xed\xdb\x85\x54\xb6\x39\x67\x99\x81\x80\x47\xdc\xfa\x23\x19\xde\x78\xa1\x11\x8e\x21\x0d\x07\x00\x3e\xcc\x19\xa1\x8a\x44\x81\x09\x34\x67\x37\x03\xed\xc7\x97\xd6\x97\x47\x8c\x00\xe8\x88\x88\x88\x83\x66\xc3\xe5\x65\x45\x15\x3c\x07\x76\xea\x53\xc4\x03\x7c\xd2\x46\x5b\x43\x3e\x2a\x19\x69\x3b\xea\x2f\x42\x09\x12\x3d\x69\x3a\xed\x32\xec\x0e\x84\xfc\x42\xbd\xb0\xcd\x8e\x40\x2b\x9b\xcb\x6a\x74\x42\xb4\x59\x0d\x8e\x24\x9c\xac\x7d\x10\x85\xcc\xda\x21\xbc\xbe\xde\xa9\x19\xbc\xd5\xee\x2e\x1a\x84\x77\xe5\xf5\xb3\xb4\xce\xfb\x96\x57\x1a\xed\x5b\xed\xb8\xe6\x93\x09\xcc\xb3\x79\x92\xb8\x7c\x97\x06\xa8\xfb\xbd\xa4\x5e\xed\x1d\xc8\xd9\x04\xee\x56\xe3\xd2\xca\x71\x27\x7a\x69\xe1\x4e\x11\x8e\xf0\x72\xf1\xc7\xa9\x7e\x20\x3f\x44\x59\xdb\x58\xa0\xdd\x95\x25\x32\x32\x61\x87\x4a\x3c\x1c\x8d\xd1\x88\x53\x9b\x3d\x69\x8e\x2f\x43\x2f\x3b\x34\x5c\x33\xd4\x03\xed\x49\xfc\x17\x7f\xdc\x5b\x34\xe7\xd4\x23\x72\xad\x59\x68\x7c\x9c\x29\x1c\xae\x65\x0a\x25\x9a\x35\xd2\x96\x3d\xcd\xc7\x16\x79\xd4\xaf\x35\xbc\x4f\xd5\x85\xb1\xb0\x1b\x4a\x7c\xd3\xb2\x2b\x73\xb2\x9f\xc1\xef\x61\x59\x06\x1a\x0d\xec\x4f\x4e\xe1\xb9\x13\xa4\xe3\x2c\x77\x93\x3c\x53\xbc\xe6\x24\xa9\x1e\xc7\xc3\x06\x2b\x70\x1c\x29\x45\x45\x96\xf3\x3b\x85\x04\x56\xae\x3f\xa0\x12\xd2\x50\x9c\x1f\x18\xb8\x41\xf1\xdd\x5e\x52\xb1\x82\x76\x07\x20\xda\xd2\x02\xad\xd4\x46\x14\x87\xa7\xd5\x07\x53\xd1\x64\xc9\x58\xf8\x10\x17\x50\x4d\x27\x72\xcf\xe0\x29\xd7\xd6\x47\x9e\x95\xc4\x82\xcf\xe0\x2f\x1f\x71\x7b\x39\x64\x39\x87\xb6\x77\x79\xa7\x2e\x7d\xe8\x3b\xb2\xa6\x36\x4e\x6a\x55\x0c\x69\xcd\x25\xf7\xba\x3c\x0f\x06\x8c\x6a\xd3\x48\x83\x10\xd7\xce\xde\x2e\x7a\x54\x3e\x01\xb4\x3e\x6c\x2b\xfc\x1e\x9d\x68\x6a\x97\xd8\x9c\x58\x65\x72\x23\xb3\x5a\x04\x40\x48\xeb\x2e\x14\xdc\xbc\xbb\x8b\x6e\x12\x6d\xa5\x95\x45\x68\x50\x0c\x5a\xe7\x4f\x00\x3d\x8b\xf6\x38\x71\xc2\xe7\x2d\x7c\xb6\x15\x3d\xf8\xf2\x43\xd3\x32\x4a\x67\x69\xf0\x90\xec\x6a\xce\x4d\x12\xb8\x77\xa6\xe6\x9c\x44\x73\x6e\x45\x6b\xd3\xa6\xc6\x62\x64\x0d\x54\xd4\xc4\xf2\x71\x96\x3f\xa8\x92\xaa\x90\x0a\x5b\x69\x9c\x0b\x5f\xfb\x13\x7a\x03\xe2\x3f\x35\xbd\x37\xe8\x2a\xbb\x89\xbf\x53\x93\x7d\x43\x5e\xa0\x2f\x0d\x78\x4a\xea\x6f\x80\xf6\x97\x4c\x0a\xee\x97\x09\x4e\xf6\x30\x61\xb8\x5f\xce\x4f\x1f\x0e\x2e\x6a\x27\xb1\x38\x3d\x99\x38\x40\x71\x28\xcd\x38\x94\x5a\x1c\x20\xf9\x67\x4b\x3a\xee\x97\x89\x07\x00\x51\x5f\xec\x5d\x29\xc7\xba\x89\xe9\x91\x17\x19\xae\x44\x5d\xb8\x45\x9b\xe9\x4c\x78\xef\x12\x61\x72\xcc\x95\x47\x4f\xb4\xcf\x3e\xee\xf5\x53\xfa\x8c\x47\xb7\x32\xed\xd9\xbb\xef\xc9\xe8\x2e\x0d\xdb\xf4\x4e\x52\x30\x6c\xd8\xb9\x3b\x9f\x75\x13\xec\xa7\xad\x5b\x1a\xc9\x1f\x11\x19\x56\x30\x53\x7a\xb7\x27\x55\x5a\xd4\x19\x76\xcf\xfc\x66\x60\xe5\x6f\xc1\x43\xc9\xd2\x47\x9e\x43\x52\x53\x72\xc5\x77\xa5\x58\xe3\x7d\x85\xe9\x0e\x44\x74\x87\x16\x4b\x5d\x3b\x8e\x9e\xd4\x0e\x6a\x8b\x19\x88\x3e\x53\xa4\x26\x69\xdf\x6c\x86\x03\x0c\xcd\x62\x82\xf6\x7d\x15\xc4\xf9\x35\x43\x17\xea\xd5\x84\xd9\x86\x31\xa9\xa0\x92\xcf\x84\xdc\xbe\xc2\x64\x9d\xc4\x70\xd5\xf5\xb7\xcf\xd7\xdf\x7e\x7d\x56\xaa\xcc\x9a\x74\x0a\xf4\xd8\x9d\xca\xd2\x52\x7b\xe6\x5a\xdf\xc6\xc2\x03\x94\x2e\x47\xaf\x06\x4b\xab\x8b\xda\xc5\xb0\xc2\x8f\xef\xdf\x84\x00\xe5\x09\x91\xf2\xc0\x2b\xe1\x04\x7f\x6a\x16\x2c\x78\x5d\x6e\x92\xb4\x9f\x63\xfe\x92\xd0\x84\x8f\xca\x9c\x02\x7b\xff\xdd\x2d\x5c\x7f\xf3\xaf\x7f\x9c\x25\x93\x89\x89\xa9\x83\xd5\x2b\x09\x95\x74\x90\x52\x33\x39\xbf\x70\x70\xc9\xbf\xae\x2a\xb5\xbe\x3c\x67\xa1\x86\x5c\x84\x35\xe9\xa7\x70\x04\x8f\xb8\xf5\x47\xfa\x67\x5d\xce\x38\x19\x9c\xb6\x30\xb4\x47\x18\xe7\x00\xd3\x0e\x04\xed\xa1\x78\x16\x28\xdd\x83\x9f\x7d\xb6\x33\x11\x90\x0e\xfb\x8a\x31\x20\xfa\x39\x40\xe8\x67\x01\xa0\x9f\x0b\x7c\xfe\x39\x80\xe7\x88\xcf\x88\x03\xce\x4f\x70\x57\x2d\x9a\x41\x3b\x05\x68\x1e\xc3\xc9\x81\xc4\xdc\x38\xc8\x8c\x43\xc9\x08\xd9\x2f\x0d\x30\x07\x17\xf0\x0c\x7f\x59\x48\xf5\x38\x0c\x9c\xde\x50\x8b\x26\xab\x1b\xf0\x52\xdb\x60\x83\x1c\xce\xd8\x3e\x54\xe6\xef\x88\x2c\x3d\x14\x39\xe6\x4e\x83\xad\xcd\x8a\x53\xa4\xed\x81\x14\x64\x3a\xad\xcb\x60\xc3\x33\xc8\x84\xcd\x97\x5a\x98\xcc\xce\x00\x5d\x7a\x0e\x6a\x22\x86\x47\x00\x53\x13\xbc\x03\x3b\xfb\x3c\xf4\xc8\x7d\x84\xab\x4f\x7b\xdb\xee\xd5\xfe\x6d\xbb\x83\xdb\x45\x34\x31\x52\x6a\x7c\xae\x0a\x0f\x30\x06\x0c\xa0\xaa\x4d\xa5\x6d\x1b\x5c\x68\xb1\xcf\xc2\x12\xb5\x89\xdc\x63\xd8\x63\xfc\x47\x53\x10\x21\x99\x8a\xa2\xd8\x02\x1b\xa7\x05\x0a\x47\xf0\x84\x4b\x2b\x1d\xaf\xbb\x41\x1b\xc9\xeb\x7c\x6a\xd5\x2e\x85\x64\x1d\x40\x33\xac\xe0\xdf\xef\xda\xc5\xf6\x07\x1d\x52\x41\x96\x83\xbb\x04\x46\x52\x5d\xb2\xca\x3b\x2f\xa6\x17\x3e\xf8\xd8\xbf\xc3\xa4\xa9\xce\x70\xe6\xf1\x96\x3d\x56\x42\x8a\x38\x95\x48\x1f\x09\x8d\x85\xcb\x65\x03\x3c\x4c\x31\x94\x5b\xb2\x91\xd4\x31\x28\x1d\xd9\x60\xb4\x08\xa8\x1f\xe7\x68\xb3\x16\x4a\xfe\x16\x91\xc6\xb0\x71\x60\x29\xe4\x14\xed\x7a\x4d\xed\x42\x6a\x8e\x3b\x7d\x84\x3a\x4d\xbe\xe5\xf7\xb6\xb9\xe5\x47\x83\x76\xbd\x1e\xdf\x1e\xf9\xbc\xa6\x94\x32\x70\xd9\x19\xd4\x12\xbf\x94\x25\x85\x08\x3e\x68\x46\xe1\x62\xf3\xdd\xab\xdd\x05\x58\x42\x08\x91\xf0\x3f\xa4\x14\xfe\xcd\x47\x6c\x4f\x13\x9b\x58\x0c\xbc\xec\x71\x19\x90\x8b\x4f\xe6\x1a\x29\x0a\x7e\x2b\x40\x7d\x83\x65\xef\x6e\xec\xc0\x9d\x83\x08\xfa\x54\x64\x24\xb4\x2a\x2d\x62\x60\x0b\x0d\x8f\x41\xfa\xd6\x65\x90\xf9\x98\x3a\xee\xe3\xae\xae\x36\xfa\x07\x2b\xfe\xb6\xe9\x01\xe7\x1c\xdb\x08\x3e\xf1\x14\x7a\xf9\xe7\x1c\x1f\xe7\x42\xb8\x67\x97\x4a\xec\x15\xd2\xdc\xe0\x5a\x72\x32\xef\xf0\x9d\x47\x3b\x43\xad\x17\x0d\x5f\x09\xaf\x61\x73\x94\xde\x66\x14\x45\x51\x40\xa1\x9f\xd0\xa4\x04\xb1\xfa\x10\xce\xa0\x90\xfc\x63\x99\x31\x29\x05\xce\xef\xa9\x75\xe3\xd1\x51\xd5\xe5\xfe\x75\xd8\x4c\xae\xf8\x96\x99\xf3\x54\xfb\x57\x59\x6c\x84\x2c\x38\xec\x72\xfa\x3b\x65\xc9\x0c\x01\xda\x41\xf6\x37\xf1\x0d\xd2\x40\xbf\xd8\xa6\x79\xee\xcd\xe4\xa8\x96\x54\xf9\xa8\xb2\xe7\x52\xe5\xdc\x4f\xfc\xa8\xb6\xe1\xb2\xd7\x35\xf4\x78\x0d\x5b\x2f\x1f\x48\x42\x9f\xe6\xca\x64\xe4\x61\xc5\x20\x45\x5f\x3e\xe6\x91\x45\xf4\x28\xd7\x3f\xbe\x18\x7b\x6a\x31\x65\x62\xc7\x4c\x0e\x9c\x6f\x9f\xfb\xf8\x62\x80\xe4\xe1\xad\xee\xa1\xfc\xfb\xa4\xf4\xef\xd0\x73\x8c\x9e\x59\x9c\xf4\x28\x63\x70\x1e\x87\xd7\xa4\x4f\x79\x9a\x31\x48\x78\xf7\x6c\xe3\xc4\x07\x1a\xc3\x54\xf7\x1f\x6f\x7c\x91\x67\x1a\x3b\xc1\xc5\xaf\x72\xfb\x32\x69\xe9\x87\x1f\x6e\xf8\xf2\x31\xcf\x37\x06\x67\x31\xfc\xb4\xe3\x84\x59\xc4\x01\x60\xcf\x0c\x3e\xfe\xb1\xc7\x09\x9c\x0d\x3f\xfc\xe8\xe1\xee\x94\xe7\x1f\xc3\xca\xda\x3e\x0d\x39\xf9\x11\xc8\x20\xdd\xfd\x07\x22\x27\x3d\x05\x19\xe6\xf7\xe0\x99\xc8\xc7\x3f\x08\xf1\x65\xfc\x59\x88\x2f\x63\x8f\x43\x3a\xad\xa6\xae\xe5\x09\x0f\x45\x06\x45\x13\xbc\xe3\x5f\xd8\xd3\x8c\xdd\xc5\x1a\x78\x46\x12\x1a\xf4\x3e\x26\x09\x1f\xa3\x4f\x4a\x42\x83\xde\xec\x64\x97\xfd\xd1\x0b\x37\x31\x7f\xc6\x37\x8d\xce\xbe\x6c\x13\x77\x5d\x67\xa7\x5f\x7a\x05\x75\xf2\x56\xb1\x6f\x84\x79\x1f\x5e\xee\xa1\x7c\x50\x15\xc0\x32\x6c\x5e\x8a\xa2\xca\xc5\xcb\x5d\x5d\xf3\xcf\x09\xfc\x3f\x0d\xe8\x7c\x6e\xf0\x5a\xb6\x00\x67\x6a\x3f\x9a\x75\xda\x88\x35\x36\x35\xff\x0f\x00\x00\xff\xff\xc2\x0b\x21\x16\x94\x41\x00\x00")

func metaAppscodeCom_resourcedescriptorsYamlBytes() ([]byte, error) {
	return bindataRead(
		_metaAppscodeCom_resourcedescriptorsYaml,
		"meta.appscode.com_resourcedescriptors.yaml",
	)
}

func metaAppscodeCom_resourcedescriptorsYaml() (*asset, error) {
	bytes, err := metaAppscodeCom_resourcedescriptorsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "meta.appscode.com_resourcedescriptors.yaml", size: 16788, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
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
	"meta.appscode.com_resourceclasses.yaml":     metaAppscodeCom_resourceclassesYaml,
	"meta.appscode.com_resourcedescriptors.yaml": metaAppscodeCom_resourcedescriptorsYaml,
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
	"meta.appscode.com_resourceclasses.yaml":     &bintree{metaAppscodeCom_resourceclassesYaml, map[string]*bintree{}},
	"meta.appscode.com_resourcedescriptors.yaml": &bintree{metaAppscodeCom_resourcedescriptorsYaml, map[string]*bintree{}},
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
