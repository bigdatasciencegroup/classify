package classify

import (
	"encoding/json"
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// Classifier is a document classifier.
type Classifier struct {
	fcCombos      map[string]map[string]int
	catCounts     map[string]int
	assumedProb   float64
	assumedWeight float64
	cutoffs       map[string]float64
}

// NewClassifier is the contstructor for Classifier
func NewClassifier(assumedProb, assumedWeight float64) *Classifier {
	return &Classifier{
		fcCombos:      make(map[string]map[string]int),
		catCounts:     make(map[string]int),
		assumedProb:   assumedProb,
		assumedWeight: assumedWeight,
		cutoffs:       make(map[string]float64),
	}
}

func (c *Classifier) incFeatureCategory(feature string, cat string) {
	if c.fcCombos[feature] == nil {
		c.fcCombos[feature] = make(map[string]int)
	}
	c.fcCombos[feature][cat]++
}

func (c *Classifier) incCatCount(cat string) {
	c.catCounts[cat]++
}

func (c *Classifier) countFeatureInCategory(feature string, cat string) int {
	if c.fcCombos[feature] == nil {
		c.fcCombos[feature] = make(map[string]int)
	}
	return c.fcCombos[feature][cat]
}

func (c *Classifier) countItemsInCategory(cat string) int {
	return c.catCounts[cat]
}

func (c *Classifier) countTotal() int {
	sum := 0
	for _, count := range c.catCounts {
		sum += count
	}
	return sum
}

func (c *Classifier) categories() []string {
	var out []string
	for cat := range c.catCounts {
		out = append(out, cat)
	}
	return out
}

// Featurer is an interface for derving features from
// an object.
type Featurer interface {
	Features() []string
}

// Train trains the classifier with a document and known
// category.
func (c *Classifier) Train(item Featurer, cat string) {
	features := item.Features()
	for _, feature := range features {
		c.incFeatureCategory(feature, cat)
	}
	c.incCatCount(cat)
}

func (c *Classifier) prob(feature string, cat string) float64 {
	if c.fcCombos[feature] == nil {
		c.fcCombos[feature] = make(map[string]int)
	}
	return float64(c.fcCombos[feature][cat]) / float64(c.catCounts[cat])
}

func (c *Classifier) cprob(feature string, cat string) float64 {
	clf := c.prob(feature, cat)
	if clf == 0.0 {
		return 0.0
	}
	freqSum := 0.0
	for _, currCat := range c.categories() {
		freqSum += c.prob(feature, currCat)
	}
	return clf / freqSum
}

// WeightedProb returns the weighted probability that a feature belongs to a certain category.
func (c *Classifier) WeightedProb(feature string, cat string) float64 {
	basicProb := c.cprob(feature, cat)
	featureSum := float64(0)
	for _, currCat := range c.categories() {
		featureSum += float64(c.countFeatureInCategory(feature, currCat))
	}
	return ((c.assumedWeight * c.assumedProb) + (featureSum * basicProb)) / (c.assumedWeight + featureSum)
}

// FisherProb computes the Fisher probability that an item belongs in a category.
func (c *Classifier) FisherProb(item Featurer, cat string) float64 {
	p := 1.0
	features := item.Features()
	for _, feature := range features {
		wp := c.WeightedProb(feature, cat)
		p *= 1 - wp
	}
	fScore := -2.0 * math.Log(p)
	return distuv.ChiSquared{K: 2.0 * float64(len(features))}.CDF(fScore)
}

func (c *Classifier) String() string {
	pretty, _ := json.MarshalIndent(&c.fcCombos, "", "   ")
	return string(pretty)
}

// SetCutoff sets the minimum probability required to be placed
// in a category.
func (c *Classifier) SetCutoff(cat string, cutoff float64) {
	c.cutoffs[cat] = cutoff
}

// Classify classifies an item.
func (c *Classifier) Classify(item Featurer) string {
	best := ""
	max := 0.0
	for _, cat := range c.categories() {
		p := c.FisherProb(item, cat)
		if p > c.cutoffs[cat] && p > max {
			best = cat
			max = p
		}
	}
	return best
}
