#include <stdint.h>

static volatile uint16_t *const VGA_TEXT = (uint16_t *)0xB8000;

void kernel_main(void) {
    const char *message = "MyOS kernel loaded";

    for (uint32_t i = 0; message[i] != '\0'; ++i) {
        VGA_TEXT[i] = (uint16_t)message[i] | 0x0F00;
    }

    for (;;) {
    }
}

