#!/bin/bash
TOTAL_UNIQ_SRC_FLOWS=4435961

# ELBOW
for i in {5..100..5}
do
    echo "ELBOW: ${i}"
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
    TRUE_POSITIVES="$(comm -12 elbow_${i}.log anomalies_20171128.txt | wc -l)"
    FALSE_POSITIVES="$(comm -23 elbow_${i}.log anomalies_20171128.txt | wc -l)"
    FALSE_NEGATIVES="$(comm -13 elbow_${i}.log anomalies_20171128.txt | wc -l)"
    TRUE_NEGATIVES="$(expr ${TOTAL_UNIQ_SRC_FLOWS} - ${NUM_ANOMALIES_PREDICTED} - ${FALSE_NEGATIVES})"
    # Insert into database
    sqlite3 roc.db  "insert into confusion_matrix_data (type,min_cluster_points_percentage,true_positives,false_positives,true_negatives,false_negatives) \
    values ('elbow',${i},${TRUE_POSITIVES},${FALSE_POSITIVES},${TRUE_NEGATIVES},${FALSE_NEGATIVES});"
    mv lumber.log lumber_elbow.log
    rm tmp.log
done

# KNEE
for i in {5..100..5}

do
    echo "KNEE: ${i}"
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
    TRUE_POSITIVES="$(comm -12 knee_${i}.log anomalies_20171128.txt | wc -l)"
    FALSE_POSITIVES="$(comm -23 knee_${i}.log anomalies_20171128.txt | wc -l)"
    FALSE_NEGATIVES="$(comm -13 knee_${i}.log anomalies_20171128.txt | wc -l)"
    TRUE_NEGATIVES="$(expr ${TOTAL_UNIQ_SRC_FLOWS} - ${NUM_ANOMALIES_PREDICTED} - ${FALSE_NEGATIVES})"
    # Insert into database
    sqlite3 roc.db  "insert into confusion_matrix_data (type,min_cluster_points_percentage,true_positives,false_positives,true_negatives,false_negatives) \
    values ('knee',${i},${TRUE_POSITIVES},${FALSE_POSITIVES},${TRUE_NEGATIVES},${FALSE_NEGATIVES});"
    mv lumber.log lumber_knee.log
    rm tmp.log
done