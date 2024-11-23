#ifndef STRUCTS_H
#define STRUCTS_H

#include <stdint.h>

//typedef struct {
//  char *clientCtx;
//  char *query;
//  int64_t blockHeight;
//  int64_t msgType;
//} XXX;

enum ErrorCode {
  EVENTS_BYTES_SYNC_ERROR,
  EVENTS_BYTES_ASYNC_ERROR,
};

#endif // STRUCTS_H