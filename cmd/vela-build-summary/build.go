// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-vela/types/constants"
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

// buildDuration is a helper function to calculate the total duration
// the build ran for in a more consumable, human readable format.
func buildDuration(b *library.Build) string {
	logrus.Debug("calculating duration of build for build summary table")

	// capture finished unix timestamp from service
	f := time.Unix(b.GetFinished(), 0)
	// capture started unix timestamp from service
	st := time.Unix(b.GetStarted(), 0)

	// check if build is in a pending or running state
	if strings.EqualFold(b.GetStatus(), constants.StatusPending) ||
		strings.EqualFold(b.GetStatus(), constants.StatusRunning) {
		// set a default value to display for the build
		f = time.Unix(time.Now().UTC().Unix(), 0)
	}

	// get the duration by subtracting the build started
	// timestamp from the build finished timestamp.
	d := f.Sub(st)

	// return duration in a human readable form
	return d.String()
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
func buildRow(table *uitable.Table, build *library.Build, buildLines *int, buildSize *uint64) {
	logrus.Debug("adding build information to build summary table")

	// calculate duration based off the build timestamps
	duration := buildDuration(build)

	// calculate rate based off build duration and size
	rate := fmt.Sprintf("%d B/s", buildRate(duration, *buildSize))

	// add a row to the table with the specified values
	//
	// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table.AddRow
	//
	// nolint: lll // ignore line length due to parameters
	table.AddRow("build", "", build.GetNumber(), build.GetStatus(), duration, *buildLines, humanize.Bytes(*buildSize), rate)
}
