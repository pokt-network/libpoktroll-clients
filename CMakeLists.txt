# CMakeLists.txt
cmake_minimum_required(VERSION 3.29)
project(libpoktroll_clients VERSION 0.1.0)

# Add option for building tests
option(BUILD_TESTS "Build test programs" OFF)

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)

# Set platform-specific library extension and other variables
if(WIN32)
    set(LIB_EXTENSION "dll")
    set(LIB_PREFIX "")
    set(OS_NAME "windows")
elseif(APPLE)
    set(LIB_EXTENSION "dylib")
    set(LIB_PREFIX "lib")
    set(OS_NAME "darwin")
else()
    set(LIB_EXTENSION "so")
    set(LIB_PREFIX "lib")
    set(OS_NAME "linux")
endif()

# Define paths
set(CLIENTS_SHARED_LIB ${CMAKE_SOURCE_DIR}/cgo/build/${LIB_PREFIX}poktroll_clients)

set(LIBPOKTROLL_CLIENTS_SRC ${CMAKE_SOURCE_DIR}/src/context.c)

# Include directories
include_directories(
        ${CMAKE_SOURCE_DIR}/include
        ${CMAKE_SOURCE_DIR}/cgo/build
)

# Get list of Go source files
file(GLOB_RECURSE GO_SOURCES
        "${CMAKE_SOURCE_DIR}/cgo/*.go"
)

# Create a command that only executes when sources change
# Note: Windows needs different shell commands
if(WIN32)
    add_custom_command(
            OUTPUT ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION}
            COMMAND ${CMAKE_COMMAND} -E chdir ${CMAKE_SOURCE_DIR}/cgo
            cmd /c "if not exist go.mod go mod init libpoktroll/clients && go build -o ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION} -buildmode=c-shared ."
            DEPENDS ${GO_SOURCES}
            WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
            COMMENT "Building Go shared library..."
            VERBATIM
    )
else()
    add_custom_command(
            OUTPUT ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION}
            COMMAND ${CMAKE_COMMAND} -E chdir ${CMAKE_SOURCE_DIR}/cgo
            /bin/sh -c "if [ ! -f go.mod ]; then go mod init libpoktroll/clients; fi && \
                   go build -o ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION} -buildmode=c-shared ."
            DEPENDS ${GO_SOURCES}
            WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
            COMMENT "Building Go shared library..."
            VERBATIM
    )
endif()

# Add custom target that depends on the output
add_custom_target(build_go_shared_lib
        DEPENDS ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION}
)

# Add main library - make it SHARED explicitly
add_library(poktroll_clients SHARED ${LIBPOKTROLL_CLIENTS_SRC})
add_dependencies(poktroll_clients build_go_shared_lib)

# Link against the Go shared library only
target_link_libraries(poktroll_clients
        PUBLIC
        ${CLIENTS_SHARED_LIB}.${LIB_EXTENSION}
)

# Set library properties for versioning and naming
if(WIN32)
    set_target_properties(poktroll_clients PROPERTIES
            VERSION ${PROJECT_VERSION}
            PUBLIC_HEADER "${CMAKE_SOURCE_DIR}/include/context.h"
            OUTPUT_NAME "poktroll_clients"
            PREFIX ""
            IMPORT_PREFIX ""
    )
else()
    set_target_properties(poktroll_clients PROPERTIES
            VERSION ${PROJECT_VERSION}
            SOVERSION ${PROJECT_VERSION_MAJOR}
            PUBLIC_HEADER "${CMAKE_SOURCE_DIR}/include/context.h"
            OUTPUT_NAME "poktroll_clients"
            PREFIX "lib"
    )
endif()

# Include the test configuration only if BUILD_TESTS is enabled
if(BUILD_TESTS)
    include(${CMAKE_SOURCE_DIR}/cmake/TestConfig.cmake)
endif()

# Include the installer configuration
include(${CMAKE_SOURCE_DIR}/cmake/InstallerConfig.cmake)