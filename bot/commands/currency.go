package commands

import (
	"fmt"
	"taxeer/service"
	"time"
)

func HandleCurrencyCommand(currency string) string {
	currentCurrencyRate, err := service.GetCurrencyAtDate(time.Now(), currency)
	if err != nil {
		return fmt.Sprintf("Ooops! Can't get today currency rate for %s, try again later:(", currency)
	}
	return fmt.Sprintf("%.4f", currentCurrencyRate)
}
