package anomalies

import (
	"github.com/jagandecapri/vision/tree"
	"github.com/jagandecapri/vision/process"
	"fmt"
)

var aggsrc_anomalies = map[string]AnomaliesInterface{}

var aggsrc_dis_vector = map[[2]string][] chan process.DissimilarityVector{}

var aggsrc_subspaces = map[[2]string] chan ProcessPackage{}

var aggdst_anomalies = map[string]AnomaliesInterface{
	"ddos": NewDDOS(),
}

var aggdst_dis_vector = map[[2]string][] chan process.DissimilarityVector{
	[2]string{"nbSrcs", "avgPktSize"}: {aggdst_anomalies["ddos"].GetChannel([2]string{"nbSrcs", "avgPktSize"})},
	[2]string{"perICMP", "perSYN"}: {aggdst_anomalies["ddos"].GetChannel([2]string{"perICMP", "perSYN"})},
	[2]string{"nbSrcPort", "perICMP"}: {aggdst_anomalies["ddos"].GetChannel([2]string{"nbSrcPort", "perICMP"})},
}

var aggdst_subspaces = map[[2]string] chan ProcessPackage{}

var aggsrcdst_anomalies = map[string]AnomaliesInterface{}

var aggsrcdst_dis_vector = map[[2]string][] chan process.DissimilarityVector{}

var aggsrcdst_subspaces = map[[2]string] chan ProcessPackage{}

type SubspaceChannels map[[2]string] chan ProcessPackage

type SubspaceChannelsContainer struct{
	AggSrc SubspaceChannels
	AggDst SubspaceChannels
	AggSrcDst SubspaceChannels
}

var subspace_channels = SubspaceChannelsContainer{
		AggSrc: aggsrc_subspaces,
		AggDst: aggdst_subspaces,
		AggSrcDst: aggsrcdst_subspaces,
}

type ProcessPackage struct{
	X_old        []tree.Point
	X_new_update []tree.Point
}

func Cluster(subspace tree.Subspace, config process.Config, done chan struct{}, outs ...chan process.DissimilarityVector) chan ProcessPackage{
	in := make(chan ProcessPackage)
	counter := 1
	go func() {
		LOOP:
			for {
				select {
				case processPackage := <-in:
					x_old := processPackage.X_old
					x_new_update := processPackage.X_new_update
					subspace.ComputeSubspace(x_old, x_new_update)
					subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
					dissimilarity_map := process.ComputeDissmilarityVector(subspace)
					if len(subspace.GetOutliers()) > 0 {
						fmt.Println("key:", subspace.Subspace_key, "outliers:", subspace.GetOutliers(), "clusters:", subspace.GetClusters())
					}
					for _, out := range outs {
						out <- process.DissimilarityVector{Id: counter, Vector: dissimilarity_map}
					}
					counter++
				case <-done:
					break LOOP
				default:
				}
			}
	}()
	return in
}

func BuildSubspace(subspace_key [2]string) tree.Subspace{
	min_interval := 0.0
	max_interval := 1.0
	interval_length := 0.1
	dim := 2
	scale_factor := 5

	Int_tree := tree.NewIntervalTree(uint64(dim))
	grid := tree.NewGrid()
	ranges := tree.RangeBuilder(min_interval, max_interval, interval_length)
	intervals := tree.IntervalBuilder(ranges, scale_factor)
	units := tree.UnitsBuilder(ranges, dim)

	subspace := tree.Subspace{Grid: &grid, Subspace_key: subspace_key, Scale_factor: scale_factor}
	subspace.SetIntervalTree(&Int_tree)
	for _, interval := range intervals{
		Int_tree.Add(interval)
	}

	for _, unit := range units{
		grid.AddUnit(&unit)
	}
	grid.SetupGrid(interval_length)

	return subspace
}

func ClusteringBuilder(config process.Config, done chan struct{}) SubspaceChannelsContainer {
	for subspace_key, channels := range aggsrc_dis_vector{
		subspace := BuildSubspace(subspace_key)
		in := Cluster(subspace, config, done, channels...)
		aggsrc_subspaces[subspace_key] = in
	}

	for subspace_key, channels := range aggdst_dis_vector{
		subspace := BuildSubspace(subspace_key)
		in := Cluster(subspace, config, done, channels...)
		aggdst_subspaces[subspace_key] = in
	}

	for subspace_key, channels := range aggsrcdst_dis_vector{
		subspace := BuildSubspace(subspace_key)
		in := Cluster(subspace, config, done, channels...)
		aggsrcdst_subspaces[subspace_key] = in
	}

	for _, anomaly := range aggsrc_anomalies{
		anomaly.WaitOnChannels(done)
	}

	for _, anomaly := range aggdst_anomalies{
		anomaly.WaitOnChannels(done)
	}

	for _, anomaly := range aggsrcdst_anomalies{
		anomaly.WaitOnChannels(done)
	}

	return subspace_channels
}