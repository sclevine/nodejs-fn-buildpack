package finalize

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Manifest interface {
	RootDir() string
}

type Stager interface {
	BuildDir() string
	DepDir() string
	DepsIdx() string
}

type Finalizer struct {
	Stager  Stager
	Log     *libbuildpack.Logger
	Logfile *os.File
}

func Run(f *Finalizer) error {
	if err := f.CopyProfileScripts(); err != nil {
		f.Log.Error("Unable to copy profile.d scripts: %s", err.Error())
		return err
	}

	ioutil.WriteFile(filepath.Join(f.Stager.

	if err := f.Logfile.Sync(); err != nil {
		f.Log.Error(err.Error())
		return err
	}

	return nil
}

func (f *Finalizer) CopyProfileScripts() error {
	profiledDir := filepath.Join(f.Stager.DepDir(), "profile.d")
	if err := os.MkdirAll(profiledDir, 0755); err != nil {
		return err
	}

	for _, fi := range files {
		if err := libbuildpack.CopyFile(filepath.Join(path, fi.Name()), filepath.Join(profiledDir, fi.Name())); err != nil {
			return err
		}
	}
	return nil
}
