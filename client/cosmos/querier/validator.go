package querier

import (
	"fmt"
	"github.com/informalsystems/stakooler/client/cosmos/model"
	"github.com/schollz/progressbar/v3"
)

func LoadValidatorStats(validator *model.Validator, bar *progressbar.ProgressBar) error {

	fmt.Printf("Validator: %s\n", validator)
	//TODO: Implement logic to fetch validator statistics
	return nil
}
