package anomalies

import "github.com/jagandecapri/vision/tree"

type AnomaliesInterface interface{
	New() AnomaliesInterface
	GetChannel([2]string) chan tree.Subspace
	WaitOnChannels(chan bool) chan bool
}