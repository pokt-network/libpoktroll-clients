// context.h
#ifndef CONTEXT_H
#define CONTEXT_H

#include <pthread.h>
#include <stdbool.h>
#include <stdlib.h>

// Forward declarations
typedef struct AsyncContext AsyncContext;
typedef struct AsyncOperation AsyncOperation;

// Callback function types
typedef void (*callback_fn)(void* ctx, const void* result);

// Define callback function types
typedef void (*success_callback)(AsyncContext* ctx, const void* result);
typedef void (*error_callback)(AsyncContext* ctx, const char* error);
typedef void (*cleanup_callback)(AsyncContext* ctx);

typedef struct AsyncContext {
    pthread_mutex_t mutex;
    pthread_cond_t cond;
    bool completed;
    bool success;
    void* data;
    size_t data_len;
    int error_code;
    char error_msg[256];
} AsyncContext;

typedef struct AsyncOperation {
    AsyncContext* ctx;
    success_callback on_success;
    error_callback on_error;
    cleanup_callback cleanup;
} AsyncOperation;

void init_context(AsyncContext* ctx);
void cleanup_context(AsyncContext* ctx);
void handle_error(AsyncContext* ctx, const char* error);
void handle_success(AsyncContext* ctx, const void* result);
bool wait_for_completion(AsyncContext* ctx, int timeout_ms);

#endif // CONTEXT_H