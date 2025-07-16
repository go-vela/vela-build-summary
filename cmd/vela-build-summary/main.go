// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/vela-build-summary/version"
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
	// Plugin Information
	cmd := cli.Command{
		Name:      "vela-build-summary",
		Usage:     "Vela Build Summary plugin for capturing a summary of a build",
		Copyright: "Copyright 2021 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Version: v.Semantic(),
		Action:  run,
	}
	// Plugin Flags

	cmd.Flags = []cli.Flag{

		&cli.StringFlag{
			Name:  "log.level",
			Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value: "info",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG_LEVEL"),
				cli.EnvVar("BUILD_SUMMARY_LOG_LEVEL"),
				cli.File("/vela/parameters/build-summary/log_level"),
				cli.File("/vela/secrets/build-summary/log_level"),
			),
		},

		// Build Flags

		&cli.IntFlag{
			Name:  "build.number",
			Usage: "provide the number for the build",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_NUMBER"),
				cli.EnvVar("BUILD_SUMMARY_NUMBER"),
				cli.EnvVar("VELA_BUILD_NUMBER"),
				cli.File("/vela/parameters/build-summary/number"),
				cli.File("/vela/secrets/build-summary/number"),
			),
		},

		// Config Flags

		&cli.StringFlag{
			Name:  "config.server",
			Usage: "Vela server to authenticate with",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_SERVER"),
				cli.EnvVar("BUILD_SUMMARY_SERVER"),
				cli.EnvVar("VELA_ADDR"),
				cli.File("/vela/parameters/build-summary/server"),
				cli.File("/vela/secrets/build-summary/server"),
			),
			Required: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				if strings.HasSuffix(v, "/") {
					return fmt.Errorf("invalid server address provided: address must not have trailing slash")
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:  "config.token",
			Usage: "user token to authenticate with the Vela server",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TOKEN"),
				cli.EnvVar("BUILD_SUMMARY_TOKEN"),
				cli.EnvVar("VELA_NETRC_PASSWORD"),
				cli.File("/vela/parameters/build-summary/token"),
				cli.File("/vela/secrets/build-summary/token"),
			),
		},

		// Repo Flags

		&cli.StringFlag{
			Name:  "repo.org",
			Usage: "provide the organization name for the build",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_ORG"),
				cli.EnvVar("BUILD_SUMMARY_ORG"),
				cli.EnvVar("VELA_REPO_ORG"),
				cli.File("/vela/parameters/build-summary/org"),
				cli.File("/vela/parameters/build-summary/org"),
			),
		},
		&cli.StringFlag{
			Name:  "repo.name",
			Usage: "provide the repository name for the build",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_REPO"),
				cli.EnvVar("BUILD_SUMMARY_REPO"),
				cli.EnvVar("VELA_REPO_NAME"),
				cli.File("/vela/parameters/build-summary/repo"),
				cli.File("/vela/parameters/build-summary/repo"),
			),
		},
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(_ context.Context, c *cli.Command) error {
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
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/build-summary",
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
			AppName:    c.Name,
			AppVersion: c.Version,
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
