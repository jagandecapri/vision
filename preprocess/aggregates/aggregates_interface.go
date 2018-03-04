package aggregates

import (
	"github.com/jagandecapri/vision/tree"
)

type AggInterface interface{
	GetKey() []tree.PointKey
	NbPacket() float64
	NbSrcPort() float64
	NbDstPort() float64
	NbSrcs() float64
	NbDsts() float64
	PerSYN() float64
	PerACK() float64
	PerRST() float64
	PerFIN() float64
	PerCWR() float64
	PerURG() float64
	AvgPktSize() float64
	MeanTTL() float64
	SimIPDst() float64
	SimIPSrc() float64
	PerICMP() float64
	PerICMPRed() float64
	PerICMPTime() float64
	PerICMPUnr() float64
	PerICMPOther() float64
}