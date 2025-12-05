#include <iec61850_client.h>
#include <stdbool.h>
#include <stdint.h>

// Forward declaration of the Go callback
extern bool goFileHandlerCallback(void* parameter, uint8_t* buffer, uint32_t bytesRead);

// C wrapper function to call IedConnection_getFile with Go callback
void callGetFile(IedConnection conn, IedClientError* error, const char* fileName, void* handlerParam) {
    IedConnection_getFile(conn, error, fileName, goFileHandlerCallback, handlerParam);
}
