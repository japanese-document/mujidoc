package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestFilter(t *testing.T) {
	type args[T any] struct {
		items      []T
		callbackFn func(item T, index int) (bool, error)
	}
	tests := []struct {
		name    string
		args    args[int] // ここでは int 型をテストケースの型として使用
		want    []int
		wantErr bool
	}{
		{
			name: "Filter even numbers",
			args: args[int]{
				items: []int{1, 2, 3, 4, 5},
				callbackFn: func(item int, index int) (bool, error) {
					return item%2 == 0, nil // 偶数をフィルタリング
				},
			},
			want:    []int{2, 4},
			wantErr: false,
		},
		{
			name: "Error case",
			args: args[int]{
				items: []int{1, 2, 3},
				callbackFn: func(item int, index int) (bool, error) {
					if item == 2 {
						return false, errors.New("error on 2")
					}
					return true, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.args.items, tt.args.callbackFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type args[T any, R any] struct {
		items      []T
		callbackFn func(item T, index int) (R, error)
	}
	tests := []struct {
		name    string
		args    args[int, string] // ここでは int 型の入力と string 型の出力をテスト
		want    []string
		wantErr bool
	}{
		{
			name: "Convert numbers to strings",
			args: args[int, string]{
				items: []int{1, 2, 3},
				callbackFn: func(item int, index int) (string, error) {
					return fmt.Sprintf("Number: %d", item), nil
				},
			},
			want:    []string{"Number: 1", "Number: 2", "Number: 3"},
			wantErr: false,
		},
		{
			name: "Error case",
			args: args[int, string]{
				items: []int{1, 2, 3},
				callbackFn: func(item int, index int) (string, error) {
					if item == 2 {
						return "", errors.New("error on 2")
					}
					return fmt.Sprintf("Number: %d", item), nil
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Map(tt.args.items, tt.args.callbackFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapToSlice(t *testing.T) {
	tests := []struct {
		name  string
		input map[int]string
		want  []string
	}{
		{
			name: "Convert map to sorted slice",
			input: map[int]string{
				2: "b",
				1: "a",
				3: "c",
			},
			want: []string{"a", "b", "c"}, // キーの昇順に対応する値
		},
		{
			name:  "Empty map",
			input: map[int]string{},
			want:  []string{}, // 空のスライス
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapToSlice(tt.input)
			if err != nil {
				t.Errorf("MapToSlice() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
