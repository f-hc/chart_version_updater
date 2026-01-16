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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

var artifactHubAPI = "https://artifacthub.io/api/v1/packages/helm"

// httpClient is a shared HTTP client with timeout for external requests.
var httpClient = &http.Client{Timeout: 60 * time.Second}

// ArtifactHubResponse represents the API response structure.
type ArtifactHubResponse struct {
	AvailableVersions []struct {
		Version string `json:"version"`
	} `json:"available_versions"`
}

// artifactHubLatestVersion fetches the latest stable version from ArtifactHub.
func artifactHubLatestVersion(repo string) (version string, err error) {
	resp, err := httpClient.Get(artifactHubAPI + "/" + repo)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("artifacthub HTTP %d", resp.StatusCode)
	}

	var data ArtifactHubResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	var versions []string
	for _, v := range data.AvailableVersions {
		if !strings.Contains(v.Version, "-") { // skip pre-release
			versions = append(versions, v.Version)
		}
	}

	if len(versions) == 0 {
		return "", errors.New("no stable versions found")
	}

	sort.Slice(versions, func(i, j int) bool {
		return versionLess(versions[i], versions[j])
	})

	return versions[len(versions)-1], nil
}
