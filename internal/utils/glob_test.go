package utils

import (
	"io/fs"
	"reflect"
	"testing"

	"github.com/japanese-document/mujidoc/internal/utils/mock_utils"
	"go.uber.org/mock/gomock"
)

type MockFilePath struct {
	WalkDirFunc func(root string, fn fs.WalkDirFunc) error
}

func (m MockFilePath) WalkDir(root string, fn fs.WalkDirFunc) error {
	return m.WalkDirFunc(root, fn)
}

type mockDirEntry struct {
	isDir bool
	name  string
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func TestGetMarkDownFileNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	mfp := mock_utils.NewMockIFilePath(ctrl)
	mfp.EXPECT().WalkDir(gomock.Any(), gomock.Any()).Do(func(root string, fn fs.WalkDirFunc) error {
		err := fn("/path/to/markdown1.md", mockDirEntry{isDir: false}, nil)
		if err != nil {
			return err
		}
		err = fn("/path/to/not_markdown.txt", mockDirEntry{isDir: false}, nil)
		if err != nil {
			return err
		}
		err = fn("/path/to/markdown2.md", mockDirEntry{isDir: false}, nil)
		if err != nil {
			return err
		}
		err = fn("/path/to/markdown3.md", mockDirEntry{isDir: true}, nil)
		if err != nil {
			return err
		}
		return nil
	}).Return(nil)

	type args struct {
		fp   IFilePath
		root string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Find markdown files successfully",
			args: args{
				fp:   mfp,
				root: "/path/to",
			},
			want:    []string{"/path/to/markdown1.md", "/path/to/markdown2.md"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMarkDownFileNames(tt.args.fp, tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMarkDownFileNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMarkDownFileNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
