package anomalies

import (
	"github.com/jagandecapri/vision/tree"
	"sync"
)

type DissimilarityVector struct{
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

func ComputeDissmilarityVector(subspace tree.Subspace) []DissimilarityVector{
	dissimilarity_vectors := []DissimilarityVector{}

	outliers := subspace.GetOutliers()
	cluster := subspace.GetBiggestCluster()
	center_biggest_cluster := cluster.GetCenter()

	for _, outlier := range outliers{
		distance := center_biggest_cluster.Distance(&outlier)
		dissimilarity_vector := DissimilarityVector{PointKey: outlier.Key, Distance: distance}
		dissimilarity_vectors = append(dissimilarity_vectors, dissimilarity_vector)
	}

	return dissimilarity_vectors
}