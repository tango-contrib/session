pongo2 [![Build Status](https://drone.io/github.com/tango-contrib/session/status.png)](https://drone.io/github.com/tango-contrib/session/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/session)](http://gocover.io/github.com/tango-contrib/session)
======

** HEAVILLY DEVELOPMENT **
Session is a session middleware for [Tango](https://github.com/lunny/tango). 

## Installation

    go get github.com/tango-contrib/session

## Simple Example

```Go
package main

import (
    "github.com/lunny/tango"
    "gopkg.in/flosch/pongo2.v3"
    "github.com/tango-contrib/tpongo2"
)

type RenderAction struct {
    tpango2.Render
}

func (a *RenderAction) Get() error {
    return a.RenderString("Hello {{ name }}!", pongo2.Context{
        "name": "tango",
    })
}

func main() {
    o := tango.Classic()
    o.Use(tpango2.Default())
    o.Get("/", new(RenderAction))
}
```

## Getting Help

- [API Reference](https://gowalker.org/github.com/tango-contrib/tpongo2)
