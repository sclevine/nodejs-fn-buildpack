package supply

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Cache interface {
	Initialize() error
	Restore() error
	Save() error
}

type Command interface {
	Execute(string, io.Writer, io.Writer, string, ...string) error
}

type Manifest interface {
	AllDependencyVersions(string) []string
	DefaultVersion(string) (libbuildpack.Dependency, error)
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type NPM interface {
	Build() error
	Rebuild() error
}

type Stager interface {
	BuildDir() string
	DepDir() string
	DepsIdx() string
	LinkDirectoryInDepDir(string, string) error
	WriteEnvFile(string, string) error
	WriteProfileD(string, string) error
	SetStagingEnvironment() error
}

type Supplier struct {
	Stager     Stager
	Manifest   Manifest
	Log        *libbuildpack.Logger
	Logfile    *os.File
	Command    Command
	NPMRebuild bool
	Cache      Cache
	NPM        NPM
}

func Run(s *Supplier) error {
	s.Log.BeginStep("Installing invoker")

	if err := s.InstallInvoker(); err != nil {
		s.Log.Error("Unable to install node: %s", err.Error())
		return err
	}

	if err := s.CreateDefaultEnv(); err != nil {
		s.Log.Error("Unable to setup default environment: %s", err.Error())
		return err
	}

	if err := s.Stager.SetStagingEnvironment(); err != nil {
		s.Log.Error("Unable to setup environment variables: %s", err.Error())
		os.Exit(11)
	}

	if err := s.Cache.Initialize(); err != nil {
		s.Log.Error("Unable to initialize cache: %s", err.Error())
		return err
	}

	if err := s.Cache.Restore(); err != nil {
		s.Log.Error("Unable to restore cache: %s", err.Error())
		return err
	}

	defer s.Logfile.Sync()

	if err := s.BuildDependencies(); err != nil {
		s.Log.Error("Unable to build dependencies: %s", err.Error())
		return err
	}

	if err := s.Cache.Save(); err != nil {
		s.Log.Error("Unable to save cache: %s", err.Error())
		return err
	}

	if err := s.Logfile.Sync(); err != nil {
		s.Log.Error(err.Error())
		return err
	}

	return nil
}

func (s *Supplier) BuildDependencies() error {
	s.Log.BeginStep("Building dependencies")

	if s.NPMRebuild {
		s.Log.Info("Prebuild detected (node_modules already exists)")
		if err := s.NPM.Rebuild(); err != nil {
			return err
		}
	} else {
		if err := s.NPM.Build(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Supplier) InstallInvoker(tempDir string) error {
	var dep libbuildpack.Dependency

	invokerInstallDir := filepath.Join(s.Stager.DepDir(), "invoker")

	if dep, err = s.Manifest.DefaultVersion("node-function-invoker"); err != nil {
		return err
	}

	if err := s.Manifest.InstallDependency(dep, tempDir); err != nil {
		return err
	}

	matches, err := filepath.Glob(filepath.Join(tempDir, "*"))
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errors.New("invoker not found in specified dependency")
	}

	return os.Rename(matches[0], invokerInstallDir)
}

func (s *Supplier) CreateDefaultEnv() error {
	scriptContents := `export HOST=0.0.0.0
export HTTP_PORT=${HTTP_PORT:8080}
export GRPC_PORT=${GRPC_PORT:10382}
export INVOKER_DIR="$DEPS_DIR/%s/invoker"
`
	return s.Stager.WriteProfileD("fn.sh", fmt.Sprintf(scriptContents, s.DepsIdx()))
}
