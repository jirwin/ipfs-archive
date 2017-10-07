
all: clean
	go install ./cmd/ipfs-archive

clean:
	rm -rf snapshots

.PHONY: all clean