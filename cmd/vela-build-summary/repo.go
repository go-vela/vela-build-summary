// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Repo represents the plugin configuration for repo information.
type Repo struct {
	// repository for the build
	Name string
	// organization for the build
	Org string
}

// Validate verifies the Repo is properly configured.
func (r *Repo) Validate() error {
	logrus.Trace("validating repo plugin configuration")

	// verify org is provided
	if len(r.Org) == 0 {
		return fmt.Errorf("no repo org provided")
	}

	// verify repo is provided
	if len(r.Name) == 0 {
		return fmt.Errorf("no repo name provided")
	}

	return nil
}
