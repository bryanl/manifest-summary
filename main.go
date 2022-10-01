//    Copyright 2018 Bryan Liles
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	if err := run(os.Stdin); err != nil {
		logrus.WithError(err).Error("summarizing")
		os.Exit(1)
	}
}

func run(r io.Reader) error {
	objects, err := decode(r)
	if err != nil {
		return errors.Wrap(err, "decoding stdin")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"API Version", "Kind", "Name"})

	for _, object := range objects {
		s, err := summarize(object)
		if err != nil {
			return errors.Wrap(err, "creating summary for object")
		}

		table.Append(s.row())
	}

	table.Render()

	return nil
}

func decode(r io.Reader) ([]map[string]interface{}, error) {
	decoder := yaml.NewDecoder(r)

	var out []map[string]interface{}

	for {
		var decoded map[string]interface{}
		err := decoder.Decode(&decoded)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, errors.Wrap(err, "decoding YAML input")
		}

		out = append(out, decoded)
	}

	return out, nil
}

type summary struct {
	APIVersion string
	Kind       string
	Name       string
}

func (s *summary) row() []string {
	return []string{s.APIVersion, s.Kind, s.Name}
}

func summarizeItem(m map[string]interface{}) (*summary, error) {

	kind, ok := m["kind"].(string)
	if !ok {
		return nil, errors.New("finding kind")
	}

	apiVersion, ok := m["apiVersion"].(string)
	if !ok {
		return nil, errors.New("finding apiVersion")
	}

	metadata, ok := m["metadata"].(map[interface{}]interface{})
	if !ok {
		return nil, errors.Errorf("finding metadata")
	}

	s := &summary{
		APIVersion: apiVersion,
		Kind:       kind,
	}

	nameRaw := metadata["name"]
	if nameRaw != nil {
		s.Name, ok = nameRaw.(string)
		if !ok {
			return nil, errors.New("finding name")
		}

		return s, nil
	}

	generateNameRaw := metadata["generateName"]
	if generateNameRaw != nil {
		s.Name, ok = generateNameRaw.(string)
		if !ok {
			return nil, errors.New("finding generateName")
		}

		return s, nil
	}

	return nil, errors.New("unable to find object name")
}

func summarize(m map[string]interface{}) (*summary, error) {

	kind, ok := m["kind"].(string)
	if !ok {
		return nil, errors.New("finding kind")
	}

	if kind == "List" {
		items, ok := m["items"].([]interface{})
		if !ok {
			return nil, errors.New("finding items")
		}
		for _, item := range items {
			summarizeItem(item.(map[string]interface{}))
		}
	} else {
		summarizeItem(m)
	}

}
