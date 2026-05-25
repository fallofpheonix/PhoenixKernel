; Minimal BIOS boot sector placeholder.
; Build target and kernel loading policy are defined in docs/01-os-from-scratch.md.

[org 0x7c00]

start:
    mov si, message

.print:
    lodsb
    cmp al, 0
    je .halt
    mov ah, 0x0e
    int 0x10
    jmp .print

.halt:
    cli
    hlt
    jmp .halt

message db "MyOS boot sector loaded", 0

times 510 - ($ - $$) db 0
dw 0xaa55

