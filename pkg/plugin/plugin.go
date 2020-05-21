/*
Copyright © Armory, Inc

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

package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//go:generate easytags $GOFILE json:camel

type Repository struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Metadata struct {
	LastModified   uint       `json:"lastModified"`
	CreateTS       uint       `json:"createTS"`
	ID             string     `json:"id"`
	LastModifiedBy string     `json:"lastModifiedBy"`
	Provider       string     `json:"provider"`
	Description    string     `json:"description"`
	Releases       []*Release `json:"releases"`
}

type Release struct {
	SHA512Sum      string `json:"sha512Sum"`
	LastModified   uint   `json:"lastModified"`
	URL            string `json:"url"`
	LastModifiedBy string `json:"lastModifiedBy"`
	Version        string `json:"version"`
	Date           string `json:"date"`
	State          string `json:"state"`
	Requires       string `json:"requires"`
}

type Loader struct {
	PluginDir     string
	BinaryAddress string
}

func (l *Loader) LoadMetadata() (*Metadata, error) {
	r, err := os.Open(path.Join(l.PluginDir, "plugin-info.json"))
	if err != nil {
		return nil, err
	}

	var m Metadata
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}

	if len(m.Releases) == 0 {
		return nil, fmt.Errorf("could not find release in plugin metadata")
	}

	// Pick the first release. If it's generated from `gradle releaseBundle`,
	// it should be the only one.
	rel := m.Releases[0]
	rel.URL = l.BinaryAddress

	return &m, nil
}

func (l *Loader) LocateBinary() (string, error) {
	files, err := ioutil.ReadDir(l.PluginDir)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".zip") {
			return path.Join(l.PluginDir, f.Name()), nil
		}
	}

	return "", fmt.Errorf("could not locate plugin binary zip in %q", l.PluginDir)
}
