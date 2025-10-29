package iec61850

/*
#include <iec61850_client.h>
#include <linked_list.h>
#include <stdbool.h>
#include "client_file.h"
*/
import "C"
import (
	"io"
	"time"
	"unsafe"
)

//export downloadHandler
func downloadHandler(parameter unsafe.Pointer, buffer *C.uint8_t, bytesRead C.uint32_t) C.bool {
	/*
	 * `parameter` is an unsafe pointer passed by the C code when invoking this callback.
	 * This `parameter` is actually the memory address that we allocated in the `GetFile` method
	 * (using `C.malloc`) and contains a reference to a Go object (`downloadWriter`).
	 *
	 * Step 1: Cast `parameter` from `unsafe.Pointer` to `*unsafe.Pointer`. This is necessary
	 * because `parameter` is a pointer to a memory location holding another pointer, which is
	 * the real reference to the `downloadWriter` object.
	 *
	 * Step 2: Dereference the first pointer (`*(*unsafe.Pointer)(parameter)`) to extract the
	 * actual pointer stored within the allocated memory and cast it to `*downloadWriter`.
	 *
	 * This two-step process safely retrieves the `downloadWriter` object we originally passed
	 * from Go to C.
	 */
	writer := (*downloadWriter)(*(*unsafe.Pointer)(parameter))

	if bytesRead == 0 {
		return C.bool(true)
	}

	// If bytes were received, convert the C buffer `buffer` to a Go byte slice.

	/*
	 * `buffer` is a `*C.uint8_t` (a pointer to a C byte array), and `bytesRead` is the
	 * number of bytes available in this buffer. To work with this data in Go, we need to
	 * convert it into a Go byte slice.
	 *
	 * Step 1: Cast `buffer` (a C pointer to the first byte) to `unsafe.Pointer` for general use.
	 * Step 2: Use `C.GoBytes` to convert the pointer and the length (`C.int(bytesRead)`)
	 * to a Go byte slice. This creates a Go-managed copy of the underlying C data.
	 */
	data := C.GoBytes(unsafe.Pointer(buffer), C.int(bytesRead))

	// Write the data to the `downloadWriter`. Handle any errors during the writing process.
	n, err := writer.Write(data)
	if err != nil {
		// If an error occurs, register it using the `Err` method of `downloadWriter`
		// and return `false` to indicate the failure to the C side.
		writer.Err(err)
		return C.bool(false)
	}

	// Verify if the number of bytes written matches the number of bytes read.
	if n != int(bytesRead) {
		// If there's a mismatch, report it as a short write error and return `false`.
		writer.Err(io.ErrShortWrite)
		return C.bool(false)
	}

	// If everything is successful, return `true` to inform the C code.
	return C.bool(true)
}

type FileDirectoryEntry struct {
	Name         string
	Size         int
	LastModified time.Time
}

// GetFileDirectory retrieves a list of file directory entries from the server for the specified directory.
// The directory parameter specifies the path of the directory to list; an empty string retrieves the root directory.
// Returns a slice of FileDirectoryEntry containing file name, size, and last modified time.
// Returns an error if the retrieval fails or encounters client/server communication errors.
func (c *Client) GetFileDirectory(directory string) ([]FileDirectoryEntry, error) {
	var clientError C.IedClientError
	var dirName *C.char = nil

	if directory != "" {
		dirName = C.CString(directory)
		defer C.free(unsafe.Pointer(dirName))
	}

	var root C.LinkedList
	root = C.IedConnection_getFileDirectory(c.conn, &clientError, dirName)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}

	entries := make([]FileDirectoryEntry, 0)

	var directoryEntry C.LinkedList
	directoryEntry = C.LinkedList_getNext(root)
	for directoryEntry != nil {
		var entry C.FileDirectoryEntry
		entry = (C.FileDirectoryEntry)(directoryEntry.data)

		cFileName := C.FileDirectoryEntry_getFileName(entry)
		cFileSize := C.FileDirectoryEntry_getFileSize(entry)
		// UTC timestamp in milliseconds
		cLastModifiedTimeStamp := C.FileDirectoryEntry_getLastModified(entry)

		// Convert UTC timestamp in milliseconds to time.Time
		lastModified := time.UnixMilli(int64(cLastModifiedTimeStamp))

		fileName := C.GoString(cFileName)
		fileSize := int(cFileSize)
		entries = append(entries, FileDirectoryEntry{fileName, fileSize, lastModified})

		directoryEntry = C.LinkedList_getNext(directoryEntry)
	}

	C.LinkedList_destroyDeep(root, (C.LinkedListValueDeleteFunction)(C.FileDirectoryEntry_destroy))
	return entries, nil
}

type downloadWriter struct {
	writer io.Writer
	err    error
}

func (w *downloadWriter) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *downloadWriter) Err(err error) {
	w.err = err
}

// GetFile downloads a file from the server and writes its content to the provided io.Writer.
func (c *Client) GetFile(w io.Writer, filename string) error {
	var clientError C.IedClientError
	localFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(localFilename))

	// Wrap the writer in a downloadWriter.
	dw := &downloadWriter{writer: w}

	// Allocate memory to hold the downloadWriter and pass the C pointer instead of Go pointer.
	cDw := C.malloc(C.size_t(unsafe.Sizeof(dw)))
	defer C.free(cDw)
	*(*unsafe.Pointer)(cDw) = unsafe.Pointer(dw)

	// Pass the Go callback and the pinned context (cDw) to the C function.
	C.IedConnection_getFile(
		c.conn,
		&clientError,
		localFilename,
		C.IedClientGetFileHandler(C.downloadHandlerCGo),
		cDw,
	)

	if dw.err != nil {
		return dw.err
	}

	return GetIedClientError(clientError)
}
