package operations

import "github.com/emrebdr/dirgod-core/models"

func DecideErrorOutput(options *models.OperationOptions, result *OperationResult, err error) {
	if options.WorkingMode == models.Force {
		result.Completed = true
		result.Err = err
	} else if options.WorkingMode == models.Default {
		// TODO: check for dependency tree
		result.Completed = false
		result.Err = err
	} else {
		result.Completed = false
		result.Err = err
	}
}
