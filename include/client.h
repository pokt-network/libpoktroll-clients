#ifndef CLIENT_H
#define CLIENT_H

#include <stddef.h>
#include <stdint.h>

//typedef struct {
//  BlockID block_id;
//  Block *Block;
//} BlockResult;

// typedef void (callback_fn)(void *ctx);

typedef int64_t go_ref;

// typedef struct {
//     char **err;
//     callback_fn *on_success;
//     callback_fn *on_error;
// } tx_client_context;

// TODO_IN_THIS_COMMIT: move the below to protobuf.h and update CMakeFiles.txt.

// Opaque structure to represent a serialized protobuf message
typedef struct {
    char* type_url;
    uint8_t* data;
    size_t length;
} serialized_proto;

// Structure to hold metadata about the message array.
typedef struct {
    serialized_proto* messages;  // Pointer to array of message pointers.
    size_t num_messages;         // Number of messages in the array.
} proto_message_array;

#endif