#include <stdlib.h>
#include "libclients.h"
#include "unity.h"

static void test_events_query_client(void) {
    const go_ref clientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));
    const go_ref eventsBytesObsRef = EventsQueryClientEventsBytes(clientRef, "", &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsBytesObsRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsBytesObsRef, 0);

    FreeGoMem(eventsBytesObsRef);
    FreeGoMem(clientRef);
    free(err);
}G

static void test_depinject_supply(void) {
    const go_ref clientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));

    const go_ref supplyCfgRef = Supply(clientRef, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, 0);

    FreeGoMem(supplyCfgRef);
    FreeGoMem(clientRef);
    free(err);
}

static void test_block_query_client(void) {
    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));

    const go_ref clientRef = NewBlockQueryClient("http://127.0.0.1:26657", &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, 0);

    const go_ref blockRef = BlockQuery_ClientBlock(clientRef, NULL, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockRef, 0);

    FreeGoMem(blockRef);
    FreeGoMem(clientRef);
    free(err);
}

static void test_depinject_supply_many(void) {
    const go_ref eventsQueryClientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));

    const go_ref blockQueryClientRef = NewBlockQueryClient("http://127.0.0.1:26657", &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, 0);

    const int numSupplyRefs = 2;
    go_ref *toSupply = calloc(numSupplyRefs, sizeof(go_ref));
    toSupply[0] = eventsQueryClientRef;
    toSupply[1] = blockQueryClientRef;
    const go_ref supplyCfgRef = SupplyMany(toSupply, numSupplyRefs, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, 0);

    FreeGoMem(supplyCfgRef);
    free(toSupply);
    FreeGoMem(blockQueryClientRef);
    FreeGoMem(eventsQueryClientRef);
    free(err);
}


static void test_block_client(void) {
    const go_ref eventsQueryClientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));

    const go_ref blockQueryClientRef = NewBlockQueryClient("http://127.0.0.1:26657", &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, 0);

    const int numSupplyRefs = 2;
    go_ref *toSupply = calloc(numSupplyRefs, sizeof(go_ref));
    toSupply[0] = eventsQueryClientRef;
    toSupply[1] = blockQueryClientRef;
    const go_ref supplyCfgRef = SupplyMany(toSupply, numSupplyRefs, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, 0);

    const go_ref blockClientRef = NewBlockClient(supplyCfgRef, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockClientRef, 0);

    FreeGoMem(blockClientRef);
    FreeGoMem(supplyCfgRef);
    free(toSupply);
    FreeGoMem(blockQueryClientRef);
    FreeGoMem(eventsQueryClientRef);
    free(err);

    // TODO_IN_THIS_COMMIT: test block client methods...
}

static void test_tx_context(void) {
    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));

    const go_ref txCtxRef = NewTxContext("tcp://127.0.0.1:26657", &err);

    FreeGoMem(txCtxRef);
    free(err);
}


// TODO_IN_THIS_COMMIT: how to free encapsulated refs?
// TODO_IN_THIS_COMMIT: move...
static const go_ref getTxClientDeps(char **err) {
    const go_ref eventsQueryClientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, 0);

    const go_ref blockQueryClientRef = NewBlockQueryClient("http://127.0.0.1:26657", err);
    TEST_ASSERT_EQUAL_STRING("", *err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, 0);

    go_ref *toSupply = calloc(4, sizeof(go_ref));

    toSupply[0] = eventsQueryClientRef;
    toSupply[1] = blockQueryClientRef;
    go_ref supplyCfgRef = SupplyMany(toSupply, 2, err);
    TEST_ASSERT_EQUAL_STRING("", *err);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, 0);

    const go_ref blockClientRef = NewBlockClient(supplyCfgRef, err);
    TEST_ASSERT_EQUAL_STRING("", *err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockClientRef, 0);

    const go_ref txCtxRef = NewTxContext("tcp://127.0.0.1:26657", err);
    TEST_ASSERT_EQUAL_STRING("", *err);
    TEST_ASSERT_NOT_EQUAL_INT64(txCtxRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(txCtxRef, 0);

    toSupply[2] = blockClientRef;
    toSupply[3] = txCtxRef;
    supplyCfgRef = SupplyMany(toSupply, 4, err);
    TEST_ASSERT_EQUAL_STRING("", *err);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, 0);

    return supplyCfgRef;
}

static void test_tx_client(void) {
    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));

    const go_ref supplyCfgRef = getTxClientDeps(&err);
    TEST_ASSERT_EQUAL_STRING("", err);

    char *signingKeyName = "faucet";
    const go_ref txClientRef = NewTxClient(supplyCfgRef, signingKeyName, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, 0);

    FreeGoMem(txClientRef);
    FreeGoMem(supplyCfgRef);
    free(err);
}

static void test_sign_and_broadcast_any(void) {
    // NB: From faucet to app1.
    char *sendMsgAnyJSON = "{\"@type\":\"type.googleapis.com/cosmos.bank.v1beta1.MsgSend\", \"from_address\":\"pokt1awtlw5sjmw2f5lgj8ekdkaqezphgz88rdk93sk\", \"to_address\":\"pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4\", \"amount\":[{\"denom\":\"upokt\", \"amount\":\"100000000\"}]}";

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char *err = calloc(1024, sizeof(char));

    const go_ref supplyCfgRef = getTxClientDeps(&err);
    TEST_ASSERT_EQUAL_STRING("", err);

    char *signingKeyName = "faucet";
    const go_ref txClientRef = NewTxClient(supplyCfgRef, signingKeyName, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, 0);

    const go_ref errChRef = TxClient_SignAndBroadcastAny(txClientRef, sendMsgAnyJSON, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(errChRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(errChRef, 0);

    FreeGoMem(errChRef);
    FreeGoMem(txClientRef);
    FreeGoMem(supplyCfgRef);
    free(err);
}

// static void test_sign_and_broadcast(void) {
//     // NB: From faucet to app1.
//     uint8_t *sendMsgBz = // TODO: integrate buf.build/pokt-network/poktroll
//
//     // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
//     char *err = calloc(1024, sizeof(char));
//
//     const go_ref supplyCfgRef = getTxClientDeps(&err);
//     TEST_ASSERT_EQUAL_STRING("", err);
//
//     char *signingKeyName = "faucet";
//     const go_ref txClientRef = NewTxClient(supplyCfgRef, signingKeyName, &err);
//     TEST_ASSERT_EQUAL_STRING("", err);
//     TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, -1);
//     TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, 0);
//
//     const go_ref errChRef = TxClient_SignAndBroadcast(txClientRef, sendMsgBz, &err);
//     TEST_ASSERT_EQUAL_STRING("", err);
//     TEST_ASSERT_NOT_EQUAL_INT64(errChRef, -1);
//     TEST_ASSERT_NOT_EQUAL_INT64(errChRef, 0);
//
//     FreeGoMem(errChRef);
//     FreeGoMem(txClientRef);
//     FreeGoMem(supplyCfgRef);
//     free(err);
// }


void setUp(void) {
    // Code to run before each test (if any)
}

void tearDown(void) {
    // Code to run after each test (if any)
}

int main(void) {
    UNITY_BEGIN();            // Initialize Unity
    RUN_TEST(test_events_query_client);
    RUN_TEST(test_depinject_supply);
    RUN_TEST(test_block_query_client);
    RUN_TEST(test_depinject_supply_many);
    RUN_TEST(test_block_client);
    RUN_TEST(test_tx_context);
    RUN_TEST(test_tx_client);
    RUN_TEST(test_sign_and_broadcast_any);
    // RUN_TEST(test_sign_and_broadcast);
    return UNITY_END();       // End Unity and report results
}
