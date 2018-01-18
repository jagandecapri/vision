package process

import "math"

/**
 * Given set of values look for the elbow/knee points.
 * See paper: "Finding a Kneedle in a Haystack: Detecting Knee Points in System Behavior"
 * @author Luke Bermingham
 */

 type Kneedle struct{}



// Finds the indices of all local minimum or local maximum values.
func (k Kneedle) findCandidateIndices(data []float64, findMinima bool) []int{
	candidates := []int{}
	//an index is considered a candidate if both of its adjacent indices are
	//greater or less (depending on whether we want local minima or local maxima)
	for i := 1; i < len(data) - 1; i++ {
		prev := data[i-1]
		cur := data[i]
		next := data[i+1];
		var isCandidate bool
		if findMinima {
			isCandidate = (prev > cur && next > cur)
		} else {
			isCandidate = (prev < cur && next < cur)
		}
		if(isCandidate){
			candidates = append(candidates, i)
		}
	}
	return candidates
}


// Find the index in the data the represents a most exaggerated elbow point.

func findElbowIndex(data []float64) int{
	bestIdx := 0;
	bestScore := 0.0;
	for i := 0; i < len(data); i++ {
		score := math.Abs(data[i]);
		if(score > bestScore){
			bestScore = score;
			bestIdx = i;
		}
	}
	return bestIdx;
}

// In this step we prepare the data by normalising into unit range 0-1
// and also subtracting the value from its normalised index value.
func prepare(data []float64) []float64{
	//find min and max value
	curMin := math.Inf(1);
	curMax := math.Inf(-1);
	for _, v := range data {
		if v < curMin {
			curMin = v;
		}
		if v > curMax {
			curMax = v;
		}
}

//normalise the data using min-max normalisation
//and also subtract each value from its normalised index
range_ := curMax - curMin;
normalisedData := []float64{};

for i := 0; i < len(normalisedData); i++ {
	normalisedIndex := float64(i) / float64(len(data));
	normalisedData[i] = ((data[i] - curMin) / range) - normalisedIndex;
}
return normalisedData;
}



public double run(double[] data){
if(data.length <= 1){
return 0;
}
double[] normalisedData = prepare(gaussianSmooth(data, 3));
int elbowIdx = findElbowIndex(normalisedData);
return data[elbowIdx];
}

/**
 * This algorithm finds the so-called elbow/knee in the data.
 * It does this by sorting the data, then making a line between the start
 * and end data points in the sorted data. Each point in the data is the projected
 * onto this line, and the point with the biggest euclidean distance is considered
 * the most likely elbow.
 * See paper: "Finding a Kneedle in a Haystack: Detecting Knee Points in System Behavior"
 * for more details.
 * @param data The data to find an elbow in.
 * @param s How many "flat" points to require before we consider it a knee/elbow.
 * @param findElbows Whether to find elbows or knees.
 * @return The elbow or knee values.
 */
public ArrayList<Double> run(double[] data, double s, boolean findElbows){
ArrayList<Double> localMinMaxPts = new ArrayList<>();
//smooth the data to make local minimum/maximum easier to find (this is Step 1 in the paper)
double[] smoothedData = gaussianSmooth(data, 3);
//prepare the data into the unit range and subtract normalised index (this is step 2 & 3 in the paper)
double[] normalisedData = prepare(smoothedData);
//find candidate indices (this is step 4 in the paper)
{
ArrayList<Integer> candidateIndices = findCandidateIndices(normalisedData, findElbows);
//go through each candidate index, i, and see if the indices after i are satisfy the threshold requirement
//(this is step 5 in the paper)
double step = 1.0/data.length;
step = findElbows ? step * s : step * -s;

//check each candidate to see if it is a real elbow/knee
for (int i = 0; i < candidateIndices.size(); i++) {
Integer candidateIdx = candidateIndices.get(i);
Integer endIdx = (i + 1 < candidateIndices.size()) ? candidateIndices.get(i+1) : data.length;

double threshold = normalisedData[candidateIdx] + step;

for (int j = candidateIdx + 1; j < endIdx; j++) {
boolean isRealElbowOrKnee = (findElbows) ?
normalisedData[j] > threshold : normalisedData[j] < threshold;
if(isRealElbowOrKnee) {
localMinMaxPts.add(data[candidateIdx]);
break;
}
}
}
}
return localMinMaxPts;
}

}

