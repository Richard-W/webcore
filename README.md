Webcore
=======

A lightweight and extensible CMS-Framework written in Go. The codebase is still under heavy development
and the API ist likely to change in the future (i am definitely going to change the RegisterFragment-API)

Installation
------------

You need some kind of template to display your contents:

template.html:
```html
<!DOCTYPE html>
<html>
<head>
        <title>Webcore example</title>
        <meta charset="utf-8" />
</head>
<body>

<header>
        <nav>
                <ul>
                        {{range .MainMenu}}
                                <li>
                                        <a href="{{.Path}}">{{.Name}}</a>
                                        {{if .SubMenu}}
                                                <ul>
                                                {{range .SubMenu}}
                                                        <li><a href="{{.Path}}">{{.Name}}</a></li>
                                                {{end}}
                                                </ul>
                                        {{end}}
                                </li>
                        {{end}}
                </ul>
        </nav>
</header>

<main>
{{.Body}}
</main>

</html>
```

The MainMenu contains the pages up to the third node level. You may choose to not display the third level. The top page and the secondary pages
are treated equally.

webcore uses the [text/template](http://golang.org/pkg/text/template/) package from the go standard library. It is passed a TemplateInterface-struct.
The definition of the structs used can be obtained using `godoc github.com/richard-w/webcore`

Now you can use the following code to display the site:

```go
package main

import (
        webcore "github.com/richard-w/webcore"
        sql     "database/sql"
        _       "github.com/mattn/go-sqlite3"
)

func main () {
        webcore.SetAddress ("localhost:8080")
        webcore.SetBaseTemplate ("./template.html")
        realDB, err := sql.Open ("sqlite3", "./sqlite.db")
        if err != nil {
                panic (err.Error ())
        }
	database, err := webcore.NewDatabase (realDB)
        if err != nil {
                panic (err.Error ())
        }
        webcore.SetDatabase (database)
        err = webcore.Run ()
        if err != nil {
                print (err.Error () + "\n")
        }
}
```

In the example i used sqlite3 but you can use every other working sql-driver. See the documentation for the package [database/sql](http://golang.org/pkg/database/sql/). Webcore automatically installs a core database and maintains it.
