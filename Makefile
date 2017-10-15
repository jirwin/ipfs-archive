CMDS=$(shell find ./cmd/* -maxdepth 1 -type d -exec basename {} \; )

cmd_targets = $(addprefix ./cmd/, $(CMDS))

all: 
	go install -v $(cmd_targets)

frontend:
	cd frontend && npm run build && statik -src=build && npm run build:clean

$(CMDS):
	go install -v $(addprefix ./cmd/, $@)

.PHONY: all frontend
