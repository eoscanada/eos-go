#include <string>
#include <vector>
#include <eosio/eosio.hpp>

using namespace eosio;
using std::string;  

class [[eosio::contract]]  migrator : public contract {
  public:
    migrator(name receiver, name code, datastream<const char*> ds)
      :contract(receiver, code, ds)
      {}
  
    // Actions      
    [[eosio::action]]
    void inject(const uint64_t scope,const uint64_t table,const uint64_t payer,const uint64_t id);
  
  private:
};