# mujidoc

## Usage

### page layout

```html
<!DOCTYPE html>
<html lang="ja">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <meta name="twitter:card" content="summary" />
    <meta property="og:url" content="__URL__" />
    <meta property="og:title" content="__TITLE__" />
    <meta property="og:description" content="__DESCRIPTION__" />
    <meta property="og:image" content="https://example.com" />
    <meta name="theme-color" content="#f1f7fe" />
    <meta name="description" content="__DESCRIPTION__" />
    <link rel="alternate" type="application/rss+xml" title="RSS" href="rss.xml" />
    <link rel="icon" type="image/png" href="/images/favicon.png" />
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
      <a href="/tips">Top</a>
    </footer>
  </body>
</html>
```

### .env.mujidoc

```
CSS_PATH=/app.css?v=001
BASE_URL=http://127.0.0.1:8000/
PAGE_LAYOUT=src/layout.html
INDEX_PAGE_HEADER=Lit
INDEX_PAGE_TITLE=Lit
INDEX_PAGE_DESCRIPTION=JavaScriptのWeb UI フレームワーク
INDEX_PAGE_LAYOUT=src/layout.html
OUTPUT_DIR=docs
SOURCE_DIR=src
SINGLE_PAGE=false
RSS=true
```