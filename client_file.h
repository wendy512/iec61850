#include <iec61850_client.h>
#include <stdbool.h>

// Declare the Go function
bool downloadHandler(void* parameter, uint8_t* buffer, uint32_t bytesRead);
// Declare the wrapper function
bool downloadHandlerCGo(void* parameter, uint8_t* buffer, uint32_t bytesRead);
