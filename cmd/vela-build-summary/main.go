// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"time"

	"github.com/go-vela/vela-build-summary/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-build-summary"
	app.HelpName = "vela-build-summary"
	app.Usage = "Vela Build Summary plugin for capturing a summary of a build"
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "BUILD_SUMMARY_LOG_LEVEL"},
			FilePath: "/vela/parameters/build-summary/log_level,/vela/secrets/build-summary/log_level",
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},

		// Build Flags

		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_NUMBER", "BUILD_SUMMARY_NUMBER", "VELA_BUILD_NUMBER"},
			FilePath: "/vela/parameters/build-summary/number,/vela/secrets/build-summary/number",
			Name:     "build.number",
			Usage:    "provide the number for the build",
		},

		// Config Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_SERVER", "BUILD_SUMMARY_SERVER", "VELA_ADDR"},
			FilePath: "/vela/parameters/build-summary/server,/vela/secrets/build-summary/server",
			Name:     "config.server",
			Usage:    "Vela server to authenticate with",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TOKEN", "BUILD_SUMMARY_TOKEN", "VELA_NETRC_PASSWORD"},
			FilePath: "/vela/parameters/build-summary/token,/vela/secrets/build-summary/token",
			Name:     "config.token",
			Usage:    "user token to authenticate with the Vela server",
		},

		// Repo Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_ORG", "BUILD_SUMMARY_ORG", "VELA_REPO_ORG"},
			FilePath: "/vela/parameters/build-summary/org,/vela/secrets/build-summary/org",
			Name:     "repo.org",
			Usage:    "provide the organization name for the build",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REPO", "BUILD_SUMMARY_REPO", "VELA_REPO_NAME"},
			FilePath: "/vela/parameters/build-summary/repo,/vela/secrets/build-summary/repo",
			Name:     "repo.name",
			Usage:    "provide the repository name for the build",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-build-summary",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/build-summary",
		"registry": "https://hub.docker.com/r/target/vela-build-summary",
	}).Info("Vela Build Summary Plugin")

	// create the plugin
	p := &Plugin{
		// build configuration
		Build: &Build{
			Number: c.Int("build.number"),
		},
		// config configuration
		Config: &Config{
			AppName:    c.App.Name,
			AppVersion: c.App.Version,
			Server:     c.String("config.server"),
			Token:      c.String("config.token"),
		},
		// repo configuration
		Repo: &Repo{
			Org:  c.String("repo.org"),
			Name: c.String("repo.name"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
