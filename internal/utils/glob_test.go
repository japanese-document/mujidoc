package utils

import (
	"io/fs"
	"reflect"
	"testing"

	"github.com/japanese-document/mujidoc/internal/utils/mock_utils"
	"go.uber.org/mock/gomock"
)

// #region Mock
type mockDirEntry struct {
	isDir bool
	name  string
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func newMockFilePath(t *testing.T) (*mock_utils.MockIFilePath, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	return mock_utils.NewMockIFilePath(ctrl), ctrl
}

func mfp1(mfp *mock_utils.MockIFilePath) *mock_utils.MockIFilePath {
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
	})
	return mfp
}

func mfp2(mfp *mock_utils.MockIFilePath) *mock_utils.MockIFilePath {
	mfp.EXPECT().WalkDir(gomock.Any(), gomock.Any()).Do(func(root string, fn fs.WalkDirFunc) error {
		err := fn("/path/to/mark down1.md", mockDirEntry{isDir: false}, nil)
		if err != nil {
			return err
		}
		return nil
	})
	return mfp
}

func TestGetMarkDownFileNames(t *testing.T) {
	type args struct {
		fp   func(mfp *mock_utils.MockIFilePath) *mock_utils.MockIFilePath
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
				fp:   mfp1,
				root: "/path/to",
			},
			want:    []string{"/path/to/markdown1.md", "/path/to/markdown2.md"},
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				fp:   mfp2,
				root: "/path/to",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mfp, ctrl := newMockFilePath(t)
			defer ctrl.Finish()
			got, err := GetMarkDownFileNames(tt.args.fp(mfp), tt.args.root)
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
