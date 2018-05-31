package classify

import (
	"testing"

	"github.com/tomjcleveland/classify/spam"
)

func TestClassifier(t *testing.T) {
	c := NewClassifier(0.5, 1.0)

	// Train
	sources := []string{
		"Youtube01-Psy.csv",
		"Youtube02-KatyPerry.csv",
		"Youtube03-LMFAO.csv",
		"Youtube04-Eminem.csv",
	}
	for _, src := range sources {
		trainingSet, err := spam.LoadComments("./fixtures/" + src)
		if err != nil {
			t.Fatal(err)
		}
		for _, comment := range trainingSet {
			cat := "good"
			if comment.IsSpam {
				cat = "bad"
			}
			c.Train(comment, cat)
		}
	}

	// Test
	testSet, err := spam.LoadComments("./fixtures/Youtube02-KatyPerry.csv")
	if err != nil {
		t.Fatal(err)
	}
	failures := 0
	for _, comment := range testSet {
		cat := c.Classify(comment)
		correct := "good"
		if comment.IsSpam {
			correct = "bad"
		}
		if cat != correct {
			failures++
		}
	}
	t.Logf("Classified %d/%d correctly", len(testSet)-failures, len(testSet))
}
