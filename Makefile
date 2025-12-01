data:
	python3 ./create_measurements.py 1000000000

.PHONY: go
go:
	cd go && go build -o 1brc-go main.go
	cd go && fish -c "time ./1brc-go"
