# About

Github repository for my Masters project. Papers written on my work is at [here](http://dx.doi.org/10.14299/ijser.2019.01.03) and [here](http://dx.doi.org/10.26782/jmcms.2019.12.00078). This a system written using Go concurrency paradigm to parallelize network anomaly detection.

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

`sed -n -e 's/.*network_scan_syn anomalies.*SrcIP: \[\(.*\)] DstIP:.*/\1/p' lumber.log > output_network_scan_syn_srcIP.log`

## Benchmarking IGDCA

`cd tree`
`go test -v -bench=BenchmarkIGDCA -run=^a` => `-run=^a` is to avoid running other tests

Alternatively, to build the test binary and run the benchmarks:

`cd tree`
`go test -c`
`./test.test.exe -test.v -test.bench=BenchmarkIGDCA -test.run=^a` => `-run=^a` is to avoid running other tests

## Extracting Packets from PCAP File Using TCPDump

`tcpdump -r 201711281400.pcap '(src 101.153.157.180 or 58.54.221.13 or 203.189.147.172 or 175.4.177.51 or 223.156.186.71 or 118.129.17.181 or 175.9.65.254 or 222.24.31.148 or 118.133.85.205 or 222.66.165.112 or 115.202.212.156 or 58.237.50.137 or 181.229.23.185 or 189.43.123.254)' -F pcap -w out_ntscsyn_all_21_04_2018.pcap`
