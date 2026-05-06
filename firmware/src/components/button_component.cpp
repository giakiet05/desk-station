#include "components/button_component.hpp"

#include <Arduino.h>

#include "core/outbox.hpp"
#include "protocol/message_schema.hpp"
#include "protocol/topics.hpp"

namespace deskstation {
namespace components {

ButtonComponent::ButtonComponent(uint8_t buttonPin, uint32_t debounceMs, core::Outbox& outbox)
    : buttonPin_(buttonPin), debounceMs_(debounceMs), outbox_(outbox) {}

const char* ButtonComponent::name() const { return "button"; }

void ButtonComponent::begin() {
  pinMode(buttonPin_, INPUT_PULLUP);
  const bool pressed = digitalRead(buttonPin_) == LOW;
  stablePressed_ = pressed;
  lastRawPressed_ = pressed;
  lastEdgeMs_ = millis();
}

void ButtonComponent::tick(uint32_t nowMs) {
  const bool rawPressed = digitalRead(buttonPin_) == LOW;
  if (rawPressed != lastRawPressed_) {
    lastRawPressed_ = rawPressed;
    lastEdgeMs_ = nowMs;
  }

  if (rawPressed == stablePressed_) {
    return;
  }

  if ((nowMs - lastEdgeMs_) < debounceMs_) {
    return;
  }

  stablePressed_ = rawPressed;

  protocol::MessageEnvelope envelope{
      protocol::topics::kProtocolVersion,
      protocol::topics::kEvtButton,
      protocol::MessageType::Event,
      ++sequence_,
      nowMs,
  };

  protocol::ButtonEventPayload payload{stablePressed_};

  char line[protocol::kMaxJsonLineLength] = {};
  if (protocol::encodeButtonEventJson(envelope, payload, line, sizeof(line))) {
    outbox_.push(line);
  }
}

} // namespace components
} // namespace deskstation
