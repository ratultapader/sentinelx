package detection

import (
	"fmt"
	"time"
)

func generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}