// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package exec provides template functions for encoding content.
package exec

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"strings"

	"github.com/spf13/cast"
)

// New returns a new instance of the encoding-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "encoding" namespace.
type Namespace struct{}

// Graphviz will execute Graphviz with the input and return an SVG.
func (ns *Namespace) Graphviz(input interface{}) (template.HTML, error) {
	gvout, err := ns.External(input, "dot", "-Tsvg")
	if err != nil {
		return "", err
	}
	return template.HTML(gvout[strings.Index(string(gvout), "<svg"):]), nil
}

// External will execute the specified command, passing stdin to the exec and returning stdout.
func (ns *Namespace) External(a ...interface{}) (template.HTML, error) {
	input, err := cast.ToStringE(a[0])
	if err != nil {
		return "", err
	}
	lookpath, err := cast.ToStringE(a[1])
	if err != nil {
		return "", err
	}
	path, err := exec.LookPath(lookpath)
	if err != nil {
		return "", err
	}
	var args []string
	if len(a) > 2 {
		for _, v := range a[2:len(a)] {
			args = append(args, cast.ToString(v))
		}
	}
	cmd := exec.Command(path, args...)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return template.HTML(fmt.Sprintf("%s", stdout.String())), nil
	// return template.HTML(fmt.Sprintf("%#v", cmd)), nil
}
