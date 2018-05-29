# Installation

* Install Go following [here](https://golang.org/doc/install)
* Install WinPCAP following [here](https://www.winpcap.org/)

# Accessing app server in browser

1. Navigate to localhost:3001/main.html

# Running from command-line

## Prepare Data

`go run main.go prepareData --db-name="201711281400.db" --pcap-file-path="C:\Users\Jack\Downloads\201705021400.pcap" --delta-t=300ms`

## Cluster Data

`go run main.go clusterData --db-name="201711281400.db" --log-path="C:\Users\Jack\go\src\github.com\jagandecapri\vision\logs\lumber.log" --num-cpu=0 --min-dense-points=10 --min-cluster-points=10 --delta-t=300ms
--window-array-len=2 num-knee-flat-points=1 knee-smoothing-window=1 knee-find-elbow=true`

## Extracting Anomalies from Log

Given that network scan of syn type needs to be extracted

`sed -n -e 's/.*network_scan_syn anomalies.*SrcIP: \[\(.*\)] DstIP:.*/\1/p' lumber.log > output_nextwork_scan_syn_srcIP.log`

## Benchmarking IGDCA

`cd tree`
`go test -v -bench=BenchmarkIGDCA -run=^a` => `-run=^a` is to avoid running other tests