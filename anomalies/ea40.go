package anomalies

import (
	"github.com/jagandecapri/vision/tree"
	"sync"
	"github.com/jagandecapri/kneedle"
	"log"
	"sort"
	"github.com/jagandecapri/vision/utils"
	"github.com/jagandecapri/vision/cmd"
)

type DissimilarityVector struct{
	Id int
	PointKey tree.PointKey
	Distance float64
}

type DissimilarityVectorContainer struct{
	Id int
	DissimilarityVectors []DissimilarityVector
}

type DissimilarityMapPackage struct{
	Key int
	Dis_vector []DissimilarityVectorContainer
}

type DissimilarityMapContainer struct{
	cm sync.Map
	Identifier string
	Num_chan int
	Counter int
}

func NewDissimilarityMapContainer2(identifier string, num_chan int) *DissimilarityMapContainer {
	return &DissimilarityMapContainer{
		cm:sync.Map{},
		Identifier: identifier,
		Num_chan: num_chan,
	}
}

func (dmc2 *DissimilarityMapContainer) Load(key int) (value []DissimilarityVectorContainer) {
	result, ok := dmc2.cm.Load(key)
	tmp := []DissimilarityVectorContainer{}
	log.Println("Is stored result found: ", ok)
	if ok {
		tmp = result.([]DissimilarityVectorContainer)
	}
	return tmp
}

func (dmc2 *DissimilarityMapContainer) Delete(key int) {
	dmc2.cm.Delete(key)
}

func (dmc2 *DissimilarityMapContainer) Store(key int, value DissimilarityVectorContainer) {
	result := dmc2.Load(key)
	log.Println("Len Stored Result:", len(result))
	if len(result) < (dmc2.Num_chan - 1){
		log.Println("Store result")
		result = append(result, value)
		dmc2.cm.Store(key, result)
	} else {
		log.Println("Store result == num_chan")
		result = append(result, value)
		go EvidenceAccumulationForOutliers(result, dmc2.Identifier, dmc2.Counter)
		dmc2.Counter++
		dmc2.cm.Delete(key)
	}
}

func EvidenceAccumulationForOutliers(result []DissimilarityVectorContainer, identifier string, counter int){
	tmp := make(map[int]DissimilarityVector)

	for _, dissimilarity_vector_container := range result{
		for _, dissimilarity_vector := range dissimilarity_vector_container.DissimilarityVectors{
			if val, ok := tmp[dissimilarity_vector.Id]; ok{
				val.Distance += dissimilarity_vector.Distance
				tmp[dissimilarity_vector.Id] = val
			} else {
				tmp[dissimilarity_vector.Id] = dissimilarity_vector
			}
		}
	}

	type kv struct {
		Key   int
		Value float64
	}

	var ss []kv

	for k, v := range tmp {
		ss = append(ss, kv{Key: k, Value: v.Distance})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	var knee_data [][]float64

	for i := 0; i < len(ss); i++ {
		kv := ss[i]
		knee_data = append(knee_data, []float64{float64(i), kv.Value})
	}

	if len(knee_data) > 0{
		knee_points, _ := kneedle.Run(knee_data, cmd.NumKneeFlatPoints, cmd.KneeSmoothingWindow,cmd.KneeFindElbow) //finding elbows in data

		log.Printf("num_knee_flat_points: %v knee_smoothing_window %v knee_find_elbow: %v\n", cmd.NumKneeFlatPoints, cmd.KneeSmoothingWindow, cmd.KneeFindElbow)
		log.Printf("%v batch: %v data sort: %v knee: %v\n", identifier, counter, knee_data, knee_points)
		if len(knee_points) > 0{
			knee_idx := int(knee_points[len(knee_points) - 1][0])

			log.Printf("%v knee_idx searched: %v knee points: %v", identifier, knee_idx, knee_points)
			anomalies := ss[knee_idx + 1:]
			for _, anomaly := range anomalies{
				dis_vec := tmp[anomaly.Key]
				srcIP := utils.UniqString(dis_vec.PointKey.SrcIP)
				dstIP := utils.UniqString(dis_vec.PointKey.DstIP)
				srcPort := utils.UniqString(dis_vec.PointKey.SrcPort)
				dstPort := utils.UniqString(dis_vec.PointKey.DstPort)
				log.Printf("%v anomalies: Batch: %v Keys: SrcIP: %+v DstIP: %+v SrcPort: %+v DstPort: %+v Distance: %+v", identifier, counter, srcIP, dstIP, srcPort, dstPort, dis_vec.Distance)
			}

			non_anomalies := ss[:knee_idx + 1]

			for _, non_anomaly := range non_anomalies{
				dis_vec := tmp[non_anomaly.Key]
				srcIP := utils.UniqString(dis_vec.PointKey.SrcIP)
				dstIP := utils.UniqString(dis_vec.PointKey.DstIP)
				srcPort := utils.UniqString(dis_vec.PointKey.SrcPort)
				dstPort := utils.UniqString(dis_vec.PointKey.DstPort)
				log.Printf("%v non_anomalies: Batch: %v Keys: SrcIP: %+v DstIP: %+v SrcPort: %+v DstPort: %+v Distance: %+v", identifier, counter, srcIP, dstIP, srcPort, dstPort, dis_vec.Distance)
			}
		}
	}
}

func ComputeDissmilarityVector(subspace tree.Subspace) []DissimilarityVector{
	dissimilarity_vectors := []DissimilarityVector{}

	outliers := subspace.GetOutliers()
	cluster := subspace.GetBiggestCluster()
	center_biggest_cluster := cluster.GetCenter()

	for _, outlier := range outliers{
		distance := center_biggest_cluster.Distance(&outlier)
		dissimilarity_vector := DissimilarityVector{Id: outlier.Id, PointKey: outlier.Key, Distance: distance}
		dissimilarity_vectors = append(dissimilarity_vectors, dissimilarity_vector)
	}

	return dissimilarity_vectors
}
