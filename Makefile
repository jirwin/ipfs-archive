
all: clean
	go install github.com/jirwin/ipfs-archive

clean:
	rm -rf snapshots

.PHONY: all clean