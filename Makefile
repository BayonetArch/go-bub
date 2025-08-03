CC=go


all : main

main : main.go
	go build ./main.go && ./$@  



