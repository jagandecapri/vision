package preprocess

import (
	"github.com/jagandecapri/vision/utils"
)

type PacketAcc []PacketFeature

func (pa PacketAcc) ExtractDeltaPacketFeature() (map[string]float64){
	nbSYN, nbACK, nbRST, nbFIN, nbCWR, nbURG, totalPktSize, totalTTL  := 0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0
	srcPorts := []float64{}
	dstPorts := []float64{}
	srcIPs := []string{}
	dstIPs := []string{}

	for _, fp := range pa{
		if fp.SrcPort != 0 {
			srcPorts = append(srcPorts, float64(fp.SrcPort))
		}
		if fp.DstPort != 0 {
			dstPorts = append(dstPorts, float64(fp.DstPort))
		}
		if fp.SrcIP != nil {
			srcIPs = append(srcIPs, string(fp.SrcIP))
		}
		if fp.DstIP != nil {
			dstIPs = append(dstIPs, string(fp.DstIP))
		}
		if fp.SYN == true{
			nbSYN++
		}
		if fp.ACK == true{
			nbACK++
		}
		if fp.RST == true{
			nbRST++
		}
		if fp.ACK == true{
			nbFIN++
		}
		if fp.CWR == true{
			nbCWR++
		}
		if fp.URG == true{
			nbURG++
		}
		totalPktSize += float64(fp.length)
		totalTTL += float64(fp.TTL)
	}

	nbPacket := float64(len(pa))
	nbSrcPort := float64(len(utils.UniqFloat64(srcPorts)))
	nbDstPort := float64(len(utils.UniqFloat64(dstPorts)))
	nbSrcs := float64(len(utils.UniqString(srcIPs)))
	nbDsts := float64(len(utils.UniqString(dstIPs)))
	perSyn := nbSYN/nbPacket*100
	perAck := nbACK/nbPacket*100
	perRST := nbRST/nbPacket*100
	perFIN := nbFIN/nbPacket*100
	perCWR := nbCWR/nbPacket*100
	perURG := nbURG/nbPacket*100
	avgPktSize := totalPktSize/nbPacket
	meanTTL := totalTTL/nbPacket
	x := map[string]float64{
		"nbPacket": nbPacket,
		"nbSrcPort": nbSrcPort,
		"nbDstPort": nbDstPort,
		"nbSrcs": nbSrcs,
		"nbDsts": nbDsts,
		"perSyn": perSyn,
		"perAck": perAck,
		"perRST": perRST,
		"perFIN": perFIN,
		"perCWR": perCWR,
		"perURG": perURG,
		"avgPktSize": avgPktSize,
		"meanTTL": meanTTL,
	}

	return x
}