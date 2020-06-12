#include "migrator.hpp"

using namespace eosio;

#define IMPORT extern "C" __attribute__((eosio_wasm_import))


IMPORT int32_t db_find_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
IMPORT int32_t db_store_i64(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id,  const void* data, uint32_t len);


void migrator::inject(const uint64_t scope,const uint64_t table,const uint64_t payer,const uint64_t id) {
    printf("scope: %d, table: %d, payer: %d,id: %d", scope,table,payer,id);
};
