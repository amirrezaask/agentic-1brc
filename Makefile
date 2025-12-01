data:
	python3 ./create_measurements.py 1000000000

.PHONY: go
go-cursor-auto-medium:
	cd go-cursor-auto && go build -o 1brc-go main.go
	cd go && fish -c "time ./1brc-go ../data/medium.txt"

go-cursor-auto-small:
	cd go-cursor-auto && go build -o 1brc-go main.go
	cd go && fish -c "time ./1brc-go ../data/small.txt"

go-cursor-auto-measurements:
	cd go-cursor-auto && go build -o 1brc-go main.go
	cd go && fish -c "time ./1brc-go ../data/measurements.txt"