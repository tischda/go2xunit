@echo off
go test -check.vv > gocheck.out
type gocheck.out
sleep 10
