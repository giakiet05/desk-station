#pragma once

#include <cstdint>

#include "core/component.hpp"

namespace deskstation {
namespace core {
class Outbox;
}
namespace components {

class ButtonComponent : public core::IComponent {
public:
  ButtonComponent(uint8_t buttonPin, uint32_t debounceMs, core::Outbox& outbox);

  const char* name() const override;
  void begin() override;
  void tick(uint32_t nowMs) override;

private:
  uint8_t buttonPin_;
  uint32_t debounceMs_;
  core::Outbox& outbox_;
  uint32_t sequence_ = 0;
  bool stablePressed_ = false;
  bool lastRawPressed_ = false;
  uint32_t lastEdgeMs_ = 0;
};

} // namespace components
} // namespace deskstation
