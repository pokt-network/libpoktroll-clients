// gas.h
#ifndef GAS_H
#define GAS_H

#include <stdint.h>

typedef struct gas_settings {
    uint64_t gas_limit;
    bool simulate;
    char *gas_prices;
    double gas_adjustment;
    char *fees;
} gas_settings;

#endif // GAS_H
