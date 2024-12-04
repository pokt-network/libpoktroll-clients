cmake_minimum_required(VERSION 3.29)
project(libpoktroll_clients VERSION 0.1.0)  # Added version

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
        COMMAND ${CMAKE_COMMAND} -E chdir ${CMAKE_SOURCE_DIR}/cgo
        /bin/sh -c "if [ ! -f go.mod ]; then go mod init libpoktroll/clients; fi && \
                       go build -o ${CLIENTS_SHARED_LIB} -buildmode=c-shared ."
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

# Set library properties for versioning
set_target_properties(libpoktroll_clients PROPERTIES
        VERSION ${PROJECT_VERSION}
        SOVERSION ${PROJECT_VERSION_MAJOR}
        PUBLIC_HEADER "${CMAKE_SOURCE_DIR}/include/context.h"
)

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

# Installation configuration
include(GNUInstallDirs)

# Include CPack AFTER configuration. 👆
include(CPack)

install(TARGETS libpoktroll_clients
        LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
        PUBLIC_HEADER DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/poktroll
        COMPONENT library
)

# Install the Go shared library
install(FILES ${CLIENTS_SHARED_LIB}
        DESTINATION ${CMAKE_INSTALL_LIBDIR}
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
set(CPACK_PACKAGE_CONTACT "your@email.com")  # TODO: Update this
set(CPACK_RESOURCE_FILE_LICENSE "${CMAKE_SOURCE_DIR}/LICENSE")  # If you have a LICENSE file

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