#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libclients.h"

static void test_events_query_client_constructor(void **state) {
  (void) state; // how to use?

  const int64_t clientId = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
  assert_int_not_equal(clientId, -1);
  assert_int_not_equal(clientId, 0);

  char *err = malloc(1024);
  int64_t eventsBytesObsID = EventsQueryClientEventsBytes(clientId, "", &err);
  printf("err: %s\n", err);
  assert_int_not_equal(eventsBytesObsID, -1);
  assert_int_equal(0, strcmp("" ,err));
}

// Main test runner
int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_events_query_client_constructor),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}
