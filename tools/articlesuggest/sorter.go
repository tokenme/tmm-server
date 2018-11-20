package articlesuggest

type ScoreSorter []Item

type Item struct {
	Key uint64
	Val float64
}

func NewScoreSorter(m map[uint64]float64) ScoreSorter {

	ms := make(ScoreSorter, 0, len(m))

	for k, v := range m {

		ms = append(ms, Item{k, v})

	}

	return ms

}

func (ms ScoreSorter) Len() int {

	return len(ms)

}

func (ms ScoreSorter) Less(i, j int) bool {

	return ms[i].Val < ms[j].Val // 按值排序

	//return ms[i].Key < ms[j].Key // 按键排序

}

func (ms ScoreSorter) Swap(i, j int) {

	ms[i], ms[j] = ms[j], ms[i]

}
