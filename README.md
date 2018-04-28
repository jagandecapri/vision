# Installation

* Install Go following [here](https://golang.org/doc/install)
* Install WinPCAP following [here](https://www.winpcap.org/)

# Accessing app server in browser

1. Navigate to localhost:3001/main.html

# Running from command-line

## Prepare Data

`go run main.go prepareData --db-name="201705021400.db" --pcap-file-path="C:\Users\Jack\Downloads\201705021400.pcap" --delta-t=300ms`

## Cluster Data

`go run main.go clusterData --db-name="2017012281400.db" --log-path="C:\Users\Jack\go\src\github.com\jagandecapri\vision\logs\lumber.log" --num-cpu=0 --min-dense-points=10 --min-cluster-points=10 --delta-t=300ms`