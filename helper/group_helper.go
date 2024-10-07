package helper

import (
	"fmt"
	"github.com/holypvp/primal/service"
)

func HighestThan(src, dest string) (bool, error) {
	srcGroup := service.Groups().LookupById(src)
	if srcGroup == nil {
		return false, fmt.Errorf("an error occurred while looking your highest group: %s", src)
	}

	destGroup := service.Groups().LookupById(dest)
	if destGroup == nil {
		return false, fmt.Errorf("an error occurred while looking the highest group of the target: %s", dest)
	}

	return srcGroup.Weight() > destGroup.Weight(), nil
}
