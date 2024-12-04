cmake_minimum_required(VERSION 3.29)
project(libpoktroll_clients)

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)

# Find protobuf package
find_package(Protobuf REQUIRED)

# Define paths
set(CLIENTS_SHARED_LIB ${CMAKE_SOURCE_DIR}/cgo/build/libclients)
set(LIBPOKTROLL_CLIENTS_SRC ${CMAKE_SOURCE_DIR}/src/context.c)
set(LIBPOKTROLL_CLIENTS_TESTS ${CMAKE_SOURCE_DIR}/tests/test_main.c)
set(UNITY_DIR ${CMAKE_SOURCE_DIR}/libs/unity/src)
set(UNITY_SRC ${UNITY_DIR}/unity.c)
set(PROTO_GEN_DIR ${CMAKE_SOURCE_DIR}/gen)

# Find all generated protobuf source files
file(GLOB_RECURSE PROTO_SOURCES "${PROTO_GEN_DIR}/**/*.pb-c.c")

# Include directories
include_directories(
        ${UNITY_DIR}
        ${CMAKE_SOURCE_DIR}/include
        ${CMAKE_SOURCE_DIR}/cgo/build
        ${PROTO_GEN_DIR}
        ${Protobuf_INCLUDE_DIRS}
)

# Get list of Go source files
file(GLOB_RECURSE GO_SOURCES
        "${CMAKE_SOURCE_DIR}/cgo/*.go"
)

# Create a command that only executes when sources change
add_custom_command(
        OUTPUT ${CLIENTS_SHARED_LIB}
        COMMAND go build -o ${CLIENTS_SHARED_LIB} -buildmode=c-shared ${CMAKE_SOURCE_DIR}/cgo
        DEPENDS ${GO_SOURCES}
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
        COMMENT "Building Go shared library..."
        VERBATIM
)

# Add custom target that depends on the output
add_custom_target(build_go_shared_lib
        DEPENDS ${CLIENTS_SHARED_LIB}
)

# Add main library
add_library(libpoktroll_clients ${LIBPOKTROLL_CLIENTS_SRC} ${PROTO_SOURCES})
add_dependencies(libpoktroll_clients build_go_shared_lib)

# Add your test executable
add_executable(libpoktroll_clients_tests
        ${LIBPOKTROLL_CLIENTS_SRC}
        ${LIBPOKTROLL_CLIENTS_TESTS}
        ${UNITY_SRC}
        ${PROTO_SOURCES}
)

target_compile_options(libpoktroll_clients_tests PRIVATE -g)

# Link the test executable with your library and the Poktroll shared object
target_link_libraries(libpoktroll_clients_tests
        PRIVATE
        ${CLIENTS_SHARED_LIB}    # Path to the shared object clients.so
        ${Protobuf_LIBRARIES}
        protobuf-c              # Add protobuf-c library
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

# TODO: add a custom target to install unity for testing.