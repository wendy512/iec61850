#include "client_file.h"

// A C wrapper for the Go function 'downloadHandler'. This is necessary because Go does not allow direct usage of Go function pointers as callbacks in C.
bool downloadHandlerCGo(void* parameter, uint8_t* buffer, uint32_t bytesRead) {
    return downloadHandler(parameter, buffer, bytesRead);
}