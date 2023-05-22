package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/emrebdr/dirgod-core/constants"
	"github.com/emrebdr/dirgod-core/models"
)

func Contains(array []any, value any) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func SetWorkingMode(workingMode string) (models.Options, error) {
	if workingMode != "" {
		switch workingMode {
		case "force":
			return models.Force, nil
		case "strict":
			return models.Strict, nil
		case "default":
			return models.Default, nil
		default:
			return -1, errors.New("unknown working mode")
		}
	}

	return models.Default, nil
}

func CreateDirgodTempFolder() error {
	if _, err := os.Stat(constants.EnaTmp); os.IsNotExist(err) {
		return os.MkdirAll(constants.EnaTmp, os.FileMode(0755))
	}

	return nil
}

func GenerateId() string {
	time := time.Now().UnixNano()
	id := fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%d", time))))
	return id
}
