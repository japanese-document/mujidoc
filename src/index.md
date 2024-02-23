{}
---
# Mujidoc

Mujidoc is a simple static site generator.

## Installation

```bash
go install github.com/japanese-document/mujidoc/cmd/mujidoc@0.0.8
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

```
{ "category": "Go",  "order": 0, "date": "2024-01-03 15:00" }
---
# Title 

something
```

#### category

This is the name of the category to which the page belongs.
This is one of `CATEGORIES`.
If `SINGLE_PAGE` is `true`, `category` is unnecessary.

#### order

This specifies the position at which the page is displayed within the category.
If `SINGLE_PAGE` is `true`, `order` is unnecessary.

#### date

This is the value for `pubDate` in the RSS feed.
If `RSS` is `false`, `date` is unnecessary.

### Configuration file

You need to place a configuration file named `.env.mujidoc` in working directory. Here is an example:

```
CATEGORIES=Go,Python,Ubuntu
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

#### CATEGORIES

This specifies categories separated by commas.
The categories will be displayed in the order specified.

#### BASE_URL

This is the base URL of the generated site.

#### PAGE_LAYOUT

This specifies the layout for the generated HTML.

#### INDEX_PAGE_HEADER

This is the value of the h1 element for `index.html`.

#### INDEX_PAGE_TITLE

This is the title of the h1 element for `index.html`.

#### INDEX_PAGE_DESCRIPTION

This is the description of the h1 element for `index.html`.

#### INDEX_PAGE_LAYOUT

This specifies the layout for `index.html`.

#### OUTPUT_DIR

This is the directory where the generated HTML is output.

#### SOURCE_DIR

This is the directory where markdown files are placed.

#### SINGLE_PAGE

If you want a single page, specify `true` for this option.

#### RSS

If you want to generate an RSS feed, specify `true` for this option.
The generated RSS feed file name is `rss.xml` in `OUTPUT_DIR`.

#### TIME_ZONE

This specifies the timezone to be used for the RSS feed.

### Images

When providing image files, you need to place the image files in the `SOURCE_DIR/images` directory.
The image files move to `OUTPUT_DIR/images`.

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