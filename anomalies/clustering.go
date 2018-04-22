package anomalies

import (
	"github.com/jagandecapri/vision/tree"
	"github.com/jagandecapri/vision/utils"
	"log"
	"sync"
)

var aggsrc_anomalies = map[string]AnomaliesInterface{
	"network_scan_syn": NewNetworkScanSYN(),
}

var aggsrc_dis_vector = map[[2]string][] chan DissimilarityVectorContainer{
	[2]string{"perSYN", "nbDstPort"}: {aggsrc_anomalies["network_scan_syn"].GetChannel([2]string{"perSYN", "nbDstPort"})},
	[2]string{"nbDstPort", "nbDsts"}: {aggsrc_anomalies["network_scan_syn"].GetChannel([2]string{"nbDstPort", "nbDsts"})},
	[2]string{"nbDstPort", "avgPktSize"}: {aggsrc_anomalies["network_scan_syn"].GetChannel([2]string{"nbDstPort", "avgPktSize"})},
}

var aggsrc_subspaces = map[[2]string] chan ProcessPackage{}

var aggdst_anomalies = map[string]AnomaliesInterface{
	"ddos": NewDDOS(),
}

var aggdst_dis_vector = map[[2]string][] chan DissimilarityVectorContainer{
	[2]string{"nbSrcs", "avgPktSize"}: {aggdst_anomalies["ddos"].GetChannel([2]string{"nbSrcs", "avgPktSize"})},
	[2]string{"perICMP", "perSYN"}: {aggdst_anomalies["ddos"].GetChannel([2]string{"perICMP", "perSYN"})},
	[2]string{"nbSrcPort", "perICMP"}: {aggdst_anomalies["ddos"].GetChannel([2]string{"nbSrcPort", "perICMP"})},
}

var aggdst_subspaces = map[[2]string] chan ProcessPackage{}

var aggsrcdst_anomalies = map[string]AnomaliesInterface{}

var aggsrcdst_dis_vector = map[[2]string][] chan DissimilarityVectorContainer{}

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

func Cluster(subspace tree.Subspace, config utils.Config, outs ...chan DissimilarityVectorContainer) chan ProcessPackage{
	in := make(chan ProcessPackage)
	counter := 1

	go func() {
			defer func(){
				log.Println("Closing anomalies channels")
				for _, out := range outs{
					close(out)
				}
			}()

			for {
				select {
				case processPackage, open := <-in:
					if open{
						x_old := processPackage.X_old
						x_new_update := processPackage.X_new_update
						log.Println("Start Cluster: ", subspace.Subspace_key)
						subspace.ComputeSubspace(x_old, x_new_update)
						subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
						dissimilarity_vectors := ComputeDissmilarityVector(subspace)
						if len(subspace.GetOutliers()) > 0 {
							log.Printf("counter: %v key: %v outliers: %v clusters: %+v", counter, subspace.Subspace_key, len(subspace.GetOutliers()), subspace.GetClusters())
							for _, cluster := range subspace.GetClusters(){
								validate_cluster := tree.ValidateCluster(cluster)
								log.Printf(" cluster %v - validate cluster: %v", cluster.Cluster_id, validate_cluster)
							}
							log.Printf("\n")
						} else {
							log.Println("counter: ", counter, " key:", subspace.Subspace_key, " No outliers Cluster: ")
							clusters := subspace.GetClusters()
							for _, cluster := range clusters{
								log.Printf("cluster ID: %+v", cluster.Cluster_id)
								for rg, unit := range cluster.ListOfUnits{
									log.Printf(" unit ID: %+v rg: %+v", unit.Id, rg)
								}
								log.Printf("\n")
							}
						}
						for _, out := range outs {
							out <- DissimilarityVectorContainer{Id: counter, DissimilarityVectors: dissimilarity_vectors}
						}
						counter++
					} else{
						return
					}
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
		grid.AddUnit(unit)
	}
	grid.SetupGrid(interval_length)

	return subspace
}

func ClusteringBuilder(config utils.Config, done chan struct{}) SubspaceChannelsContainer {

	for subspace_key, channels := range aggsrc_dis_vector{
		subspace := BuildSubspace(subspace_key)
		in := Cluster(subspace, config, channels...)
		aggsrc_subspaces[subspace_key] = in
	}

	for subspace_key, channels := range aggdst_dis_vector{
		subspace := BuildSubspace(subspace_key)
		in := Cluster(subspace, config, channels...)
		aggdst_subspaces[subspace_key] = in
	}

	for subspace_key, channels := range aggsrcdst_dis_vector{
		subspace := BuildSubspace(subspace_key)
		in := Cluster(subspace, config, channels...)
		aggsrcdst_subspaces[subspace_key] = in
	}

	var wg sync.WaitGroup
	wg.Add(len(aggsrc_anomalies) + len(aggdst_anomalies) + len(aggsrcdst_anomalies))

	for _, anomaly := range aggsrc_anomalies{
		anomaly.WaitOnChannels(&wg)
	}

	for _, anomaly := range aggdst_anomalies{
		anomaly.WaitOnChannels(&wg)
	}

	for _, anomaly := range aggsrcdst_anomalies{
		anomaly.WaitOnChannels(&wg)
	}

	go func(){
		wg.Wait()
		log.Println("Signal anomalies done. Closing done channel")
		close(done)
	}()

	return subspace_channels
}