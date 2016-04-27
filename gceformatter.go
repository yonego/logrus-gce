package logrusgce

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
)

type severity string

const (
	logrusToCallerSkip = 4
)

const (
	severityDEBUG     severity = "DEBUG"
	severityINFO      severity = "INFO"
	severityNOTICE    severity = "NOTICE"
	severityWARNING   severity = "WARNING"
	severityERROR     severity = "ERROR"
	severityCRITICAL  severity = "CRITICAL"
	severityALERT     severity = "ALERT"
	severityEMERGENCY severity = "EMERGENCY"
)

var (
	levelsLogrusToGCE = map[logrus.Level]severity{
		logrus.DebugLevel: severityDEBUG,
		logrus.InfoLevel:  severityINFO,
		logrus.WarnLevel:  severityWARNING,
		logrus.ErrorLevel: severityERROR,
		logrus.FatalLevel: severityCRITICAL,
		logrus.PanicLevel: severityALERT,
	}
)

type sourceLocation struct {
	File         string `json:"file"`
	Line         int    `json:"line"`
	FunctionName string `json:"functionName"`
}

type GCEFormatter struct {
	withSourceInfo bool
}

func NewGCEFormatter(withSourceInfo bool) *GCEFormatter {
	return &GCEFormatter{withSourceInfo: withSourceInfo}
}

func (f *GCEFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	data["time"] = entry.Time.Format(time.RFC3339Nano)
	data["severity"] = string(levelsLogrusToGCE[entry.Level])
	data["logMessage"] = entry.Message

	if f.withSourceInfo == true {
		pc, file, line, ok := runtime.Caller(logrusToCallerSkip)
		if ok {
			f := runtime.FuncForPC(pc)
			data["sourceLocation"] = sourceLocation{
				File:         file,
				Line:         line,
				FunctionName: f.Name(),
			}
		}
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
