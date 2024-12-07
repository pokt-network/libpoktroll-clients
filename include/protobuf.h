#include <stddef.h>
#include <stdint.h>

#ifndef PROTOBUF_H
#define PROTOBUF_H

// Opaque structure to represent a serialized protobuf message
typedef struct {
    uint8_t* type_url;
    size_t type_url_length;
    uint8_t* data;
    size_t data_length;
} serialized_proto;

// Structure to hold metadata about the message array.
typedef struct {
    serialized_proto* messages;  // Pointer to array of message pointers.
    size_t num_messages;         // Number of messages in the array.
} proto_message_array;

#endif //PROTOBUF_H
