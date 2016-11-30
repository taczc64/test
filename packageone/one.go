package one

import(
  "github.com/cihub/seelog"
)

type A struct{
  _a int
  _B string
  C  int
  d  string
}

type B struct {
  say func(word string)
}

func One(){
    seelog.Error("this is package one function write to log file")
}
