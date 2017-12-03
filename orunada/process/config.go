package process

const SEQUENTIAL = 1
const PARALLEL = 2

type Config struct{
	Min_dense_points int
	Min_cluster_points int
	Execution_type int
	Num_cpu int
}