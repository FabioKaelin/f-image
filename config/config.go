package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fabiokaelin/fcommon/pkg/values"
	"github.com/fabiokaelin/ferror"
	"github.com/joho/godotenv"
)

var (
	// GinMode is the mode of the gin server (Default: debug)
	GinMode string
	// NotificationID is the notification id for the push notifications (Default: "")
	NotificationID string
	// FVersion is the version of the backend
	FVersion string
	// JsonLogs is the flag for json logs
	JsonLogs bool
)

func getString(key string) (string, ferror.FError) {
	value := os.Getenv(key)
	if value == "" {
		ferr := ferror.New(fmt.Sprintf("key '%s' not found", key))
		ferr.SetLayer("config")
		ferr.SetKind("read string from env")
		return "", ferr
	}
	return value, nil
}

func getBool(key string) (bool, ferror.FError) {
	value := os.Getenv(key)
	if value == "" {
		ferr := ferror.New(fmt.Sprintf("key '%s' not found", key))
		ferr.SetLayer("config")
		ferr.SetKind("read bool from env")
		return false, ferr
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		ferr := ferror.FromError(err)
		ferr.SetLayer("config")
		ferr.SetKind("parse bool from env")
		return false, ferr
	}
	return boolVal, nil
}

// Load loads the config from the .env or .env.test file
func Load(environment string) ferror.FError {
	if environment == "test" || os.Getenv("GIN_MODE") == "test" {
		godotenv.Load(".env.test")
	} else {
		godotenv.Load(".env")
	}

	ginmode, ferr := getString("GIN_MODE")
	if ferr != nil {
		return ferr
	}
	if ginmode != "debug" && ginmode != "release" {
		ginmode = "debug"
	}
	GinMode = ginmode

	NotificationID, ferr = getString("NOTIFICATION_ID")
	if ferr != nil {
		return ferr
	}

	FVersion, ferr = getString("F_VERSION")
	if ferr != nil {
		return ferr
	}

	JsonLogs, ferr = getBool("JSON_LOGS")
	if ferr != nil {
		return ferr
	}

	values.V = values.Values{
		GinMode:        GinMode,
		JsonLogs:       JsonLogs,
		NotificationID: NotificationID,
		FVersion:       FVersion,
	}

	return nil
}
