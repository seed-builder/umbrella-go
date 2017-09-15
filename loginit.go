package umbrella

import (
	"github.com/op/go-logging"
	"os"
	"umbrella/models"
	"umbrella/utilities"
)


var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	logfile := &utilities.LogFile{}
	backend1 := logging.NewLogBackend(logfile, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend3 := &utilities.DBackend{}
	l := new(models.EquipmentLog)
	backend3.Logger = l

	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.INFO, "")

	logging.SetBackend(backend1Leveled, backend2Formatter, backend3)
}