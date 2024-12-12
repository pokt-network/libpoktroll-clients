# cmake/InstallerConfig.cmake

# Installation configuration
include(GNUInstallDirs)

# Determine system architecture
if(CMAKE_SYSTEM_PROCESSOR MATCHES "aarch64|ARM64")
    set(ARCH "arm64")
elseif(CMAKE_SYSTEM_PROCESSOR MATCHES "x86_64")
    set(ARCH "x86_64")
else()
    set(ARCH ${CMAKE_SYSTEM_PROCESSOR})
endif()

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
if(APPLE)
    set(SHARED_LIB_EXTENSION "dylib")
    set(PACKAGE_FILE_NAME "${PROJECT_NAME}-${PROJECT_VERSION}-${ARCH}-darwin")
else()
    set(SHARED_LIB_EXTENSION "so")
    set(PACKAGE_FILE_NAME "${PROJECT_NAME}-${PROJECT_VERSION}-${ARCH}-linux")
endif()

install(FILES ${CLIENTS_SHARED_LIB}.${SHARED_LIB_EXTENSION}
        DESTINATION ${CMAKE_INSTALL_LIBDIR}
        RENAME libpoktroll_clients.${SHARED_LIB_EXTENSION}.${PROJECT_VERSION}
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
set(CPACK_PACKAGE_CONTACT "bryanchriswhite+libpoktroll_clients@gmail.com")
set(CPACK_RESOURCE_FILE_LICENSE "${CMAKE_SOURCE_DIR}/LICENSE")

# Set package file name format
set(CPACK_PACKAGE_FILE_NAME ${PACKAGE_FILE_NAME})

# Platform-specific configuration
if(APPLE)
    # macOS specific settings
    set(CPACK_GENERATOR "productbuild;TGZ")
    set(CPACK_PACKAGING_INSTALL_PREFIX "/usr/local")

    # Set macOS package identifiers
    set(CPACK_OSX_PACKAGE_VERSION "${PROJECT_VERSION}")
    set(CPACK_BUNDLE_NAME "${PROJECT_NAME}")
    set(CPACK_BUNDLE_IDENTIFIER "network.pokt.clients")

    # Component-based installation for macOS
    set(CPACK_COMPONENTS_ALL library headers)
    set(CPACK_COMPONENT_LIBRARY_DISPLAY_NAME "POKT Network Client Library")
    set(CPACK_COMPONENT_HEADERS_DISPLAY_NAME "Development Headers")

    # macOS package naming
    set(CPACK_PRODUCTBUILD_DOMAINS TRUE)
    set(CPACK_DMG_FILE_NAME "${PACKAGE_FILE_NAME}")
    set(CPACK_PRODUCTBUILD_FILE_NAME "${PACKAGE_FILE_NAME}")

    # Dependencies for macOS (using Homebrew package names)
    set(CPACK_PRODUCTBUILD_DEPENDENCIES
            "protobuf-c"
    )

else()
    # Linux specific settings
    if(${CMAKE_SYSTEM_NAME} MATCHES "Linux")
        set(CPACK_GENERATOR "TGZ;DEB;RPM")

        # Debian-specific
        set(CPACK_DEBIAN_PACKAGE_MAINTAINER "Bryan White <bryanchriswhite+libpoktroll_clients@gmail.com>")
        set(CPACK_DEBIAN_PACKAGE_DEPENDS "libprotobuf-c-dev")
        set(CPACK_DEBIAN_PACKAGE_SECTION "libs")
        set(CPACK_DEBIAN_FILE_NAME "${PACKAGE_FILE_NAME}.deb")

        # RPM-specific
        set(CPACK_RPM_PACKAGE_REQUIRES "protobuf-c-devel")
        set(CPACK_RPM_PACKAGE_GROUP "Development/Libraries")
        set(CPACK_RPM_FILE_NAME "${PACKAGE_FILE_NAME}.rpm")

        # TGZ naming
        set(CPACK_ARCHIVE_FILE_NAME "${PACKAGE_FILE_NAME}")

        # Create pkg directory for Arch Linux
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
    endif()
endif()

# Include CPack
include(CPack)