package spam

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

// Comment describes a YouTube comment
type Comment struct {
	CommentID string
	Author    string
	Content   string
	IsSpam    bool
}

// Features generates features of the comment for document classification
func (c Comment) Features() []string {
	var out []string
	out = append(out, "[author]"+c.Author)
	n := 2
	words := strings.Split(removeSpecialChars(strings.ToLower(c.Content)), " ")
	for i := 0; i < len(words)-n; i++ {
		feature := "[content]"
		for j := 0; j < n; j++ {
			feature += words[i+j] + " "
		}
		out = append(out, strings.TrimSpace(feature))
	}
	return out
}

func removeSpecialChars(s string) string {
	alpha := "abcdefghijklmnopqrstuvwxyz "
	var out string
	for _, char := range strings.ToLower(s) {
		if strings.ContainsRune(alpha, char) {
			out += string(char)
			continue
		}
	}
	return out
}

// LoadComments loads comments from a CSV
func LoadComments(src string) ([]Comment, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	_, err = r.Read() // Headers
	if err != nil {
		return nil, err
	}
	var out []Comment
	for {
		fields, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		c := Comment{
			CommentID: fields[0],
			Author:    fields[1],
			Content:   fields[3],
			IsSpam:    fields[4] == "1",
		}
		out = append(out, c)
	}
	return out, nil
}
