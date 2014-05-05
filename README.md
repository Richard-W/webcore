Webcore
=======

A lightweigt and extensible CMS-Framework written in Go. The codebase is still under heavy development
and the API ist likely to change in the future (i am definitely going to change the RegisterFragment-API)

Installation
------------

At this time installation of this is king of tough, because you have to set up
the SQL tables manually. At some point i will add a tool to simplify the task.

For a simple setup execute the following SQL in your database:

```
CREATE TABLE `nodes` (
        `uuid` TEXT NOT NULL PRIMARY KEY,
        `parentId` TEXT,
        `name` TEXT,
        `displayName` TEXT,
        `fragment` TEXT,
        `fragmentOptions` TEXT,
	`deep` INTEGER
);

CREATE TABLE `fragment_markdown` (
        `uuid` TEXT NOT NULL PRIMARY KEY,
        `content` TEXT
);

INSERT INTO nodes (`uuid`, `parentId`, `name`, `displayName`, `fragment`, `fragmentOptions`, `deep`) VALUES
        (
                'e2c9ee58-b3e9-4762-a46f-3dc69905bc5f',
                '',
                'root',
                'Home',
                'markdown',
                'fragment_markdown,uuid,content',
		'0'
        ), (
                '002ac299-3d5d-4460-99d4-58db6c3179e1',
                'e2c9ee58-b3e9-4762-a46f-3dc69905bc5f',
                'secondary',
                'Secondary',
                'markdown',
                'fragment_markdown,uuid,content',
		'0'
        );

INSERT INTO `fragment_markdown` (`uuid`, `content`) VALUES
('e2c9ee58-b3e9-4762-a46f-3dc69905bc5f', '
Home
====

This is some sample site written in markdown
'),

('002ac299-3d5d-4460-99d4-58db6c3179e1', '
Secondary
=========

Some secondary page
');
```

This is a simple database setup with a top page (Home) and a second level page. Now you need a template to actually display something:

template.html:
```
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

Now you can use the following code to display the site:

```
package main

import (
        webcore "github.com/richard-w/webcore"
        sql     "database/sql"
        _       "github.com/mattn/go-sqlite3"
)

func main () {
        webcore.SetAddress ("localhost:8080")
        webcore.SetBaseTemplate ("./template.html")
        db, err := sql.Open ("sqlite3", "./sqlite.db")
        if err != nil {
                panic (err.Error ())
        }
        webcore.SetDatabase (db)
        err = webcore.Run ()
        if err != nil {
                print (err.Error () + "\n")
        }
}
```

In the example i used sqlite3 but you can use every other working sql-driver. See the documentation for the package [database/sql](http://golang.org/pkg/database/sql/)
