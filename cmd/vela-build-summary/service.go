// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// nolint: dupl // ignore similar code with step
package main

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-vela/types/library"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
)

// serviceLines is a helper function to calculate the total lines of logs
// a service produced by measuring the newlines (\n) in that log entry.
func serviceLines(s *library.Service, logs *[]library.Log) int {
	logrus.Debugf("calculating lines of logs for service %s for build summary table", s.GetName())

	// create a variable to track the lines of logs for the service
	var lines int

	// iterate through all logs in the list
	for _, log := range *logs {
		// check if the log service ID matches the service ID
		if log.GetServiceID() != s.GetID() {
			continue
		}

		// capture the total lines for the logs
		lines = bytes.Count(log.GetData(), []byte("\n"))

		// break out of the for loop
		break
	}

	return lines
}

// serviceRate is a helper function to calculate the total size of logs
// a service produced over the total duration a service ran for.
func serviceRate(duration string, size uint64) int64 {
	// parse the string duration into a timestamp duration
	d, _ := time.ParseDuration(duration)

	// calculate the timestamp duration in seconds
	s := (float64(d) / float64(time.Second))

	// return the rate of bytes per second
	return int64(float64(size) / s)
}

// serviceReverse is a helper function to sort the services based off the
// service number and then flip the order they get displayed in.
func serviceReverse(s []library.Service) []library.Service {
	logrus.Debug("reversing order of services for build summary table")

	// sort the list of services based off the service number
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].GetNumber() < s[j].GetNumber()
	})

	return s
}

// serviceRows is a helper function to produce service rows in the build summary table.
//
// nolint: lll // ignore line length due to parameters
func serviceRows(table *uitable.Table, logs *[]library.Log, services *[]library.Service, buildLines *int, buildSize *uint64) {
	logrus.Debug("adding service information to build summary table")

	// iterate through all services in the list
	for _, s := range serviceReverse(*services) {
		logrus.Tracef("adding service %s to build summary table", s.GetName())

		// calculate lines based off the service logs
		//
		// nolint: gosec // ignore memory aliasing
		lines := serviceLines(&s, logs)

		// calculate size based off the service logs
		//
		// nolint: gosec // ignore memory aliasing
		size := serviceSize(&s, logs)

		// calculate duration based off the service timestamps
		duration := s.Duration()

		// calculate rate based off service duration and size
		rate := fmt.Sprintf("%d B/s", serviceRate(duration, size))

		// update the lines of build logs with the lines of service logs
		*buildLines = *buildLines + lines

		// update the size of the build logs with the size of the service logs
		*buildSize = *buildSize + size

		// add a row to the table with the specified values
		//
		// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table.AddRow
		//
		// nolint: lll // ignore line length due to parameters
		table.AddRow("service", s.GetName(), s.GetNumber(), s.GetStatus(), duration, lines, humanize.Bytes(size), rate)
	}
}

// serviceSize is a helper function to calculate the total size of logs
// a service produced by measuring the data in that log entry.
func serviceSize(s *library.Service, logs *[]library.Log) uint64 {
	logrus.Debugf("calculating size of logs for service %s for build summary table", s.GetName())

	// create a variable to track the size of logs for the service
	var size uint64

	// iterate through all logs in the list
	for _, log := range *logs {
		// check if the log service ID matches the service ID
		if log.GetServiceID() != s.GetID() {
			continue
		}

		// capture the total size for the logs
		size = uint64(len(log.GetData()))

		// break out of the for loop
		break
	}

	return size
}
