#include "context.h"
#include <errno.h>
#include <pthread.h>
#include <stdbool.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

void init_context(AsyncContext* ctx) {
    memset(ctx, 0, sizeof(AsyncContext));
    pthread_mutex_init(&ctx->mutex, NULL);
    pthread_cond_init(&ctx->cond, NULL);
    ctx->completed = false;
    ctx->success = false;
}

void cleanup_context(AsyncContext* ctx) {
    if (ctx) {
        if (ctx->data) {
            free(ctx->data);
            ctx->data = NULL;
        }
        ctx->data_len = 0;
        ctx->error_code = 0;
        memset(ctx->error_msg, 0, sizeof(ctx->error_msg));

        pthread_mutex_destroy(&ctx->mutex);
        pthread_cond_destroy(&ctx->cond);
    }
}

void handle_error(AsyncContext* ctx, const char* error) {
    if (ctx && error) {
        pthread_mutex_lock(&ctx->mutex);

        ctx->error_code = errno;
        strncpy(ctx->error_msg, error, sizeof(ctx->error_msg) - 1);
        ctx->completed = true;
        ctx->success = false;

        pthread_cond_signal(&ctx->cond);
        pthread_mutex_unlock(&ctx->mutex);

        fprintf(stderr, "Error %d: %s\n", ctx->error_code, ctx->error_msg);
    }
}

void handle_success(AsyncContext* ctx, const void* result) {
    if (ctx && result) {
        pthread_mutex_lock(&ctx->mutex);

        ctx->completed = true;
        ctx->success = true;

        pthread_cond_signal(&ctx->cond);
        pthread_mutex_unlock(&ctx->mutex);

        printf("Operation completed successfully with result\n");
    }
}

void* async_worker(void* arg) {
    AsyncOperation* op = (AsyncOperation*)arg;

    // Simulate some work
    usleep(100000);  // 100ms delay

    if (op->on_success) {
        op->on_success(op->ctx, "Operation result");
    }

    return NULL;
}

void perform_async_operation(AsyncOperation* op) {
    if (!op || !op->ctx) {
        return;
    }

    pthread_t worker_thread;
    pthread_create(&worker_thread, NULL, async_worker, op);
    pthread_detach(worker_thread);
}

bool wait_for_completion(AsyncContext* ctx, int timeout_ms) {
    struct timespec ts;
    clock_gettime(CLOCK_REALTIME, &ts);
    ts.tv_sec += timeout_ms / 1000;
    ts.tv_nsec += (timeout_ms % 1000) * 1000000;
    
    pthread_mutex_lock(&ctx->mutex);
    
    while (!ctx->completed) {
        int rc = pthread_cond_timedwait(&ctx->cond, &ctx->mutex, &ts);
        if (rc == ETIMEDOUT) {
            pthread_mutex_unlock(&ctx->mutex);
            return false;
        }
    }
    
    bool success = ctx->success;
    pthread_mutex_unlock(&ctx->mutex);
    return success;
}