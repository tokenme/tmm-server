package tfidf

import (
	"bytes"
	"encoding/csv"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

type Weighter struct {
	words map[string]float64
}

func NewWeighter() *Weighter {
	return &Weighter{
		words: make(map[string]float64),
	}
}

func (this *Weighter) AddDoc(words []string) {
	founder := make(map[string]struct{}, len(words))
	for _, w := range words {
		if _, found := founder[w]; found {
			continue
		}
		this.words[w] += 1
	}
}

func (this *Weighter) Score(terms []string) ScoreSorter {
	tf := make(map[string]float64)
	tt := float64(len(terms))
	docs := float64(len(this.words))

	for _, term := range terms {
		tf[term] += 1
	}

	// tft tf(t) term frequency of t
	tfidf := make(map[string]float64, len(tf))

	for term, freq := range tf {
		var dwt float64
		if count, found := this.words[term]; found {
			dwt = count
		}
		tft := freq / tt
		var idf float64
		if 0 == dwt {
			idf = 0
		} else {
			idf = math.Log10(docs / dwt)
		}
		tfidf[term] = tft * idf
	}
	sorter := NewScoreSorter(tfidf)
	sort.Sort(sort.Reverse(sorter))
	return sorter
}

func (this *Weighter) Save(filename string) error {
	fi, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer fi.Close()
	buf := new(bytes.Buffer)
	r2 := csv.NewWriter(buf)
	for word, score := range this.words {
		r2.Write([]string{word, strconv.FormatInt(int64(score), 10)})
		r2.Flush()
	}
	fi.WriteString(buf.String())
	return nil
}

func (this *Weighter) Load(filename string) error {
	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fi.Close()
	reader := csv.NewReader(fi)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if len(record) != 2 {
			continue
		}
		this.words[record[0]], _ = strconv.ParseFloat(record[1], 64)
	}
	return nil
}
