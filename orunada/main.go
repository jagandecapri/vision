package main

import (
	"net"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
	"time"
	"fmt"
	//"github.com/cockroachdb/apd"
	//"strconv"
	"sort"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"os"
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

func ExtractDeltaPacketFeature(feature_packets PacketAcc) (map[string]float64, []string){
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

	sorter := []string{}
	for k := range x{
		sorter = append(sorter, k)
	}
	sort.Strings(sorter)
	return x, sorter
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
}

func Normalize(mat []Point, sorter []string) ([]Point, map[string]DimMinMax){
	rows := len(mat)
	//cols := len(mat[0])

	dim_min_max := map[string]DimMinMax{}
	for _,c := range sorter {
		min := mat[0].vec[c]
		max := mat[0].vec[c]
		for j := 0; j < rows; j++{
			val := mat[j].vec[c]
			if val < min{
				min = val
			} else if  val > max{
				max = val
			}
		}
		range_ := max - min
		dim_min_max[c] = DimMinMax{min, max, range_}
	}

	for i := 0; i < rows; i++{
		for _, c := range sorter{
			col_min := dim_min_max[c].Min
			col_max := dim_min_max[c].Max
			elem := mat[i].vec[c]
			if col_min == 0 && col_max == 0{
				mat[i].norm_vec[c] = elem
			} else {
				mat[i].norm_vec[c] = (elem - col_min)/(col_max - col_min)
			}
		}
	}

	//Assign normalized min-max
	for _,c := range sorter{
		norm_col_min := 0.0
		tmp := dim_min_max[c]
		tmp.Min = norm_col_min
		norm_col_max := 1.0
		tmp.Max = norm_col_max
		dim_min_max[c] = tmp
	}

	return mat, dim_min_max
}

func comb(n, m int, emit func([]int)) {
	s := make([]int, m)
	last := m - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			s[i] = j
			if i == last {
				emit(s)
			} else {
				rc(i+1, j+1)
			}
		}
		return
	}
	rc(0, 0)
}

func getSubspaceKey(sorter []string, feature_cnt int) [][]string {
	all := [][]string{}
	comb(len(sorter), feature_cnt, func (c []int){
		tmp := []string{}
		for _, v := range c {
			tmp = append(tmp, sorter[v])
		}
		all = append(all, tmp)
	})
	return all
}

func getSubspace(subspace_keys [][]string, mat []Point) map[[2]string][]Point{
	subspaces := map[[2]string][]Point{}
	for _, subspace_k := range subspace_keys{
		key := [2]string{}
		copy(key[:], subspace_k)
		subspace := []Point{}
		for _, p:= range mat{
			sub_point := Point{id: p.id}
			tmp := make(map[string]float64)
			for i := 0; i < len(subspace_k); i++{
				key := subspace_k[i]
				tmp[key] = p.norm_vec[key]
			}
			sub_point.norm_vec = tmp
			subspace = append(subspace, sub_point)
		}
		subspaces[key] = subspace
	}
	return subspaces
}

var point_ctr int = 0

func UpdateFS(acc chan PacketAcc){
	base_matrix := []Point{}
	for{
		select{
		case packet_acc := <- acc:
			fmt.Print(".")
			x, sorter := ExtractDeltaPacketFeature(packet_acc)
			point_ctr += 1
			p := Point{id: point_ctr, vec: x, norm_vec: make(map[string]float64)}
			if len(base_matrix) < window_arr_len - 1{
				base_matrix = append(base_matrix, p)
			} else if len(base_matrix) == window_arr_len - 1{
				base_matrix = append(base_matrix, p)
				//TODO: normalization need to be done here, x_old and x_new will be what here?
			} else {
				base_matrix = append(base_matrix, p)
				norm_mat, dim_min_max := Normalize(base_matrix, sorter)
				subspace_keys := getSubspaceKey(sorter, 2)
				subspaces := getSubspace(subspace_keys, norm_mat)
				for keys, subspace := range subspaces{
					dims := map[string]DimMinMax{}
					for _, key := range keys{
						dims[key] = dim_min_max[key]
					}
					x_old, x_update, x_new := []Point{subspace[0]}, subspace[1:len(subspace)-2], []Point{subspace[len(subspace)-1]}
					build2Dgrid(x_old, x_update, x_new, dims)
				}
				base_matrix = base_matrix[1:]
				os.Exit(2)
			}

		}
	}
}

type Point struct{
	id int
	vec map[string]float64
	norm_vec map[string]float64
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

//func (g *Grid) intersect(p Point) *Unit{
//	vec := p.vec
//	ctx := apd.BaseContext.WithPrecision(6)
//	for i := 0; i < len(g.units); i++{
//			unit := &g.units[i]
//			inside_interval_ctr := false
//			lower_bound_x := unit.intervals[0].range_[0]
//			upper_bound_x := unit.intervals[0].range_[1]
//			lower_bound_y := unit.intervals[1].range_[0]
//			upper_bound_y := unit.intervals[1].range_[1]
//
//			lb_x, _ := new(apd.Decimal).SetFloat64(lower_bound_x)
//			ub_x, _ := new(apd.Decimal).SetFloat64(upper_bound_x)
//			lb_y, _ := new(apd.Decimal).SetFloat64(lower_bound_y)
//			ub_y, _ := new(apd.Decimal).SetFloat64(upper_bound_y)
//
//			vec_0, _ := new(apd.Decimal).SetFloat64(vec[0])
//			vec_1, _ := new(apd.Decimal).SetFloat64(vec[1])
//
//			cmp_vec_0_lb_x := new(apd.Decimal)
//			cmp_vec_0_ub_x := new(apd.Decimal)
//			ctx.Cmp(cmp_vec_0_lb_x, vec_0, lb_x)
//			ctx.Cmp(cmp_vec_0_ub_x, vec_0, ub_x)
//
//			cmp_vec_1_lb_y := new(apd.Decimal)
//			cmp_vec_1_ub_y := new(apd.Decimal)
//			ctx.Cmp(cmp_vec_1_lb_y, vec_1, lb_y)
//			ctx.Cmp(cmp_vec_1_ub_y, vec_1, ub_y)
//
//			int_cmp_vec_0_lb_x, _ := cmp_vec_0_lb_x.Int64()
//			int_cmp_vec_0_ub_x, _ := cmp_vec_0_ub_x.Int64()
//			int_cmp_vec_1_lb_y, _ := cmp_vec_1_lb_y.Int64()
//			int_cmp_vec_1_ub_y, _ := cmp_vec_1_ub_y.Int64()
//
//			if i == len(g.units) - 1{
//				if (int_cmp_vec_0_lb_x == 1 || int_cmp_vec_0_lb_x == 0) && (int_cmp_vec_0_ub_x == -1 || int_cmp_vec_0_ub_x == 0) && (int_cmp_vec_1_lb_y == 1 || int_cmp_vec_1_lb_y == 0) && (int_cmp_vec_1_ub_y == -1 || int_cmp_vec_1_ub_y == 0){
//					inside_interval_ctr = true
//				}
//			} else {
//				if (int_cmp_vec_0_lb_x == 1 || int_cmp_vec_0_lb_x == 0) && (int_cmp_vec_0_ub_x == -1) && (int_cmp_vec_1_lb_y == 1 || int_cmp_vec_1_lb_y == 0) && (int_cmp_vec_1_ub_y == -1){
//					inside_interval_ctr = true
//				}
//			}
//
//			if inside_interval_ctr == true{
//				fmt.Println("Intersected", unit)
//				return unit
//			}
//		}
//	return nil
//}

func build2Dgrid(x_old, x_update, x_new []Point, dim_min_max map[string]DimMinMax){
	grid := Grid{}
	unit_id := 0
	ctx := apd.BaseContext.WithPrecision(6)
	dim_x := []Interval{}
	dim_j := []Interval{}
	fmt.Println(len(dim_min_max))
	axis := 0
	for _, dim := range dim_min_max {
		interval_l, _ := new(apd.Decimal).SetFloat64(0.1)
		incr, _ := new(apd.Decimal).SetFloat64(0.0)
		range_, _ := new(apd.Decimal).SetFloat64(1.0)
		for i := 0; i < 10; i++ {
			interval := Interval{}
			//TODO: Range here only need to be 1
			//range_, _, _:= apd.NewFromString(strconv.FormatFloat(dim.Range, 'f', -1, 64))
			// interval_l here is the same number
			min, _, _ := apd.NewFromString(strconv.FormatFloat(dim.Min, 'f', -1, 64))
			lb := new(apd.Decimal)
			ctx.Add(lb, lb, min)
			lower_bound, _ := lb.Float64()
			ub := new(apd.Decimal)
			ctx.Add(ub, incr, interval_l)
			ctx.Mul(ub, ub, range_)
			ctx.Add(ub, ub, min)
			upper_bound, _ := ub.Float64()
			interval.range_ = []float64{lower_bound, upper_bound}
			if axis == 0 {
				dim_x = append(dim_x, interval)
			} else {
				dim_j = append(dim_j, interval)
			}
			ctx.Add(incr, incr, interval_l)
		}
		axis += 1
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

	//for _, elem := range x_old{
	//	unit := grid.intersect(elem)
	//	elem.unit_id = unit.id
	//	unit.points = append(unit.points, elem)
	//
	//}
	//
	//for _, elem := range x_update{
	//	unit := grid.intersect(elem)
	//	unit.points = append(unit.points, elem)
	//}
	//
	//for _, elem := range x_new{
	//	unit := grid.intersect(elem)
	//	unit.points = append(unit.points, elem)
	//}
}
func main() {
	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	ch := make(chan PacketData)
	acc := make(chan PacketAcc)
	quit := make(chan int)
	go WindowTimeSlide(ch, acc, quit)
	go UpdateFS(acc)

	if(err != nil){
		log.Fatal(err)
	}

	for {
		data, ci, err := handleRead.ReadPacketData()
		if err != nil && err != io.EOF {
			quit <- 0
			log.Fatal(err)
		} else if err == io.EOF {
			quit <- 0
			break
		} else {
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			ch <- PacketData{packet, ci}
		}
	}
	//dim_min_max := []DimMinMax{{
	//	Min: 1.003,
	//	Max: 10.0004,
	//},{
	//	Min: 1.003,
	//	Max: 10.0004,
	//}}
	//x_old := [][]float64{{1.12344123141234123,2.000012341234123234123},{1.12344123141234123,2.000012341234123234123}}
	//dim_min_max[0].Range = dim_min_max[0].Max - dim_min_max[0].Min
	//dim_min_max[1].Range = dim_min_max[1].Max - dim_min_max[1].Min
	//build2Dgrid(x_old, [][]float64{}, [][]float64{}, dim_min_max)

}