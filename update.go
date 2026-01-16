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
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func updateChart(file, repo string) error {
	path := filepath.Join(argoAppsDir, file)

	docs, err := readYAMLDocuments(path)
	if err != nil {
		return err
	}

	var current string
	for _, d := range docs {
		if kind(d) == "Application" {
			current = getTargetRevision(d)
			break
		}
	}

	if current == "" {
		return fmt.Errorf("failed to read current version in %s", file)
	}

	latest, err := artifactHubLatestVersion(repo)
	if err != nil {
		return err
	}

	if versionLess(current, latest) {
		logf("%s: %s → %s", file, current, latest)

		for _, d := range docs {
			if kind(d) == "Application" {
				setTargetRevision(d, latest)
			}
		}

		if dryRun {
			return showDiff(path, docs)
		}
		return writeYAMLDocuments(path, docs)
	}

	logf("%s: already up to date (%s)", file, current)
	return nil
}

func showDiff(path string, docs []*yaml.Node) (err error) {
	return showDiffInternal(path, docs)
}
