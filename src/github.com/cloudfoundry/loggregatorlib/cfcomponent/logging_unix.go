// +build !windows,!plan9

package cfcomponent

import (
	"github.com/cloudfoundry/gosteno"
)

func GetNewSyslogSink(namespace string) *gosteno.Syslog {
	return gosteno.NewSyslogSink(namespace)
}
