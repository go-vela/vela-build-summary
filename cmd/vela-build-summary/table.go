// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/library"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
)

// table is a helper function to output the provided build summary in a table.
//
// The summary includes generic information on the steps and services in the
// build, such as name, number, status and duration of runtime. Also in the
// table are some more fine grained metrics on log size and rate of logs
// produced throughout the lifecycle of each resource.
func table(build *api.Build, logs *[]library.Log, services *[]library.Service, steps *[]library.Step) error {
	logrus.Debug("creating table for build summary")

	// create a new table
	//
	// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#New
	table := uitable.New()

	// set column width for table to 50
	//
	// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table
	table.MaxColWidth = 50

	// ensure the table is always wrapped
	//
	// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table
	table.Wrap = true

	logrus.Trace("adding headers to build summary table")
	// set of build fields we display in a table
	//
	// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table.AddRow
	table.AddRow("TYPE", "NAME", "NUMBER", "STATUS", "DURATION", "LOG LINES", "LOG SIZE", "LOG RATE")

	// create a variable to track the lines of logs for the build
	buildLines := new(int)
	// create a variable to track the size of logs for the build
	buildSize := new(uint64)

	// add the service rows to the table
	serviceRows(table, logs, services, buildLines, buildSize)

	// add the step rows to the table
	stepRows(table, logs, steps, buildLines, buildSize)

	// add a separation row to the table with the specified values
	//
	// https://pkg.go.dev/github.com/gosuri/uitable?tab=doc#Table.AddRow
	table.AddRow("----------", "--------------------", "----------", "----------", "----------", "----------", "---------------", "---------------")

	// add the build row to the table
	buildRow(table, build, buildLines, buildSize)

	// ensure we output table to stdout
	fmt.Fprintln(os.Stdout, table)

	return nil
}
