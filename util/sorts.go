package util

import "strconv"

type SortSlice struct {
	FieldValues []interface{}
	ComparseFn  func(a, b interface{}) bool
}

func (ss SortSlice) Len() int { return len(ss.FieldValues) }

func (ss SortSlice) Swap(i, j int) {
	ss.FieldValues[i], ss.FieldValues[j] = ss.FieldValues[j], ss.FieldValues[i]
}

func (ss SortSlice) Less(i, j int) bool {
	if ss.ComparseFn == nil {
		panic("ss.fn is nil")
	}
	return ss.ComparseFn(ss.FieldValues[i], ss.FieldValues[j])
}

func (ss *SortSlice) Register(fn func(a, b interface{}) bool) {
	ss.ComparseFn = fn
}

func StringInt64BigSmallComparseFn(a, b interface{}) bool {
	aStr, bStr := a.(string), b.(string)
	aValue, err := strconv.ParseInt(aStr, 10, 64)
	if err != nil {
		panic(err)
	}
	bValue, err := strconv.ParseInt(bStr, 10, 64)
	if err != nil {
		panic(err)
	}

	return aValue > bValue
}

func StringInt64SmallBigComparseFn(a, b interface{}) bool {
	aStr, bStr := a.(string), b.(string)
	aValue, err := strconv.ParseInt(aStr, 10, 64)
	if err != nil {
		panic(err)
	}
	bValue, err := strconv.ParseInt(bStr, 10, 64)
	if err != nil {
		panic(err)
	}

	return aValue < bValue
}

func Uint64SmallBigComparseFn(a, b interface{}) bool {
	return a.(uint64) < b.(uint64)
}
