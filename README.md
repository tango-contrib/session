session [![Build Status](https://drone.io/github.com/tango-contrib/session/status.png)](https://drone.io/github.com/tango-contrib/session/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/session)](http://gocover.io/github.com/tango-contrib/session)
======

Session is a session middleware for [Tango](https://github.com/lunny/tango). 

## Installation

    go get github.com/tango-contrib/session

## Simple Example

```Go
package main

import (
    "github.com/lunny/tango"
    "github.com/tango-contrib/session"
)

type SessionAction struct {
    session.Session
}

func (a *SessionAction) Get() string {
    a.Session.Set("test", "1")
    return a.Session.Get("test").(string)
}

func main() {
    o := tango.Classic()
    o.Use(session.New(session.Options{
        MaxAge:time.Minute * 20,
        }))
    o.Get("/", new(SessionAction))
}
```

## Getting Help

- [API Reference](https://gowalker.org/github.com/tango-contrib/session)
