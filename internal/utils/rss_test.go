package utils

import (
	"reflect"
	"testing"
)

func TestSortedLimitedArray_Push(t *testing.T) {
	compareFn := func(a, b int) bool {
		return a < b
	}
	type args struct {
		items []int
	}
	tests := []struct {
		name  string
		array *SortedLimitedArray[int]
		args  args
		want  []int
	}{
		{
			name:  "Add single item",
			array: NewSortedLimitedArray([]int{}, 5, compareFn),
			args:  args{items: []int{3}},
			want:  []int{3},
		},
		{
			name:  "Add multiple items",
			array: NewSortedLimitedArray([]int{9, 9, 6, 9, 9}, 5, compareFn),
			args:  args{items: []int{3, 1, 4, 7}},
			want:  []int{1, 3, 4, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.array.Push(tt.args.items...)
			if got := tt.array.Items; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortedLimitedArray.Push() = %v, want %v", got, tt.want)
			}
		})
	}
}
