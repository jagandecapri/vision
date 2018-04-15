package anomalies

import (
	"github.com/jagandecapri/vision/tree"
	"sync"
	"github.com/jagandecapri/kneedle"
	"log"
	"sort"
	"github.com/jagandecapri/vision/utils"
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

type DissimilarityMapContainer struct{
	sync.RWMutex
	internal map[int][]DissimilarityVectorContainer
}

type DissimilarityMapPackage struct{
	Key int
	Dis_vector []DissimilarityVectorContainer
}

func NewDissimilarityMapContainer() *DissimilarityMapContainer {
	return &DissimilarityMapContainer{
		internal: make(map[int][]DissimilarityVectorContainer),
	}
}

func (dmc *DissimilarityMapContainer) Load(key int) (value []DissimilarityVectorContainer, ok bool) {
	dmc.RLock()
	result, ok := dmc.internal[key]
	dmc.RUnlock()
	return result, ok
}

func (dmc *DissimilarityMapContainer) Len() (value int) {
	dmc.RLock()
	result := len(dmc.internal)
	dmc.RUnlock()
	return result
}

func (dmc *DissimilarityMapContainer) Delete(key int) {
	dmc.Lock()
	delete(dmc.internal, key)
	dmc.Unlock()
}

func (dmc *DissimilarityMapContainer) Store(key int, value DissimilarityVectorContainer) {
	dmc.Lock()
	dmc.internal[key] = append(dmc.internal[key], value)
	dmc.Unlock()
}

func (dmc *DissimilarityMapContainer) IterateDissimilarityMapContainer() chan DissimilarityMapPackage {
	out := make(chan DissimilarityMapPackage, 10)
	go func(){
		dmc.RLock()
		defer dmc.RUnlock()
		for key, dis_vector := range dmc.internal {
			dmc.RUnlock()
			out <- DissimilarityMapPackage{
					Key: key,
					Dis_vector: dis_vector,
			}
			dmc.RLock()
		}
		close(out)
	}()
	return out
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

func EvidenceAccummulationForOutliers(identifier string, store *DissimilarityMapContainer, num_channels int, done chan struct{}, wg_channels *sync.WaitGroup){
	go func(identifier string, store *DissimilarityMapContainer, num_channels int, done chan struct{}, wg_channels *sync.WaitGroup){
		defer func(){
			log.Println(identifier, "closing iterating channel")
			wg_channels.Done()
		}()

		counter := 0

		for{
			out := store.IterateDissimilarityMapContainer()
			for dmp := range out{
				if len(dmp.Dis_vector) == num_channels{
					log.Println(identifier, "All subspaces processed in disimilarity vector")

					tmp := make(map[int]DissimilarityVector)

					for _, dissimilarity_vector_container := range dmp.Dis_vector{
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
						knee_points, _ := kneedle.Run(knee_data, 1, 1,true) //finding elbows in data

						log.Printf("%v batch: %v data sort: %v knee: %v\n", identifier, counter, knee_data, knee_points)
						if len(knee_points) > 0{
							knee_idx := int(knee_points[0][0])

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
					store.Delete(dmp.Key)
					counter++
				}
			}

			select{
			case <-done:
				return
			default:
			}
		}
	}(identifier, store, num_channels, done, wg_channels)
}