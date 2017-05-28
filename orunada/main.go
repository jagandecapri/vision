package main

import (
	"net"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
	"time"
	"fmt"
	"github.com/cockroachdb/apd"
	"strconv"
)


type PacketFeature struct{
	SrcIP, DstIP net.IP
	SrcPort, DstPort layers.TCPPort
	SYN, ACK, RST, FIN, CWR, URG bool
	length uint16
	TTL uint8
}
func extractFeature(packet gopacket.Packet) PacketFeature{
	packet_feature := PacketFeature{}
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		// ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	}

	// Let's see if the packet is IPv4 (even though the ether type told us)
	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		packet_feature.SrcIP = ipv4.SrcIP
		packet_feature.DstIP = ipv4.DstIP
		packet_feature.TTL = ipv4.TTL
	}
	// Let's see if the packet is IPv6 (even though the ether type told us)
	ipv6Layer := packet.Layer(layers.LayerTypeIPv4)
	if ipv6Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv6, _ := ipv6Layer.(*layers.IPv4)
		packet_feature.SrcIP = ipv6.SrcIP
		packet_feature.DstIP = ipv6.DstIP
		packet_feature.TTL = ipv6.TTL
	}

	// Let's see if the packet is TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		//fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		packet_feature.SrcPort = tcp.SrcPort
		packet_feature.DstPort = tcp.DstPort
		packet_feature.SYN = tcp.SYN
		packet_feature.ACK = tcp.ACK
		packet_feature.RST = tcp.RST
		packet_feature.FIN = tcp.FIN
		packet_feature.CWR = tcp.CWR
		packet_feature.URG = tcp.URG
	}

	// Let's see if the packet is icmp4
	icmp4Layer := packet.Layer(layers.LayerTypeICMPv4)
	if icmp4Layer != nil {
		//fmt.Println("ICMP4 Layer Detected")
		//icmp4, _ := icmp4Layer.(*layers.ICMPv4)
		//fmt.Println(icmp4.TypeCode.String())
	}

	// Let's see if the packet is icmp6
	icmp6Layer := packet.Layer(layers.LayerTypeICMPv6)
	if icmp6Layer != nil {
		//fmt.Println("ICMP6 Layer Detected")
		//icmp6, _ := icmp6Layer.(*layers.ICMPv6)
		//fmt.Println(icmp6.TypeCode.String())
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		//fmt.Println("Error decoding some part of the packet:", err)
	}
	return packet_feature
}

var delta_t time.Duration = 300 * time.Millisecond
var window time.Duration = 15 * time.Second
var window_arr_len = int(window.Seconds()/delta_t.Seconds())
var time_counter time.Time

type PacketData struct{
	data gopacket.Packet
	metadata gopacket.CaptureInfo
}

type PacketAcc []PacketFeature

func UniqFloat64(input []float64) []float64 {
	u := make([]float64, 0, len(input))
	m := make(map[float64]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

func UniqString(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

func ExtractDeltaPacketFeature(feature_packets PacketAcc) []float64{
	nbSYN, nbACK, nbRST, nbFIN, nbCWR, nbURG, totalPktSize, totalTTL  := 0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0
	srcPorts := []float64{}
	dstPorts := []float64{}
	srcIPs := []string{}
	dstIPs := []string{}

	for _, fp := range feature_packets{
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
	x := []float64{}
	nbPacket := float64(len(feature_packets))
	nbSrcPort := float64(len(UniqFloat64(srcPorts)))
	nbDstPort := float64(len(UniqFloat64(dstPorts)))
	nbSrcs := float64(len(UniqString(srcIPs)))
	nbDsts := float64(len(UniqString(dstIPs)))
	perSyn := nbSYN/nbPacket*100
	perAck := nbACK/nbPacket*100
	perRST := nbRST/nbPacket*100
	perFIN := nbFIN/nbPacket*100
	perCWR := nbCWR/nbPacket*100
	perURG := nbURG/nbPacket*100
	avgPktSize := totalPktSize/nbPacket
	meanTTL := totalTTL/nbPacket
	x = append(x, nbPacket, nbSrcPort, nbDstPort, nbSrcs, nbDsts, perSyn, perAck, perRST, perFIN, perCWR, perURG, avgPktSize, meanTTL)
	return x
}

func WindowTimeSlide(ch chan PacketData, acc chan PacketAcc, quit chan int){
	var tmp PacketAcc
	initializePacketAcc := func(){
		tmp = []PacketFeature{}
	}
	for{
		select{
		case pd := <- ch:
			packet_time := pd.metadata.Timestamp
			if time_counter.IsZero() {
				//fmt.Println("Initialize Time")
				time_counter = packet_time
				pf := extractFeature(pd.data)
				initializePacketAcc()
				tmp = append(tmp, pf)
			} else if packet_time.Before(time_counter.Add(delta_t)) || packet_time.Equal(time_counter.Add(delta_t)) {
				//fmt.Println("packet_time <= time_counter")
				pf := extractFeature(pd.data)
				tmp = append(tmp, pf)
			} else if packet_time.After(time_counter.Add(delta_t)) {
				//fmt.Println("packet_time > time_counter")
				acc <- tmp
				initializePacketAcc()
				time_counter = time_counter.Add(delta_t)
			}
		case <-quit:
			fmt.Println("No Data")
			close(ch)
			close(acc)
		}
	}
}

type DimMinMax struct{
	Min, Max float64
	Range float64
	ID int
}

func Normalize(mat [][]float64) ([][]float64, []DimMinMax){
	rows := len(mat)
	cols := len(mat[0])

	dim_min_max := []DimMinMax{}
	for c := 0; c < cols; c++ {
		min := mat[0][c]
		max := mat[0][c]
		for j := 0; j < rows; j++{
			val := mat[j][c]
			if val < min{
				min = val
			} else if  val > max{
				max = val
			}
		}
		range_ := max - min
		dim_min_max = append(dim_min_max, DimMinMax{min, max, range_, c})
	}

	for i := 0; i < rows; i++{
		for j := 0; j < cols; j++{
			col_min := dim_min_max[j].Min
			col_max := dim_min_max[j].Max
			elem := mat[i][j]
			if col_min == 0 && col_max == 0{
				mat[i][j] = elem
			} else {
				mat[i][j] = (elem - col_min)/(col_max - col_min)
			}

		}
	}
	return mat, dim_min_max
}
func UpdateFS(acc chan PacketAcc){
	base_matrix := [][]float64{}
	for{
		select{
		case packet_acc := <- acc:
			fmt.Print(".")
			x := ExtractDeltaPacketFeature(packet_acc)
			if len(base_matrix) < window_arr_len{
				base_matrix = append(base_matrix, x)
			} else if len(base_matrix) == window_arr_len{
				base_matrix = append(base_matrix, x)
				//TODO: normalization need to be done here, x_old and x_new will be what here?
			} else {
				base_matrix = append(base_matrix, x)
				norm_mat, dim_min_max := Normalize(base_matrix)
				x_old, x_update, x_new := [][]float64{norm_mat[0]}, norm_mat[1:len(norm_mat)-2], [][]float64{norm_mat[len(norm_mat)-1]}
				build2Dgrid(x_old, x_update, x_new, dim_min_max)
				//base_matrix = base_matrix[1:]
			}

		}
	}
}

type Point struct{
	vec []float64
	unit_id int
}

type Interval struct{
	range_ []float64
}
type Unit struct{
	id int
	intervals []Interval
	density int
	points []Point
}

type Grid struct{
	units []Unit
}

func (g *Grid) intersect(vec []float64) *Unit{
		for i := 0; i < len(g.units); i++{
			unit := &g.units[i]
			inside_interval_ctr := false
			lower_bound_x := unit.intervals[0].range_[0]
			upper_bound_x := unit.intervals[0].range_[1]
			lower_bound_y := unit.intervals[1].range_[0]
			upper_bound_y := unit.intervals[1].range_[1]
			if i == len(g.units) - 1{
				if vec[0] >= lower_bound_x && vec[0] <= upper_bound_x && vec[1] >= lower_bound_y && vec[1] <= upper_bound_y {
					inside_interval_ctr = true
				}
			} else {
				if vec[0] >= lower_bound_x && vec[0] < upper_bound_x && vec[1] >= lower_bound_y && vec[1] < upper_bound_y {
					inside_interval_ctr = true
				}
			}
			if inside_interval_ctr == true{
				fmt.Println("Intersected", unit)
				return unit
			}
		}
	return nil
}

func build2Dgrid(x_old, x_update, x_new [][]float64, dim_min_max []DimMinMax){
	grid := Grid{}
	unit_id := 0
	ctx := apd.BaseContext.WithPrecision(6)
	dim_x := []Interval{}
	dim_j := []Interval{}
	for j := 0; j < len(dim_min_max); j++ {
		interval_l, _, _ := apd.NewFromString("0.1")
		incr, _, _ := apd.NewFromString("0.0")
		for i := 0; i < 10; i++ {
			interval := Interval{}
			range_, _, _:= apd.NewFromString(strconv.FormatFloat(dim_min_max[j].Range, 'f', -1, 64))
			// interval_l here is the same number
			min, _, _ := apd.NewFromString(strconv.FormatFloat(dim_min_max[j].Min, 'f', -1, 64))
			lb := new(apd.Decimal)
			ctx.Mul(lb,incr,range_)
			ctx.Add(lb, lb, min)
			lower_bound, _ := lb.Float64()
			ub := new(apd.Decimal)
			ctx.Add(ub, incr, interval_l)
			ctx.Mul(ub, ub, range_)
			ctx.Add(ub, ub, min)
			upper_bound, _ := ub.Float64()
			interval.range_ = []float64{lower_bound, upper_bound}
			if j == 0 {
				dim_x = append(dim_x, interval)
			} else {
				dim_j = append(dim_j, interval)
			}
			ctx.Add(incr, incr, interval_l)
		}

	}

	for i := 0; i < len(dim_x); i++{
		for j := 0; j < len(dim_j); j++{
			unit := Unit{}
			unit_id += 1
			unit.id = unit_id
			unit.intervals = append(unit.intervals, dim_x[i], dim_j[j])
			grid.units = append(grid.units, unit)
		}
	}

	fmt.Println(grid.units)
	for _, elem := range x_old{
		fmt.Println(elem)
		grid.intersect(elem)
		//fmt.Println(*val)
	}
}
func main() {
	//handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	//ch := make(chan PacketData)
	//acc := make(chan PacketAcc)
	//quit := make(chan int)
	//go WindowTimeSlide(ch, acc, quit)
	//go UpdateFS(acc)
	//
	//if(err != nil){
	//	log.Fatal(err)
	//}
	//
	//for {
	//	data, ci, err := handleRead.ReadPacketData()
	//	if err != nil && err != io.EOF {
	//		quit <- 0
	//		log.Fatal(err)
	//	} else if err == io.EOF {
	//		quit <- 0
	//		break
	//	} else {
	//		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	//		ch <- PacketData{packet, ci}
	//	}
	//}
	dim_min_max := []DimMinMax{{
		Min: 1.003,
		Max: 10.0004,
		ID: 1,
	},{
		Min: 1.003,
		Max: 10.0004,
		ID: 1,
	}}
	x_old := [][]float64{{1.12344123141234123,2.000012341234123234123},{1.12344123141234123,2.000012341234123234123}}
	dim_min_max[0].Range = dim_min_max[0].Max - dim_min_max[0].Min
	dim_min_max[1].Range = dim_min_max[1].Max - dim_min_max[1].Min
	build2Dgrid(x_old, [][]float64{}, [][]float64{}, dim_min_max)

}