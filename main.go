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
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultArgoAppsDir = "argoapps"
	argoAppsDirEnvVar          = "UPDATE_VERSION_DIR"
)

var (
	argoAppsDir string
	dryRun      bool
	checkOnly   bool
)

func init() {
	// Default value, can be overridden by env var or flag
	argoAppsDir = defaultArgoAppsDir

	// Check environment variable
	if v := os.Getenv(argoAppsDirEnvVar); v != "" {
		argoAppsDir = v
	}

	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]

		switch a {
		case "--dry-run", "-n":
			dryRun = true

		case "--check", "-C":
			checkOnly = true

		case "--help", "-h":
			printUsage()
			os.Exit(0)

		case "--dir", "-d":
			if i+1 >= len(os.Args) {
				fail("--dir requires a directory path")
			}
			argoAppsDir = os.Args[i+1]
			i++

		default:
			if strings.HasPrefix(a, "-test.") {
				// Skip Go test framework flags
				continue
			}
			if strings.HasPrefix(a, "-") {
				fail("unknown flag: " + a)
			}
		}
	}

	if dryRun && checkOnly {
		fail("--dry-run and --check cannot be used together")
	}
}

func main() {
	charts, err := discoverCharts(argoAppsDir)
	must(err)

	if len(charts) == 0 {
		fail("no charts with artifacthub comments found in " + argoAppsDir)
	}

	if checkOnly {
		logf("discovered %d chart(s) with artifacthub comments:", len(charts))
		for _, c := range charts {
			logf("  %s → %s", c.File, c.Repo)
		}
		return
	}

	for _, c := range charts {
		must(updateChart(c.File, c.Repo))
	}
}

func printUsage() {
	exe := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, `Usage:
  %s [flags]

Description:
  Updates Argo CD Application Helm chart versions by scanning for manifests
  with "# artifacthub: org/repo" comments and fetching the latest version
  from ArtifactHub.

License:
  GNU GPL v3.0 only - https://spdx.org/licenses/GPL-3.0-only.html

Flags:
  -d, --dir <path>    Path to argoapps directory (default: %s)
  -n, --dry-run       Show git diff without modifying files
  -C, --check         Discover charts and show what would be updated
  -h, --help          Show this help message

Environment:
  %s    Directory path (used if --dir is not provided)

Exit codes:
  0  Success
  1  Error

Examples:
  %s
  %s --dir ./my-apps
  %s --dry-run
  %s=./my-apps %s --check

`, exe, defaultArgoAppsDir, argoAppsDirEnvVar, exe, exe, exe, argoAppsDirEnvVar, exe)
}
