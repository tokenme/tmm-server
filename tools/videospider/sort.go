package videospider

type VideoSorter []VideoLink

func NewVideoSorter(ll []VideoLink) VideoSorter {

	ls := make(VideoSorter, 0, len(ll))

	for _, v := range ll {
		ls = append(ls, v)
	}

	return ls

}

func (ls VideoSorter) Len() int {

	return len(ls)

}

func (ls VideoSorter) Less(i, j int) bool {

	return ls[i].Resolution < ls[j].Resolution // 按值排序

}

func (ls VideoSorter) Swap(i, j int) {

	ls[i], ls[j] = ls[j], ls[i]

}
