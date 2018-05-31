# Classify
A simple library for classifying documents using (20th-century) machine learning.

## Quickstart
Any document you want to classify must implement the `Featurer` interface:
```go
type Featurer interface {
    Features() []string
}
```
Features are used by the classifier to classify documents.

Here's an example using YouTube spam data:
```go
package main

import (
    "github.com/tomjcleveland/classify/spam"
)

func main() {
    // Create classifier
    c := NewClassifier(0.5, 1.0)

    // Train using YouTube spam training set
    // https://archive.ics.uci.edu/ml/datasets/YouTube+Spam+Collection
    sources := []string{
        "Youtube01-Psy.csv",
        "Youtube02-KatyPerry.csv",
        "Youtube03-LMFAO.csv",
        "Youtube04-Eminem.csv",
    }
    for _, src := range sources {
        trainingSet, err := spam.LoadComments("./fixtures/" + src)
        if err != nil {
            panic(err)
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
        panic(err)
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
    fmt.Printf("Classified %d/%d correctly\n", len(testSet)-failures, len(testSet))
}
```
Running the above would give you this output:
```
Classified 347/350 correctly
```
Ship it.
