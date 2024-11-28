cmake_minimum_required(VERSION 3.29)
project(libpoktroll_clients)

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)

# Define paths
set(CLIENTS_SHARED_LIB ${CMAKE_SOURCE_DIR}/libs/poktroll/libclients.so)
set(UNITY_DIR ${CMAKE_SOURCE_DIR}/libs/unity/src)
set(UNITY_SRC ${UNITY_DIR}/unity.c)

# Include directories
include_directories(
        ${UNITY_DIR}
        ${CMAKE_SOURCE_DIR}/libs/poktroll
)

# Link directories for CMocka
link_directories(${CMOCKA_LIB_DIR})

# Add your main library (if applicable)
add_library(libpoktroll_clients SHARED library.cpp)

# Add your test executable
add_executable(libpoktroll_clients_tests test_main.c ${UNITY_SRC})

# Link the test executable with your library, CMocka, and the Poktroll shared object
target_link_libraries(libpoktroll_clients_tests
        PRIVATE
        libpoktroll_clients      # Your main library (optional)
        ${CLIENTS_SHARED_LIB}    # Path to the shared object clients.so
)

# Ensure the shared library directory is in the runtime path
set_target_properties(libpoktroll_clients_tests PROPERTIES
        BUILD_RPATH ${CMAKE_SOURCE_DIR}/libs/poktroll
)

# Enable testing and add the test
enable_testing()
add_test(NAME LibPoktrollTests COMMAND libpoktroll_clients_tests)
