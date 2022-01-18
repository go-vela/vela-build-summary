// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// nolint: dupl // ignore similar code with service
package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
)

// stepDuration is a helper function to calculate the total duration
// a step ran for in a more consumable, human readable format.
func stepDuration(s *library.Step) string {
	logrus.Debugf("calculating duration of step %s for build summary table", s.GetName())

	// capture finished unix timestamp from step
	f := time.Unix(s.GetFinished(), 0)
	// capture started unix timestamp from step
	st := time.Unix(s.GetStarted(), 0)

	// check if step is in a pending or running state
	if strings.EqualFold(s.GetStatus(), constants.StatusPending) ||
		strings.EqualFold(s.GetStatus(), constants.StatusRunning) {
		// set a default value to display for the step
		f = time.Unix(time.Now().UTC().Unix(), 0)
	}

	// get the duration by subtracting the step started
	// timestamp from the step finished timestamp.
	d := f.Sub(st)

	// return duration in a human readable form
	return d.String()
}

// stepLines is a helper function to calculate the total lines of logs
// a step produced by measuring the newlines (\n) in that log entry.
func stepLines(s *library.Step, logs *[]library.Log) int {
	logrus.Debugf("calculating lines of logs for step %s for build summary table", s.GetName())

	// create a variable to track the lines of logs for the step
	var lines int

	// iterate through all logs in the list
	for _, log := range *logs {
		// check if the log step ID matches the step ID
		if log.GetStepID() != s.GetID() {
			continue
		}

		// capture the total lines for the logs
		lines = bytes.Count(log.GetData(), []byte("\n"))

		// break out of the for loop
		break
	}

	return lines
}

// stepRate is a helper function to calculate the total size of logs
// a step produced over the total duration a step ran for.
func stepRate(duration string, size uint64) int64 {
	// parse the string duration into a timestamp duration
	d, _ := time.ParseDuration(duration)

	// calculate the timestamp duration in seconds
	s := (float64(d) / float64(time.Second))

	// return the rate of bytes per second
	return int64(float64(size) / s)
}

// stepReverse is a helper function to sort the steps based off the
// step number and then flip the order they get displayed in.
func stepReverse(s []library.Step) []library.Step {
	logrus.Debug("reversing order of steps for build summary table")

	// sort the list of steps based off the step number
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].GetNumber() < s[j].GetNumber()
	})

	return s
}

// stepRows is a helper function to produce step rows in the build summary table.
//
// nolint: lll // ignore line length due to parameters
func stepRows(table *uitable.Table, logs *[]library.Log, steps *[]library.Step, buildLines *int, buildSize *uint64) {
	logrus.Debug("adding step information to build summary table")

	// iterate through all steps in the list
	for _, s := range stepReverse(*steps) {
		logrus.Tracef("adding step %s to build summary table", s.GetName())

		// calculate lines based off the step logs
		//
		// nolint: gosec // ignore memory aliasing
		lines := stepLines(&s, logs)

		// calculate size based off the step logs
		//
		// nolint: gosec // ignore memory aliasing
		size := stepSize(&s, logs)

		// calculate duration based off the step timestamps
		//
		// nolint: gosec // ignore memory aliasing
		duration := stepDuration(&s)

		// calculate rate based off step duration and size
		rate := fmt.Sprintf("%d B/s", stepRate(duration, size))

		// update the lines of build logs with the lines of step logs
		*buildLines = *buildLines + lines

		// update the size of the build logs with the size of the step logs
		*buildSize = *buildSize + size

		// add a row to the table with the specified values
		//
		// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table.AddRow
		//
		// nolint: lll // ignore line length due to parameters
		table.AddRow("step", s.GetName(), s.GetNumber(), s.GetStatus(), duration, lines, humanize.Bytes(size), rate)
	}
}

// stepSize is a helper function to calculate the total size of logs
// a step produced by measuring the data in that log entry.
func stepSize(s *library.Step, logs *[]library.Log) uint64 {
	logrus.Debugf("calculating size of logs for step %s for build summary table", s.GetName())

	// create a variable to track the size of logs for the step
	var size uint64

	// iterate through all logs in the list
	for _, log := range *logs {
		// check if the log step ID matches the step ID
		if log.GetStepID() != s.GetID() {
			continue
		}

		// capture the total size for the logs
		size = uint64(len(log.GetData()))

		// break out of the for loop
		break
	}

	return size
}
