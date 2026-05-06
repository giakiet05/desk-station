#pragma once

#include <cstdint>

#include "core/component.hpp"
#include "protocol/message_schema.hpp"

namespace deskstation {
namespace components {

class LEDComponent : public core::IComponent {
public:
  LEDComponent(uint8_t redPin, uint8_t yellowPin, uint8_t greenPin);

  const char* name() const override;
  void begin() override;
  void tick(uint32_t nowMs) override;

  void setState(bool redOn, bool yellowOn, bool greenOn);
  bool applyCommand(const protocol::LedCommandPayload& payload);

private:
  uint8_t redPin_;
  uint8_t yellowPin_;
  uint8_t greenPin_;
  bool redOn_ = false;
  bool yellowOn_ = false;
  bool greenOn_ = false;
};

} // namespace components
} // namespace deskstation
