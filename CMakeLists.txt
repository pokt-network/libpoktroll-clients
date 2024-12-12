# CMakeLists.txt
cmake_minimum_required(VERSION 3.29)
project(libpoktroll_clients VERSION 0.1.0)

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)

# Find protobuf package
find_package(Protobuf REQUIRED)
# Add pkg-config for finding protobuf-c
include(FindPkgConfig)
pkg_check_modules(PROTOBUFC REQUIRED libprotobuf-c)

# Define paths
set(CLIENTS_SHARED_LIB ${CMAKE_SOURCE_DIR}/cgo/build/libpoktroll_clients)
set(LIBPOKTROLL_CLIENTS_SRC ${CMAKE_SOURCE_DIR}/src/context.c)
set(PROTO_GEN_DIR ${CMAKE_SOURCE_DIR}/gen)

# Find all generated protobuf source files
file(GLOB_RECURSE PROTO_SOURCES "${PROTO_GEN_DIR}/**/*.pb-c.c")

# Include directories
include_directories(
        ${CMAKE_SOURCE_DIR}/include
        ${CMAKE_SOURCE_DIR}/cgo/build
        ${PROTO_GEN_DIR}
        ${Protobuf_INCLUDE_DIRS}
        ${PROTOBUFC_INCLUDE_DIRS}
)

# Get list of Go source files
file(GLOB_RECURSE GO_SOURCES
        "${CMAKE_SOURCE_DIR}/cgo/*.go"
)

# Create a command that only executes when sources change
add_custom_command(
        OUTPUT ${CLIENTS_SHARED_LIB}.so
        COMMAND ${CMAKE_COMMAND} -E chdir ${CMAKE_SOURCE_DIR}/cgo
        /bin/sh -c "if [ ! -f go.mod ]; then go mod init libpoktroll/clients; fi && \
               go build -o ${CLIENTS_SHARED_LIB}.so -buildmode=c-shared ."
        DEPENDS ${GO_SOURCES}
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
        COMMENT "Building Go shared library..."
        VERBATIM
)

# Add custom target that depends on the output
add_custom_target(build_go_shared_lib
        DEPENDS ${CLIENTS_SHARED_LIB}.so
)

# Add main library - make it SHARED explicitly
add_library(poktroll_clients SHARED ${LIBPOKTROLL_CLIENTS_SRC} ${PROTO_SOURCES})
add_dependencies(poktroll_clients build_go_shared_lib)

# Link against protobuf-c and the Go shared library
target_link_libraries(poktroll_clients
        PUBLIC
        ${PROTOBUFC_LIBRARIES}
        ${CLIENTS_SHARED_LIB}.so
)

# Set library properties for versioning and naming
set_target_properties(poktroll_clients PROPERTIES
        VERSION ${PROJECT_VERSION}
        SOVERSION ${PROJECT_VERSION_MAJOR}
        PUBLIC_HEADER "${CMAKE_SOURCE_DIR}/include/context.h"
        OUTPUT_NAME "poktroll_clients"
        PREFIX "lib"
)

# Include the test configuration
include(${CMAKE_SOURCE_DIR}/cmake/TestConfig.cmake)

# Include the installer configuration
include(${CMAKE_SOURCE_DIR}/cmake/InstallerConfig.cmake)