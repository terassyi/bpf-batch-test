
VMLINUX = bpf/include/vmlinux.h
BIN_DIR = bin
BATCH_BIN = $(BIN_DIR)/batch
NORMAL_BIN = $(BIN_DIR)/normal
LOAD_BIN = $(BIN_DIR)/load
PORG_NAME = xdpprog
BPF_PKG = pkg/bpf
GEN_GO_FILE = $(BPF_PKG)/$(PORG_NAME)_bpfeb.go

.PHONY: vmlinux
$(VMLINUX):
	bpftool btf dump file /sys/kernel/btf/vmlinux format c > $(VMLINUX)

.PHONY: build
build: $(VMLINUX) $(GEN_GO_FILE) $(BATCH_BIN) $(NORMAL_BIN) $(LOAD_BIN)

.PHONY: clean
clean:
	rm $(VMLINUX)
	rm $(BIN_DIR)/*
	rm $(BPF_PKG)/$(PORG_NAME)*

$(BATCH_BIN): $(BIN_DIR) $(GEN_GO_FILE)
	go build -o $(BATCH_BIN) batch/main.go

$(NORMAL_BIN): $(BIN_DIR) $(GEN_GO_FILE)
	go build -o $(NORMAL_BIN) normal/main.go

$(LOAD_BIN): $(BIN_DIR) $(GEN_GO_FILE)
	go build -o $(LOAD_BIN) load/main.go

.PHONY: bpf2go
$(GEN_GO_FILE):
	go generate ./...

$(BIN_DIR):
	mkdir -p bin

