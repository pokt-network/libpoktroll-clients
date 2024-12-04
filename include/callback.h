#ifndef CALLBACK_H
#define CALLBACK_H

#include <context.h>

// Declare the bridge functions that Go will call
static void bridge_success(AsyncOperation *op, void *results);
static void bridge_error(AsyncOperation *op, char *err);

#endif // CALLBACK_H