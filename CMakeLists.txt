cmake_minimum_required(VERSION 3.29)
project(libpoktroll_clients)

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)

# Define paths
set(CLIENTS_SHARED_LIB ${CMAKE_SOURCE_DIR}/libs/poktroll/libclients.so)
set(CMOCKA_INCLUDE_DIR ${CMAKE_SOURCE_DIR}/libs/cmocka-1.1.0/include)
set(CMOCKA_LIB_DIR ${CMAKE_SOURCE_DIR}/libs/cmocka-1.1.0/build/src)

# Include directories for CMocka
include_directories(${CMOCKA_INCLUDE_DIR})

include_directories(${CMAKE_SOURCE_DIR}/libs/poktroll)


# Link directories for CMocka
link_directories(${CMOCKA_LIB_DIR})

# Add your main library (if applicable)
add_library(libpoktroll_clients SHARED library.cpp)

# Add your test executable
add_executable(libpoktroll_clients_tests test_main.c)

# Link the test executable with your library, CMocka, and the Poktroll shared object
target_link_libraries(libpoktroll_clients_tests
        PRIVATE
        libpoktroll_clients      # Your main library (optional)
        cmocka                   # CMocka testing framework
        ${CLIENTS_SHARED_LIB}    # Path to the shared object clients.so
)

# Ensure the shared library directory is in the runtime path
set_target_properties(libpoktroll_clients_tests PROPERTIES
        BUILD_RPATH ${CMAKE_SOURCE_DIR}/lib/poktroll
)

# Enable testing and add the test
enable_testing()
add_test(NAME LibPoktrollTests COMMAND libpoktroll_clients_tests)
