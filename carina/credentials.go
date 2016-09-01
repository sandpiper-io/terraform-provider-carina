package carina

import (
	"os"
	"path/filepath"
	"io/ioutil"

	"github.com/getcarina/libcarina"
)

const defaultDotDir = ".carina"
const defaultNonDotDir = "carina"
const xdgDataHomeEnvVar = "XDG_DATA_HOME"
const clusterDirName = "clusters"

// CredentialsBaseDirEnvVar environment variable name for where credentials are downloaded to by default
const CredentialsBaseDirEnvVar = "CARINA_CREDENTIALS_DIR"

// CarinaHomeDirEnvVar is the environment variable name for carina data, config, etc.
const CarinaHomeDirEnvVar = "CARINA_HOME"


// CarinaCredentialsBaseDir get the current base directory for carina credentials
func CarinaCredentialsBaseDir() (string, error) {
	if os.Getenv(CarinaHomeDirEnvVar) != "" {
		return os.Getenv(CarinaHomeDirEnvVar), nil
	}
	if os.Getenv(CredentialsBaseDirEnvVar) != "" {
		return os.Getenv(CredentialsBaseDirEnvVar), nil
	}

	// Support XDG
	if os.Getenv(xdgDataHomeEnvVar) != "" {
		return filepath.Join(os.Getenv(xdgDataHomeEnvVar), defaultNonDotDir), nil
	}

	homeDir, err := userHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, defaultDotDir), nil
}

func CarinaStoreClusterCredentials(creds *libcarina.Credentials, username, clusterName string) (string, error) {
	baseDir, err := CarinaCredentialsBaseDir()
	if err != nil {
		return "", err
	}
	clusterDir := filepath.Join(baseDir,
		clusterDirName,
		username,
		clusterName,
	)
	err = os.MkdirAll(clusterDir, 0777)
	if err != nil {
		return "", err
	}

	for fname, b := range creds.Files {
		p := filepath.Join(clusterDir, fname)
		err = ioutil.WriteFile(p, b, 0600)
		if err != nil {
			return "", err
		}
	}

	return clusterDir, err
}

