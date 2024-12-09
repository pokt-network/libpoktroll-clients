cmake_minimum_required(VERSION 3.29)
project(libpoktroll_clients VERSION 0.1.0)  # Added version

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)

# Find protobuf package
find_package(Protobuf REQUIRED)

find_package(PkgConfig REQUIRED)
pkg_check_modules(PROTOBUF_C REQUIRED libprotobuf-c)

# Define paths
set(CLIENTS_SHARED_LIB ${CMAKE_SOURCE_DIR}/cgo/build/libpoktroll_clients)
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
        ${PROTOBUF_C_INCLUDE_DIRS}
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

# Set library properties for versioning and naming
set_target_properties(poktroll_clients PROPERTIES
        VERSION ${PROJECT_VERSION}
        SOVERSION ${PROJECT_VERSION_MAJOR}
        PUBLIC_HEADER "${CMAKE_SOURCE_DIR}/include/context.h"
        OUTPUT_NAME "poktroll_clients"
        PREFIX "lib"
)

target_link_libraries(poktroll_clients
        PRIVATE
        ${CLIENTS_SHARED_LIB}.so
        ${Protobuf_C_LIBRARIES}
        ${PROTOBUF_C_LIBRARIES}
        protobuf-c
)

# Add your test executable
add_executable(libpoktroll_clients_tests
        ${LIBPOKTROLL_CLIENTS_SRC}
        ${LIBPOKTROLL_CLIENTS_TESTS}
        ${UNITY_SRC}
        ${PROTO_SOURCES}
        ${PROTOBUF_LIBRARIES}
        ${Protobuf_C_LIBRARIES}
)

target_compile_options(libpoktroll_clients_tests PRIVATE -g)

# Link the test executable with your library and the Poktroll shared object
target_link_libraries(libpoktroll_clients_tests
        PRIVATE
        ${CLIENTS_SHARED_LIB}.so    # Updated path to the shared object
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

# Installation configuration
include(GNUInstallDirs)

# Install targets with proper naming and symlinks
install(TARGETS poktroll_clients
        LIBRARY
        DESTINATION ${CMAKE_INSTALL_LIBDIR}
        NAMELINK_SKIP
        PUBLIC_HEADER
        DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/poktroll
        COMPONENT library
)

# Install the Go shared library with proper naming
install(FILES ${CLIENTS_SHARED_LIB}.so
        DESTINATION ${CMAKE_INSTALL_LIBDIR}
        RENAME libpoktroll_clients.so.${PROJECT_VERSION}
        COMPONENT library
)

# Generate and install pkg-config file
configure_file(
        ${CMAKE_SOURCE_DIR}/libpoktroll_clients.pc.in
        ${CMAKE_BINARY_DIR}/libpoktroll_clients.pc
        @ONLY
)
install(FILES ${CMAKE_BINARY_DIR}/libpoktroll_clients.pc
        DESTINATION ${CMAKE_INSTALL_LIBDIR}/pkgconfig
)

# CPack configuration
set(CPACK_PACKAGE_NAME "${PROJECT_NAME}")
set(CPACK_PACKAGE_VERSION "${PROJECT_VERSION}")
set(CPACK_PACKAGE_DESCRIPTION_SUMMARY "POKT Network Client Library")
set(CPACK_PACKAGE_VENDOR "POKT Network")

# Set the maintainer contact information (required for Debian packages)
set(CPACK_DEBIAN_PACKAGE_MAINTAINER "Bryan White <bryanchriswhite+libpoktroll_clients@gmail.com>")
set(CPACK_PACKAGE_CONTACT "bryanchriswhite+libpoktroll_clients@gmail.com")

set(CPACK_RESOURCE_FILE_LICENSE "${CMAKE_SOURCE_DIR}/LICENSE")

# Enable multiple package formats
set(CPACK_GENERATOR "TGZ;DEB;RPM")

# Debian-specific
set(CPACK_DEBIAN_PACKAGE_DEPENDS "libprotobuf-c-dev")
set(CPACK_DEBIAN_PACKAGE_SECTION "libs")
set(CPACK_DEBIAN_FILE_NAME DEB-DEFAULT)  # Use standardized naming

# RPM-specific
set(CPACK_RPM_PACKAGE_REQUIRES "protobuf-c-devel")
set(CPACK_RPM_PACKAGE_GROUP "Development/Libraries")
set(CPACK_RPM_FILE_NAME RPM-DEFAULT)  # Use standardized naming

# Include CPack
include(CPack)

# Create pkg directory
add_custom_command(
        OUTPUT ${CMAKE_BINARY_DIR}/pkg
        COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_BINARY_DIR}/pkg
)

# Generate PKGBUILD file
configure_file(
        ${CMAKE_SOURCE_DIR}/PKGBUILD.in
        ${CMAKE_BINARY_DIR}/PKGBUILD
        @ONLY
)

# Custom target for building Arch package
add_custom_target(pkgbuild
        COMMAND bash ${CMAKE_SOURCE_DIR}/scripts/build_pkg.sh ${CMAKE_BINARY_DIR} ${PROJECT_VERSION}
        COMMENT "Generating Arch Linux package"
)