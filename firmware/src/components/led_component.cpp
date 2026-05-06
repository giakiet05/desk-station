#include "components/led_component.hpp"

#include <Arduino.h>

namespace deskstation {
namespace components {

LEDComponent::LEDComponent(uint8_t redPin, uint8_t yellowPin, uint8_t greenPin)
    : redPin_(redPin), yellowPin_(yellowPin), greenPin_(greenPin) {
  redOn_ = true;
  yellowOn_ = true;
  greenOn_ = true;
}

const char* LEDComponent::name() const { return "led"; }

void LEDComponent::begin() {
  pinMode(redPin_, OUTPUT);
  pinMode(yellowPin_, OUTPUT);
  pinMode(greenPin_, OUTPUT);
  setState(redOn_, yellowOn_, greenOn_);
}

void LEDComponent::tick(uint32_t nowMs) { (void)nowMs; }

void LEDComponent::setState(bool redOn, bool yellowOn, bool greenOn) {
  redOn_ = redOn;
  yellowOn_ = yellowOn;
  greenOn_ = greenOn;

  digitalWrite(redPin_, redOn_ ? HIGH : LOW);
  digitalWrite(yellowPin_, yellowOn_ ? HIGH : LOW);
  digitalWrite(greenPin_, greenOn_ ? HIGH : LOW);
}

bool LEDComponent::applyCommand(const protocol::LedCommandPayload& payload) {
  setState(payload.redOn, payload.yellowOn, payload.greenOn);
  return true;
}

} // namespace components
} // namespace deskstation
