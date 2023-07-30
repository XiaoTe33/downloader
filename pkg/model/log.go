package model

import (
	"downloader/pkg/myLog"
	"github.com/sirupsen/logrus"
)

func init() {
	myLog.Log.SetLevel(logrus.TraceLevel)
}
