// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bootstrap

import (
	meshconfig "istio.io/api/mesh/v1alpha1"
	"fmt"
	"path"
	"os"
	"io/ioutil"
	"text/template"
)

// Generate the envoy v2 bootstrap configuration, using template.
const (
	// EpochFileTemplate is a template for the root config JSON
	EpochFileTemplate = "envoy-rev%d.yaml"
)

func configFile(config string, epoch int) string {
	return path.Join(config, fmt.Sprintf(EpochFileTemplate, epoch))
}

func WriteBootstrap(config *meshconfig.ProxyConfig, epoch int) (string, error) {
	if err := os.MkdirAll(config.ConfigPath, 0700); err != nil {
		return "", err
	}
	// attempt to write file
	fname := configFile(config.ConfigPath, epoch)

	f, err := ioutil.ReadFile(config.CustomConfigFile)
	if err != nil {
		return "", err
	}
	t, err := template.New("bootstrap").Parse(string(f))
	if err != nil {
		return "", err
	}

	opts := map[string]interface{} {
		"config": config,
    }
	fout, err := os.Create(fname)
	if err != nil {
		return "", err
	}

	// Execute needs some sort of io.Writer
    err = t.Execute(fout, opts)

	return fname, nil
}