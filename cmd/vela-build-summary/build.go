// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-vela/types/library"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
)

// Build represents the plugin configuration for build information.
type Build struct {
	// number for the build
	Number int
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build plugin configuration")

	// verify number is provided
	if b.Number == 0 {
		return fmt.Errorf("no build number provided")
	}

	return nil
}

// buildRate is a helper function to calculate the total size of logs
// a build produced over the total duration a build ran for.
func buildRate(duration string, size uint64) int64 {
	// parse the string duration into a timestamp duration
	d, _ := time.ParseDuration(duration)

	// calculate the timestamp duration in seconds
	s := (float64(d) / float64(time.Second))

	// return the rate of bytes per second
	return int64(float64(size) / s)
}

// buildRow is a helper function to produce a build row in the build summary table.
func buildRow(table *uitable.Table, b *library.Build, buildLines *int, buildSize *uint64) {
	logrus.Debug("adding build information to build summary table")

	// calculate duration based off the build timestamps
	duration := b.Duration()

	// calculate rate based off build duration and size
	rate := fmt.Sprintf("%d B/s", buildRate(duration, *buildSize))

	// add a row to the table with the specified values
	//
	// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table.AddRow
	table.AddRow("build", "", b.GetNumber(), b.GetStatus(), duration, *buildLines, humanize.Bytes(*buildSize), rate)
}
