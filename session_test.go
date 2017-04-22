// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/lunny/tango"
)

type SessionAction struct {
	Session
}

var _ Sessioner = &SessionAction{}

func (action *SessionAction) Get() string {
	action.Session.Set("test", "1")
	return action.Session.Get("test").(string)
}

func TestSession(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	tg.Use(New())
	tg.Get("/", new(SessionAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type PrefixGenerator struct {
	generator IdGenerator
	prefix    string
}

func NewPrefixGenerator(generator IdGenerator, prefix string) IdGenerator {
	return &PrefixGenerator{
		generator: generator,
		prefix:    prefix,
	}
}

func (p *PrefixGenerator) Gen(req *http.Request) Id {
	return Id(p.prefix + string(p.generator.Gen(req)))
}

func (p *PrefixGenerator) IsValid(id Id) bool {
	return p.generator.IsValid(Id(strings.TrimPrefix(string(id), p.prefix)))
}

func TestSession2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	tg.Use(New(Options{
		Generator: NewPrefixGenerator(NewSha1Generator(string(GenRandKey(16))), "prefix_"),
	}))
	tg.Get("/", new(SessionAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
