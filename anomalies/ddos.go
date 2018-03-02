package anomalies

import "github.com/jagandecapri/vision/tree"

type DDOS struct{
	Channels map[[2]string]chan tree.Subspace
}

func (d *DDOS) GetChannel(subspace_key [2]string) chan tree.Subspace {
	return d.Channels[subspace_key]
}

func (DDOS) WaitOnChannels() chan bool {
	return make(chan bool)
}

func New() DDOS{
	return DDOS{Channels: map[[2]string]chan tree.Subspace{
		[2]string{"nbSrcs", "avgPktSize"}: make(chan tree.Subspace),
		[2]string{"perICMP", "perSYN"}: make(chan tree.Subspace),
		[2]string{"nbSrcPort", "perICMP"}: make(chan tree.Subspace),
	}}
}