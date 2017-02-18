all: input.proto

input.proto: irgen/input.def irgen/irgen
	cd irgen && ./irgen input.def ..
	# Sanity check
	go build
	gofmt -w .

.PHONY: irgen/irgen
irgen/irgen:
	cd irgen && make && go build
