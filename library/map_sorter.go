package library

import "sort"

/*
 * 对map进行按value进行排序
 * 按value降序排
 */
type MapValueSorter struct {
	Map  map[string]int
	Keys []string
	Vals []int
}

func SortMapByValue(m map[string]int) *MapValueSorter {
	vs := &MapValueSorter{
		Map:  m,
		Keys: make([]string, 0, len(m)),
		Vals: make([]int, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Vals = append(vs.Vals, v)
	}
	vs.Sort()

	return vs
}

func (this *MapValueSorter) Sort() {
	sort.Sort(this)
}
func (this *MapValueSorter) Len() int {
	return len(this.Vals)
}
func (this *MapValueSorter) Less(i, j int) bool {
	return this.Vals[i] > this.Vals[j]
}
func (this *MapValueSorter) Swap(i, j int) {
	this.Vals[i], this.Vals[j] = this.Vals[j], this.Vals[i]
	this.Keys[i], this.Keys[j] = this.Keys[j], this.Keys[i]
}

/*
 * 对map进行按key进行排序
 * 按key降序排
 */
type MapKeySorter struct {
	Map  map[string]int
	Keys []string
	Vals []int
}

func (this *MapKeySorter) Len() int {
	return len(this.Map)
}

func (this *MapKeySorter) Less(i, j int) bool {
	return this.Map[this.Keys[i]] < this.Map[this.Keys[j]]
}

func (this *MapKeySorter) Swap(i, j int) {
	this.Vals[i], this.Vals[j] = this.Vals[j], this.Vals[i]
	this.Keys[i], this.Keys[j] = this.Keys[j], this.Keys[i]
}

func SortMapByKey(m map[string]int) *MapKeySorter {
	ks := new(MapKeySorter)
	ks.Map = m
	ks.Keys = make([]string, 0, len(m))
	ks.Vals = make([]int, 0, len(m))
	for k, v := range m {
		ks.Keys = append(ks.Keys, k)
		ks.Vals = append(ks.Vals, v)
	}
	sort.Sort(ks)
	return ks
}
