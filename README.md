# Mujidoc

Mujidoc is a simple static site generator.

## Installation

```bash
go install github.com/japanese-document/mujidoc/cmd/mujidoc@0.0.6
```

## Usage

You create a configuration file and execute the command.

### Execute

following command:

```
mujidoc
```

Please note that this command first deletes the directory specified in `OUTPUT_DIR`, then creates a new directory at `OUTPUT_DIR`.

### Content

You place markdown files with the following metadata in `SOURCE_DIR`.
If `RSS` is `false`, `date` in the following example is unnecessary.

```
{ "category": {"name": "Go", "order": 6},  "order": 7, "date": "2024-01-03 15:00" }
---
# Title 

something
```

#### category.name

#### category.order

#### order

#### date

### Configuration file

You need to place a configuration file named `.env.mujidoc` in working directory. Here is an example:

```
BASE_URL=https://japanese-document.github.io/mujidoc
PAGE_LAYOUT=src/layout.html
INDEX_PAGE_HEADER=Mujidoc
INDEX_PAGE_TITLE=Mujidoc
INDEX_PAGE_DESCRIPTION="Mujidoc is simple html page generator."
INDEX_PAGE_LAYOUT=src/layout.html
OUTPUT_DIR=docs
SOURCE_DIR=src
SINGLE_PAGE=false
RSS=true
TIME_ZONE="Asia/Tokyo"
```

### Images

When providing image files, you need to place the image files in the `SOURCE_DIR/images` directory. The image files move to `OUTPUT_DIR`. `SOURCE_DIR` is configured in the configuration file.

### Page layout

Page layout files (`PAGE_LAYOUT` and `INDEX_PAGE_LAYOUT`) must be placed as below. Their names are configured by `PAGE_LAYOUT` and `INDEX_PAGE_LAYOUT` in the configuration file. 

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <meta name="twitter:card" content="summary" />
    <meta property="og:url" content="__URL__" />
    <meta property="og:title" content="__TITLE__" />
    <meta property="og:description" content="__DESCRIPTION__" />
    <meta property="og:image" content="https://japanese-document.github.io/mujidoc/images/favicon.png" />
    <meta name="theme-color" content="#f1f7fe" />
    <meta name="description" content="__DESCRIPTION__" />
    <link rel="icon" type="image/png" href="https://japanese-document.github.io/mujidoc/images/favicon.png" />
    <title>__TITLE__</title>
    <link rel="stylesheet" href="__CSS__" type="text/css"  media="all" />
  </head>
  <body class="container">
    <div class="left-side">__INDEX__</div>
    <main class="main markdown-body">
      __BODY__
    </main>
    <div class="right-side">__HEADER__</div>
    <footer class="footer markdown-body">
      <a href="/mujidoc">Top</a>
    </footer>
  </body>
</html>
```