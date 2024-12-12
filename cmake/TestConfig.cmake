# cmake/TestConfig.cmake

# Find protobuf package and protobuf-c for tests only
find_package(Protobuf REQUIRED)
include(FindPkgConfig)
pkg_check_modules(PROTOBUFC REQUIRED libprotobuf-c)

# Set test-specific paths
set(LIBPOKTROLL_CLIENTS_TESTS ${CMAKE_SOURCE_DIR}/tests/test_main.c)
set(UNITY_DIR ${CMAKE_SOURCE_DIR}/libs/unity/src)
set(UNITY_SRC ${UNITY_DIR}/unity.c)
set(PROTO_GEN_DIR ${CMAKE_SOURCE_DIR}/gen)

# Find all generated protobuf source files (only needed for tests)
file(GLOB_RECURSE PROTO_SOURCES "${PROTO_GEN_DIR}/**/*.pb-c.c")

# Include Unity test framework directory and protobuf directories for tests
include_directories(
        ${UNITY_DIR}
        ${PROTO_GEN_DIR}
        ${Protobuf_INCLUDE_DIRS}
        ${PROTOBUFC_INCLUDE_DIRS}
)

# Add your test executable
add_executable(libpoktroll_clients_tests
        ${LIBPOKTROLL_CLIENTS_SRC}
        ${LIBPOKTROLL_CLIENTS_TESTS}
        ${UNITY_SRC}
        ${PROTO_SOURCES}
)

target_compile_options(libpoktroll_clients_tests PRIVATE -g)

# Link the test executable with your library and test-specific dependencies
target_link_libraries(libpoktroll_clients_tests
        PRIVATE
        ${CLIENTS_SHARED_LIB}.so
        ${Protobuf_LIBRARIES}
        ${PROTOBUFC_LIBRARIES}
)

# Make tests depend on the shared library being built
add_dependencies(libpoktroll_clients_tests build_go_shared_lib)

# Ensure the shared library directory is in the runtime path
set_target_properties(libpoktroll_clients_tests PROPERTIES
        BUILD_RPATH ${CMAKE_SOURCE_DIR}/cgo/build
)

# Enable testing and add the test
enable_testing()
add_test(NAME LibPoktrollTests COMMAND libpoktroll_clients_tests)