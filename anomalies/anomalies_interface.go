package anomalies

type AnomaliesInterface interface{
	GetChannel([2]string) chan DissimilarityVectorContainer
	WaitOnChannels(chan struct{})
}