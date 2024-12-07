#include <stdlib.h>
#include <string.h>

#include "libpoktroll_clients.h"
#include "unity.h"
#include "poktroll/application/tx.pb-c.h"

static const char* msg_stake_application_type_url = "poktroll.application.MsgStakeApplication";

static void test_events_query_client(void)
{
    const go_ref clientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char* err = calloc(1024, sizeof(char));
    const go_ref eventsBytesObsRef = EventsQueryClientEventsBytes(clientRef, "", &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsBytesObsRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsBytesObsRef, 0);

    FreeGoMem(eventsBytesObsRef);
    FreeGoMem(clientRef);
    free(err);
}

static void test_depinject_supply(void)
{
    const go_ref clientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(clientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char* err = calloc(1024, sizeof(char));

    const go_ref supplyCfgRef = Supply(clientRef, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(supplyCfgRef, 0);

    FreeGoMem(supplyCfgRef);
    FreeGoMem(clientRef);
    free(err);
}

static void test_block_query_client(void)
{
    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char* err = calloc(1024, sizeof(char));

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

static void test_depinject_supply_many(void)
{
    const go_ref eventsQueryClientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char* err = calloc(1024, sizeof(char));

    const go_ref blockQueryClientRef = NewBlockQueryClient("http://127.0.0.1:26657", &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, 0);

    const int numSupplyRefs = 2;
    go_ref* toSupply = calloc(numSupplyRefs, sizeof(go_ref));
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

static void test_block_client(void)
{
    const go_ref eventsQueryClientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, 0);

    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char* err = calloc(1024, sizeof(char));

    const go_ref blockQueryClientRef = NewBlockQueryClient("http://127.0.0.1:26657", &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, 0);

    const int numSupplyRefs = 2;
    go_ref* toSupply = calloc(numSupplyRefs, sizeof(go_ref));
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

static void test_tx_context(void)
{
    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char* err = calloc(1024, sizeof(char));

    const go_ref txCtxRef = NewTxContext("tcp://127.0.0.1:26657", &err);

    FreeGoMem(txCtxRef);
    free(err);
}

// TODO_IN_THIS_COMMIT: how to free encapsulated refs?
// TODO_IN_THIS_COMMIT: move...
static const go_ref getTxClientDeps(char** err)
{
    const go_ref eventsQueryClientRef = NewEventsQueryClient("ws://127.0.0.1:26657/websocket");
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(eventsQueryClientRef, 0);

    const go_ref blockQueryClientRef = NewBlockQueryClient("http://127.0.0.1:26657", err);
    TEST_ASSERT_EQUAL_STRING("", *err);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(blockQueryClientRef, 0);

    go_ref* toSupply = calloc(4, sizeof(go_ref));

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

static void test_tx_client(void)
{
    // DEV_NOTE: using calloc to ensure the error is an empty string if no error occurred.
    char* err = calloc(1024, sizeof(char));

    const go_ref supplyCfgRef = getTxClientDeps(&err);
    TEST_ASSERT_EQUAL_STRING("", err);

    char* signingKeyName = "faucet";
    const go_ref txClientRef = NewTxClient(supplyCfgRef, signingKeyName, &err);
    TEST_ASSERT_EQUAL_STRING("", err);
    TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, -1);
    TEST_ASSERT_NOT_EQUAL_INT64(txClientRef, 0);

    FreeGoMem(txClientRef);
    FreeGoMem(supplyCfgRef);
    free(err);
}

static void test_sign_and_broadcast_success(void)
{
    char* err = calloc(1024, sizeof(char));

    // Create and populate the message
    Poktroll__Application__MsgStakeApplication msg = POKTROLL__APPLICATION__MSG_STAKE_APPLICATION__INIT;

    // Stake application 3
    msg.address = "pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex";

    // Create and set the stake coin
    Cosmos__Base__V1beta1__Coin app_stake = COSMOS__BASE__V1BETA1__COIN__INIT;
    app_stake.denom = "upokt";
    app_stake.amount = "100000000"; // 1000 POKT
    msg.stake = &app_stake;

    // Create service configs
    Poktroll__Shared__ApplicationServiceConfig service = POKTROLL__SHARED__APPLICATION_SERVICE_CONFIG__INIT;
    service.service_id = "anvil";

    // Allocate array for services
    msg.n_services = 1;
    msg.services = malloc(sizeof(Poktroll__Shared__ApplicationServiceConfig*));
    msg.services[0] = &service;

    // Calculate the serialized size
    size_t msg_bz_len = poktroll__application__msg_stake_application__get_packed_size(&msg);

    // Allocate buffer and serialize
    uint8_t* msg_bz = malloc(msg_bz_len);
    poktroll__application__msg_stake_application__pack(&msg, msg_bz);

    // Get dependencies for tx client
    const go_ref supplyCfgRef = getTxClientDeps(&err);
    TEST_ASSERT_EQUAL_STRING("", err);

    const go_ref txClientRef = NewTxClient(supplyCfgRef, "app3", &err);
    TEST_ASSERT_EQUAL_STRING("", err);

    // Set up async context
    AsyncContext ctx;
    init_context(&ctx);

    AsyncOperation op = {
        .ctx = &ctx,
        .on_error = handle_error,
        .on_success = handle_success,
        .cleanup = cleanup_context
    };

    // Send the message
    serialized_proto serialized_msg = {
        .type_url = (uint8_t*)msg_stake_application_type_url,
        .type_url_length = strlen(msg_stake_application_type_url),
        .data = msg_bz,
        .data_length = msg_bz_len
    };
    const go_ref errChRef = TxClient_SignAndBroadcast(&op, txClientRef, &serialized_msg);

    // Wait for completion
    if (wait_for_completion(&ctx, 15000))
    {
        printf("Test passed: Operation completed successfully\n");
    }
    else
    {
        printf("Test failed: Operation timed out\n");
    }

    TEST_ASSERT_EQUAL_STRING("", ctx.error_msg);

    // Cleanup
    cleanup_context(&ctx);
    free(msg_bz);
    FreeGoMem(errChRef);
    FreeGoMem(txClientRef);
    FreeGoMem(supplyCfgRef);
    free(err);
}

static void test_sign_and_broadcast_sync_error(void)
{
    char* err = calloc(1024, sizeof(char));

    // Create and populate the message
    Poktroll__Application__MsgStakeApplication msg = POKTROLL__APPLICATION__MSG_STAKE_APPLICATION__INIT;

    // Set the application address
    msg.address = "pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex";

    // Create and set the stake coin with amount that should trigger error
    Cosmos__Base__V1beta1__Coin stake = COSMOS__BASE__V1BETA1__COIN__INIT;
    stake.denom = "upokt";
    stake.amount = "100000068"; // Amount equal to previous stake
    msg.stake = &stake;

    // Create service configs
    Poktroll__Shared__ApplicationServiceConfig service1 = POKTROLL__SHARED__APPLICATION_SERVICE_CONFIG__INIT;
    service1.service_id = "svc_123";

    // Allocate array for services
    msg.n_services = 1;
    msg.services = malloc(sizeof(Poktroll__Shared__ApplicationServiceConfig*));
    msg.services[0] = &service1;

    // Calculate the serialized size
    size_t msg_bz_len = poktroll__application__msg_stake_application__get_packed_size(&msg);

    // Allocate buffer and serialize
    uint8_t* msg_bz = malloc(msg_bz_len);
    poktroll__application__msg_stake_application__pack(&msg, msg_bz);

    // Get dependencies for tx client
    const go_ref supplyCfgRef = getTxClientDeps(&err);
    TEST_ASSERT_EQUAL_STRING("", err);

    // Construct the tx client with the WRONG signing key.
    const go_ref txClientRef = NewTxClient(supplyCfgRef, "pnf", &err);
    TEST_ASSERT_EQUAL_STRING("", err);

    // Set up async context
    AsyncContext ctx;
    init_context(&ctx);

    AsyncOperation op = {
        .ctx = &ctx,
        .on_error = handle_error,
        .on_success = handle_success,
        .cleanup = cleanup_context
    };

    // Send the message
    serialized_proto serialized_msg = {
        .type_url = (uint8_t*)msg_stake_application_type_url,
        .type_url_length = strlen(msg_stake_application_type_url),
        .data = msg_bz,
        .data_length = msg_bz_len
    };
    const go_ref errChRef = TxClient_SignAndBroadcast(&op, txClientRef, &serialized_msg);

    // Check that operation failed
    TEST_ASSERT_FALSE(ctx.success);

    TEST_ASSERT_EQUAL_STRING("tx intended signer does not match the given signer: pnf", ctx.error_msg);

    // Cleanup
    cleanup_context(&ctx);
    free(msg_bz);
    free(msg.services);
    FreeGoMem(errChRef);
    FreeGoMem(txClientRef);
    FreeGoMem(supplyCfgRef);
    free(err);
}

static void test_sign_and_broadcast_async_error(void)
{
    char* err = calloc(1024, sizeof(char));

    // Create and populate the message
    Poktroll__Application__MsgStakeApplication msg = POKTROLL__APPLICATION__MSG_STAKE_APPLICATION__INIT;

    // Set the application address
    msg.address = "pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex";

    // Create and set the stake coin with amount that should trigger error
    Cosmos__Base__V1beta1__Coin stake = COSMOS__BASE__V1BETA1__COIN__INIT;
    stake.denom = "upokt";
    stake.amount = "100000068"; // Amount equal to previous stake
    msg.stake = &stake;

    // Create service configs
    Poktroll__Shared__ApplicationServiceConfig service1 = POKTROLL__SHARED__APPLICATION_SERVICE_CONFIG__INIT;
    service1.service_id = "svc_123";

    // Allocate array for services
    msg.n_services = 1;
    msg.services = malloc(sizeof(Poktroll__Shared__ApplicationServiceConfig*));
    msg.services[0] = &service1;

    // Calculate the serialized size
    size_t msg_bz_len = poktroll__application__msg_stake_application__get_packed_size(&msg);

    // Allocate buffer and serialize
    uint8_t* msg_bz = malloc(msg_bz_len);
    poktroll__application__msg_stake_application__pack(&msg, msg_bz);

    // Get dependencies for tx client
    const go_ref supply_cfg_ref = getTxClientDeps(&err);
    TEST_ASSERT_EQUAL_STRING("", err);

    const go_ref tx_client_ref = NewTxClient(supply_cfg_ref, "app3", &err);
    TEST_ASSERT_EQUAL_STRING("", err);

    // Set up async context
    AsyncContext ctx;
    init_context(&ctx);

    AsyncOperation op = {
        .ctx = &ctx,
        .on_error = handle_error,
        .on_success = handle_success,
        .cleanup = cleanup_context
    };

    // Send the message
    serialized_proto serialized_msg = {
        // .type_url = "poktroll.application.MsgStakeApplication",
        .type_url = (uint8_t*)msg_stake_application_type_url,
        .type_url_length = strlen(msg_stake_application_type_url),
        .data = msg_bz,
        .data_length = msg_bz_len
    };
    const go_ref err_ch_ref = TxClient_SignAndBroadcast(&op, tx_client_ref, &serialized_msg);

    // Wait for completion and verify error
    wait_for_completion(&ctx, 15000);

    // Check that operation failed
    TEST_ASSERT_FALSE(ctx.success);

    // Verify error message contains expected parts, ignoring the dynamic hash
    char expected_error_msg[256] =
        "failed to execute message; message index: 0: rpc error: code = InvalidArgument desc = stake amount 100000068upokt must be higher than previous stake amount 100000068upokt: invalid";
    if (strstr(ctx.error_msg,
               expected_error_msg
    ) != NULL)
    {
        char fail_msg[1024];
        snprintf(fail_msg, sizeof(fail_msg),
                 "expected \"%s\" to contain \"%s\"",
                 ctx.error_msg, expected_error_msg);
        TEST_FAIL_MESSAGE(fail_msg);
    }

    // Cleanup
    cleanup_context(&ctx);
    free(msg_bz);
    free(msg.services);
    FreeGoMem(err_ch_ref);
    FreeGoMem(tx_client_ref);
    FreeGoMem(supply_cfg_ref);
    free(err);
}

void setUp(void)
{
    // Code to run before each test (if any)
}

void tearDown(void)
{
    // Code to run after each test (if any)
}

int main(void)
{
    UNITY_BEGIN(); // Initialize Unity
    RUN_TEST(test_events_query_client);
    RUN_TEST(test_depinject_supply);
    RUN_TEST(test_block_query_client);
    RUN_TEST(test_depinject_supply_many);
    RUN_TEST(test_block_client);
    RUN_TEST(test_tx_context);
    RUN_TEST(test_tx_client);
    RUN_TEST(test_sign_and_broadcast_sync_error);
    RUN_TEST(test_sign_and_broadcast_success);
    RUN_TEST(test_sign_and_broadcast_async_error);
    return UNITY_END(); // End Unity and report results
}
