package anomalies

import "sync"

type AnomaliesInterface interface{
	GetChannel([2]string) chan DissimilarityVectorContainer
	WaitOnChannels(*sync.WaitGroup)
}