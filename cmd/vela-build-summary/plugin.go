// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// build arguments loaded for the plugin
	Build *Build
	// config arguments loaded for the plugin
	Config *Config
	// repo arguments loaded for the plugin
	Repo *Repo
}

// Exec formats and runs the commands for creating a summary of the build.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	logrus.Infof("creating client for %s", p.Config.Server)
	// create new Vela client from config configuration
	client, err := p.Config.New()
	if err != nil {
		return err
	}

	logrus.Infof("capturing build %s/%s/%d", p.Repo.Org, p.Repo.Name, p.Build.Number)
	// send API call to capture a build
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#BuildService.Get
	build, _, err := client.Build.Get(p.Repo.Org, p.Repo.Name, p.Build.Number)
	if err != nil {
		return err
	}

	// set the pagination options for list of resources
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#ListOptions
	opts := &vela.ListOptions{
		PerPage: 100,
	}

	logrus.Infof("capturing services for build %s/%s/%d", p.Repo.Org, p.Repo.Name, p.Build.Number)
	// send API call to capture a list of services
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#SvcService.GetAll
	services, _, err := client.Svc.GetAll(p.Repo.Org, p.Repo.Name, p.Build.Number, opts)
	if err != nil {
		return err
	}

	logrus.Infof("capturing steps for build %s/%s/%d", p.Repo.Org, p.Repo.Name, p.Build.Number)
	// send API call to capture a list of steps
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#StepService.GetAll
	steps, _, err := client.Step.GetAll(p.Repo.Org, p.Repo.Name, p.Build.Number, opts)
	if err != nil {
		return err
	}

	logrus.Infof("capturing logs for build %s/%s/%d", p.Repo.Org, p.Repo.Name, p.Build.Number)
	// send API call to capture a list of build logs
	//
	// https://pkg.go.dev/github.com/go-vela/sdk-go/vela?tab=doc#BuildService.GetLogs
	logs, _, err := client.Build.GetLogs(p.Repo.Org, p.Repo.Name, p.Build.Number, nil)
	if err != nil {
		return err
	}

	return table(build, logs, services, steps)
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate build configuration
	err := p.Build.Validate()
	if err != nil {
		return err
	}

	// validate config configuration
	err = p.Config.Validate()
	if err != nil {
		return err
	}

	// validate repo configuration
	err = p.Repo.Validate()
	if err != nil {
		return err
	}

	return nil
}
