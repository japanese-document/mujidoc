package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/japanese-document/mujidoc/internal/css"
	"github.com/japanese-document/mujidoc/internal/utils"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func init() {
	err := godotenv.Load(".env.mujidoc")
	if err != nil {
		log.Fatal(err)
	}
}

func createPageHtmlFileTask(markDownFileName, indexMenu, pageLayout, sourceDir, outputDir string) func() error {
	return func() error {
		content, err := os.ReadFile(markDownFileName)
		if err != nil {
			return errors.WithStack(err)
		}
		_, md, err := utils.GetMetaAndMd(string(content))
		if err != nil {
			return err
		}
		title := utils.CreateTitle(md)
		dir, name := utils.GetDirAndName(markDownFileName)
		url := utils.CreateURL(dir, name)
		headerList, err := utils.CreateHeaderList(md)
		if err != nil {
			return err
		}
		page, err := utils.CreatePage(pageLayout, md, title, url, indexMenu, headerList)
		if err != nil {
			return err
		}
		dirPath := utils.CreateHTMLFilePath(dir, sourceDir, outputDir)
		if !utils.IsDirExists(dirPath) {
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		htmlFileName := filepath.Join(dirPath, fmt.Sprintf("%s.html", name))
		err = os.WriteFile(htmlFileName, []byte(page), 0644)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func createCopyImageDirTask() func() error {
	return func() error {
		srcImageDir, distImageDir := utils.CreateSrcAndOutputDir()
		if utils.IsDirExists(srcImageDir) {
			err := os.MkdirAll(distImageDir, os.ModePerm)
			if err != nil {
				return errors.WithStack(err)
			}
			return utils.CopyDir(srcImageDir, distImageDir)
		}
		return nil
	}
}

func createIndexHtmlFileTask(layout string, outputDir string, indexItems []utils.IndexItem) func() error {
	return func() error {
		indexPageLayout, err := os.ReadFile(layout)
		if err != nil {
			return errors.WithStack(err)
		}
		indexPage, err := utils.CreateIndexPage(string(indexPageLayout), indexItems)
		if err != nil {
			return err
		}
		htmlFileName := filepath.Join(outputDir, "index.html")
		return os.WriteFile(htmlFileName, []byte(indexPage), 0644)
	}
}

func main() {
	err := os.Mkdir(os.Getenv("OUTPUT_DIR"), os.ModePerm)
	if err != nil {
		log.Fatalf("%s is existed. Please remove %s.", os.Getenv("OUTPUT_DIR"), os.Getenv("OUTPUT_DIR"))
	}
	markDownFileNames, err := utils.GetMarkDownFileNames(os.Getenv("SOURCE_DIR"), ".md")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	pages := []*utils.Page{}
	if os.Getenv("SINGLE_PAGE") != "true" {
		pages, err = utils.CreatePages(markDownFileNames)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
	indexItems, err := utils.CreateIndexItems(pages)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	indexMenu := utils.CreateIndexMenu(indexItems)
	_pageLayout, err := os.ReadFile(os.Getenv("PAGE_LAYOUT"))
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	pageLayout := string(_pageLayout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	// markdownからhtmlを生成する
	for _, markDownFileName := range markDownFileNames {
		task := createPageHtmlFileTask(markDownFileName, indexMenu, pageLayout, os.Getenv("SOURCE_DIR"), os.Getenv("OUTPUT_DIR"))
		eg.Go(task)
	}

	// 画像をコピーする
	task := createCopyImageDirTask()
	eg.Go(task)

	// index.htmlを作成する
	if os.Getenv("SINGLE_PAGE") != "true" {
		task = createIndexHtmlFileTask(os.Getenv("INDEX_PAGE_LAYOUT"), os.Getenv("OUTPUT_DIR"), indexItems)
		eg.Go(task)
	}

	// CSSファイルを作成する
	task = css.CreateWriteTask(os.Getenv("OUTPUT_DIR"), utils.CSS_FILE_NAME)
	eg.Go(task)

	if os.Getenv("RSS") == "true" {
		task := utils.CreateRssFileTask(pages)
		eg.Go(task)
	}

	if err := eg.Wait(); err != nil {
		log.Fatalf("%+v", err)
	}
}
