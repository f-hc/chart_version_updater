// SPDX-License-Identifier: GPL-3.0-only
//
// Copyright (C) 2026 f-hc <207619282+f-hc@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func readYAMLDocuments(path string) (docs []*yaml.Node, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	dec := yaml.NewDecoder(f)
	for {
		var n yaml.Node
		if err := dec.Decode(&n); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		docs = append(docs, &n)
	}
	return docs, nil
}

func writeYAMLDocuments(path string, docs []*yaml.Node) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)

	for _, d := range docs {
		if err := enc.Encode(d); err != nil {
			return err
		}
	}
	return nil
}

func docRoot(n *yaml.Node) *yaml.Node {
	if n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		return n.Content[0]
	}
	return n
}

func kind(n *yaml.Node) string {
	return lookup(docRoot(n), "kind")
}

func getTargetRevision(n *yaml.Node) string {
	return lookup(docRoot(n), "spec", "source", "targetRevision")
}

func setTargetRevision(n *yaml.Node, v string) {
	set(docRoot(n), v, "spec", "source", "targetRevision")
}

const artifactHubPrefix = "# artifacthub:"

// getArtifactHubRepo extracts the ArtifactHub repository path from a YAML comment.
// It looks for a comment in the format "# artifacthub: org/repo" at the top of the file.
// In yaml.v3, this comment is attached to the first key of the root mapping node.
func getArtifactHubRepo(n *yaml.Node) string {
	root := docRoot(n)

	// The comment is attached to the first key in a mapping node
	if root.Kind == yaml.MappingNode && len(root.Content) > 0 {
		firstKey := root.Content[0]
		if strings.HasPrefix(firstKey.HeadComment, artifactHubPrefix) {
			return strings.TrimSpace(strings.TrimPrefix(firstKey.HeadComment, artifactHubPrefix))
		}
	}

	return ""
}

func lookup(n *yaml.Node, path ...string) string {
	cur := n
	for _, p := range path {
		cur = mapGet(cur, p)
		if cur == nil {
			return ""
		}
	}
	return cur.Value
}

func set(n *yaml.Node, value string, path ...string) {
	cur := n
	for _, p := range path {
		next := mapGet(cur, p)
		if next == nil {
			next = &yaml.Node{Kind: yaml.MappingNode}
			mapSet(cur, p, next)
		}
		cur = next
	}
	cur.Value = value
}

func mapGet(n *yaml.Node, key string) *yaml.Node {
	if n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

func mapSet(n *yaml.Node, key string, val *yaml.Node) {
	n.Content = append(
		n.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		val,
	)
}
