package anomalies

type AnomaliesInterface interface{
	GetChannel([2]string) chan DissimilarityVector
	WaitOnChannels(chan struct{})
}