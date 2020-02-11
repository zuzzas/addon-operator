package log

import (
	"time"

	"github.com/flant/addon-operator/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func MeasureTimeToLog(fn func(), description string, logLabels map[string]string) {
	start := time.Now().UnixNano()
	fn()
	passed := time.Duration(time.Now().UnixNano() - start).Truncate(time.Microsecond).String()
	log.WithFields(utils.LabelsToLogFields()).Infof("%s took %s", description, passed)
}
