package anomalies

import (
	"github.com/jagandecapri/vision/tree"
	"github.com/jagandecapri/vision/process"
)

type AnomaliesInterface interface{
	GetChannel(subspace_key [2]string) chan process.DissimilarityVector
	WaitOnChannels() chan bool
}