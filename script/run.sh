#!/bin/bash
# ELBOW
for i in {30..30}
do
    # Cluster
    ./main.exe clusterData --db-name="201711281400_300ms.db" --log-path="$(cmd //c cd)\lumber.log" --delta-t=300ms \
    --min-dense-points=5 --min-cluster-points=${i} --points-mode="percentage" \
    --knee-find-elbow=true --knee-smoothing-window=1 --num-knee-flat-points=1
    # Get srcIPs
    sed -n -e 's/.*network_scan_syn anomalies.*SrcIP: \[\(.*\)] DstIP:.*/\1/p' lumber.log > tmp.log
    # Get unique srcIPs
    sort -u tmp.log > elbow_${i}.log
    NUM_ANOMALIES_PREDICTED="$(wc -l elbow_${i}.log | awk '{print $1}')"
    # Get anomalies detected by finding intersection
    NUM_ANOMALIES_ACTUAL="$(comm -12 elbow_${i}.log anomalies_20171128.txt | wc -l)"
    # Insert into database
    sqlite3 roc.db  "insert into confusion_matrix_data (type,min_cluster_points_percentage,anomalies_predicted,anomalies_actual) values ('elbow',${i},${NUM_ANOMALIES_ACTUAL},${NUM_ANOMALIES_PREDICTED});"
    rm lumber.log
    rm tmp.log
done

# KNEE
for i in {30..30}
do
    # Cluster
    ./main.exe clusterData --db-name="201711281400_300ms.db" --log-path="$(cmd //c cd)\lumber.log" --delta-t=300ms \
    --min-dense-points=5 --min-cluster-points=${i} --points-mode="percentage" \
    --knee-find-elbow=false --knee-smoothing-window=1 --num-knee-flat-points=1
    # Get srcIPs
    sed -n -e 's/.*network_scan_syn anomalies.*SrcIP: \[\(.*\)] DstIP:.*/\1/p' lumber.log > tmp.log
    # Get unique srcIPs
    sort -u tmp.log > knee_${i}.log
    NUM_ANOMALIES_PREDICTED="$(wc -l knee_${i}.log | awk '{print $1}')"
    # Get anomalies detected by finding intersection
    NUM_ANOMALIES_ACTUAL="$(comm -12 knee_${i}.log anomalies_20171128.txt | wc -l)"
    # Insert into database
    sqlite3 roc.db  "insert into confusion_matrix_data (type,min_cluster_points_percentage,anomalies_predicted,anomalies_actual) values ('knee',${i},${NUM_ANOMALIES_ACTUAL},${NUM_ANOMALIES_PREDICTED});"
    rm lumber.log
    rm tmp.log
done