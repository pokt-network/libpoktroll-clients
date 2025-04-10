cmake_minimum_required(VERSION 3.24)
project(libpoktroll_clients VERSION 0.2.0)

# Add option for building tests
option(BUILD_TESTS "Build test programs" ON)

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)

# Set platform-specific library extension and other variables
if (WIN32)
    set(LIB_EXTENSION "dll")
    set(LIB_PREFIX "")
    set(OS_NAME "windows")
elseif (APPLE)
    set(LIB_EXTENSION "dylib")
    set(LIB_PREFIX "lib")
    set(OS_NAME "darwin")
else ()
    set(LIB_EXTENSION "so")
    set(LIB_PREFIX "lib")
    set(OS_NAME "linux")
endif ()

# Define paths
set(CLIENTS_SHARED_LIB ${CMAKE_SOURCE_DIR}/cgo/build/${LIB_PREFIX}pocket_clients)

set(LIBPOCKET_CLIENTS_SRC ${CMAKE_SOURCE_DIR}/src/context.c)

# Include directories
include_directories(
        ${CMAKE_SOURCE_DIR}/include
        ${CMAKE_SOURCE_DIR}/cgo/build
)

# Get list of Go source files
file(GLOB_RECURSE GO_SOURCES
        "${CMAKE_SOURCE_DIR}/cgo/*.go"
)

# Determine Go architecture based on system
if (CMAKE_SYSTEM_PROCESSOR MATCHES "aarch64|ARM64")
    set(GO_ARCH "arm64")
elseif (CMAKE_SYSTEM_PROCESSOR MATCHES "x86_64")
    set(GO_ARCH "amd64")
else ()
    set(GO_ARCH ${CMAKE_SYSTEM_PROCESSOR})
endif ()

# Determine Go OS
if (WIN32)
    set(GO_OS "windows")
elseif (APPLE)
    set(GO_OS "darwin")
else ()
    set(GO_OS "linux")
endif ()

add_custom_command(
        OUTPUT ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION}
        COMMAND ${CMAKE_COMMAND} -E env GOOS=${GO_OS} GOARCH=${GO_ARCH} CGO_ENABLED=1
        ${CMAKE_COMMAND} -E chdir ${CMAKE_SOURCE_DIR}/cgo
        go build -v -o ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION} -buildmode=c-shared .
        DEPENDS ${GO_SOURCES}
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
        VERBATIM
)

# Add custom target that depends on the output
add_custom_target(build_go_shared_lib
        DEPENDS ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION}
)

add_custom_target(pocket_clients ALL
        COMMAND ${CMAKE_COMMAND} -E copy ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION} ${CMAKE_BINARY_DIR}/libpoktroll_clients.${LIB_EXTENSION}.${PROJECT_VERSION}
        COMMAND ${CMAKE_COMMAND} -E create_symlink libpoktroll_clients.${LIB_EXTENSION}.${PROJECT_VERSION} ${CMAKE_BINARY_DIR}/libpoktroll_clients.${LIB_EXTENSION}
        COMMAND ${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/cgo/build/libpoktroll_clients.h ${CMAKE_BINARY_DIR}/libpoktroll_clients.h
        DEPENDS ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION}
)

# Set library properties for versioning and naming
if (WIN32)
    set_target_properties(pocket_clients PROPERTIES
            VERSION ${PROJECT_VERSION}
            PUBLIC_HEADER "${CMAKE_SOURCE_DIR}/include/context.h"
            OUTPUT_NAME "pocket_clients"
            PREFIX ""
            IMPORT_PREFIX ""
    )
else ()
    set_target_properties(pocket_clients PROPERTIES
            VERSION ${PROJECT_VERSION}
            SOVERSION ${PROJECT_VERSION_MAJOR}
            PUBLIC_HEADER "${CMAKE_SOURCE_DIR}/include/context.h"
            OUTPUT_NAME "pocket_clients"
            PREFIX "lib"
    )
endif ()

# Include the test configuration only if BUILD_TESTS is enabled
if (BUILD_TESTS)
    include(${CMAKE_SOURCE_DIR}/cmake/TestConfig.cmake)
endif ()

# Include the installer configuration
include(${CMAKE_SOURCE_DIR}/cmake/InstallerConfig.cmake)