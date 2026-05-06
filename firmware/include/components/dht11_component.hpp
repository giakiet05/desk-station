#pragma once

#include <cstdint>

#include "core/component.hpp"

class DHT;

namespace deskstation {
namespace core {
class Outbox;
}
namespace components {

class DHT11Component : public core::IComponent {
public:
  DHT11Component(uint8_t dataPin,
                 uint32_t publishIntervalMs,
                 core::Outbox& outbox);

  const char* name() const override;
  void begin() override;
  void tick(uint32_t nowMs) override;

private:
  uint8_t dataPin_;
  uint32_t publishIntervalMs_;
  core::Outbox& outbox_;
  DHT* dht_ = nullptr;
  uint32_t sequence_ = 0;
  uint32_t lastPublishMs_ = 0;
};

} // namespace components
} // namespace deskstation
