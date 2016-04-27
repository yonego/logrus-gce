# logrus-gce
Google Cloud Engine Logging Formatter for Logrus

## How to use
```golang
package main

import (
	log "github.com/Sirupsen/logrus"
    logrusgce "github.com/znly/logrus-gce"
)

func main() {
    log.SetFormatter(logrusgce.NewGCEFormatter(true))
    log.WithField("myfield", "myvalue").Info("hey")
}
```