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
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

func showDiffInternal(ctx context.Context, path string, docs []*yaml.Node) (err error) {
	tmp, err := os.CreateTemp("", "update-version-*.yaml")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}

	defer func() {
		if removeErr := os.Remove(tmp.Name()); removeErr != nil && err == nil {
			err = removeErr
		}
	}()

	enc := yaml.NewEncoder(tmp)
	enc.SetIndent(yamlIndent)

	if err = encodeStream(enc, docs); err != nil {
		return err
	}

	if err = enc.Close(); err != nil {
		return fmt.Errorf("close encoder: %w", err)
	}

	if err = tmp.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	//nolint:gosec // path is validated to be within base directory in config.go
	cmd := exec.CommandContext(ctx, "git", "diff", "--no-index", "--", path, tmp.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return nil // git diff returns 1 when files differ
		}

		return fmt.Errorf("run git diff: %w", err)
	}

	return nil
}
