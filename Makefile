CLANG ?= clang
CFLAGS ?= -g -O2 -target bpf -D__TARGET_ARCH_x86 -I/usr/include/x86_64-linux-gnu

SRC_DIR := src
OUT_DIR := src
VMLINUX_H := $(SRC_DIR)/vmlinux.h

.PHONY: all clean

all: $(OUT_DIR)/phoenix_exec.o

$(VMLINUX_H):
	@if [ ! -s $(VMLINUX_H) ]; then \
		echo "vmlinux.h is missing or empty. Attempting to generate from /sys/kernel/btf/vmlinux..."; \
		if [ -f /sys/kernel/btf/vmlinux ]; then \
			bpftool btf dump file /sys/kernel/btf/vmlinux format c > $(VMLINUX_H); \
		else \
			echo "Error: /sys/kernel/btf/vmlinux not found and $(VMLINUX_H) is not provided."; \
			exit 1; \
		fi; \
	else \
		echo "Using existing $(VMLINUX_H)."; \
	fi

$(OUT_DIR)/phoenix_exec.o: $(SRC_DIR)/phoenix_exec.c $(VMLINUX_H)
	$(CLANG) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(OUT_DIR)/*.o
