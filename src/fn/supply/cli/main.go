package main

import (
	"io"
	"io/ioutil"
	"nodejs/cache"
	_ "nodejs/hooks"
	"nodejs/npm"
	"nodejs/supply"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudfoundry/libbuildpack"
)

func main() {
	logfile, err := ioutil.TempFile("", "riff.nodejs-fn-buildpack.supply")
	defer logfile.Close()
	if err != nil {
		logger := libbuildpack.NewLogger(os.Stdout)
		logger.Error("Unable to create log file: %s", err.Error())
		os.Exit(8)
	}

	stdout := io.MultiWriter(os.Stdout, logfile)
	logger := libbuildpack.NewLogger(stdout)

	buildpackDir, err := libbuildpack.GetBuildpackDir()
	if err != nil {
		logger.Error("Unable to determine buildpack directory: %s", err.Error())
		os.Exit(9)
	}

	manifest, err := libbuildpack.NewManifest(buildpackDir, logger, time.Now())
	if err != nil {
		logger.Error("Unable to load buildpack manifest: %s", err.Error())
		os.Exit(10)
	}

	stager := libbuildpack.NewStager(os.Args[1:], logger, manifest)
	if err := stager.CheckBuildpackValid(); err != nil {
		os.Exit(11)
	}

	if err := manifest.ApplyOverride(stager.DepsDir()); err != nil {
		logger.Error("Unable to apply override.yml files: %s", err)
		os.Exit(17)
	}

	if err = libbuildpack.RunBeforeCompile(stager); err != nil {
		logger.Error("Before Compile: %s", err.Error())
		os.Exit(12)
	}

	if err = stager.SetStagingEnvironment(); err != nil {
		logger.Error("Unable to setup environment variables: %s", err.Error())
		os.Exit(13)
	}

	s := supply.Supplier{
		Logfile: logfile,
		Stager:  stager,
		NPM: &npm.NPM{
			BuildDir: filepath.Join(stager.DepDir(), "invoker"),
			Command:  &libbuildpack.Command{},
			Log:      logger,
		},
		Manifest: manifest,
		Log:      logger,
		Command:  &libbuildpack.Command{},
		Cache: &cache.Cache{
			Stager:  stager,
			Command: &libbuildpack.Command{},
			Log:     logger,
		},
	}

	err = supply.Run(&s)
	if err != nil {
		os.Exit(14)
	}

	if err := stager.WriteConfigYml(nil); err != nil {
		logger.Error("Error writing config.yml: %s", err.Error())
		os.Exit(15)
	}
}
