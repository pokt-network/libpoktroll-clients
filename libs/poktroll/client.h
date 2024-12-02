#ifndef CLIENT_H
#define CLIENT_H

#include <stdint.h>

//typedef struct {
//  BlockID block_id;
//  Block *Block;
//} BlockResult;

typedef void *(callback_fn)(void *data, char **err);

typedef int64_t go_ref;

#endif