package anomalies

import (

	"github.com/jagandecapri/vision/process"
)

type AnomaliesInterface interface{
	GetChannel([2]string) chan process.DissimilarityVector
	WaitOnChannels(chan struct{})
}