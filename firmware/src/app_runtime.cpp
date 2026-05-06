#include <Arduino.h>
#include <cstdio>

#include "components/button_component.hpp"
#include "components/dht11_component.hpp"
#include "components/ir_receiver_component.hpp"
#include "components/led_component.hpp"
#include "config/pins.hpp"
#include "config/timing.hpp"
#include "core/component.hpp"
#include "core/outbox.hpp"
#include "protocol/message_schema.hpp"
#include "protocol/topics.hpp"
#include "transport/serial_reader.hpp"
#include "transport/serial_writer.hpp"

namespace deskstation {
namespace {

core::Outbox gOutbox;

char gSerialReadBuffer[protocol::kMaxJsonLineLength] = {};
transport::SerialWriter gSerialWriter(Serial);
transport::SerialReader gSerialReader(Serial, gSerialReadBuffer, sizeof(gSerialReadBuffer));

components::LEDComponent gLedComponent(
    config::kRedLedPin, config::kYellowLedPin, config::kGreenLedPin);
components::IRReceiverComponent gIrComponent(config::kIrReceivePin, gOutbox);
components::DHT11Component gDht11Component(
    config::kDht11Pin, config::kDht11PublishIntervalMs, gOutbox);
components::ButtonComponent gButtonComponent(
    config::kButtonPin, config::kButtonDebounceMs, gOutbox);

core::IComponent* gComponents[] = {
    &gLedComponent,
    &gIrComponent,
    &gDht11Component,
    &gButtonComponent,
};

uint32_t gAckSequence = 0;

void publishAck(uint32_t cmdId, bool ok, const char* reason, uint32_t nowMs) {
  protocol::MessageEnvelope envelope{
      protocol::topics::kProtocolVersion,
      protocol::topics::kAckCommand,
      protocol::MessageType::Ack,
      ++gAckSequence,
      nowMs,
  };

  char line[protocol::kMaxJsonLineLength] = {};
  if (protocol::encodeAckJson(envelope, cmdId, ok, reason, line, sizeof(line))) {
    gOutbox.push(line);
  }
}

void handleIncomingLine(const char* line, uint32_t nowMs) {
  protocol::ParsedLedCommand command{};
  if (protocol::decodeLedCommandJson(line, command)) {
    const bool applied = gLedComponent.applyCommand(command.payload);
    publishAck(command.cmdId, applied, applied ? "ok" : "apply_failed", nowMs);
    return;
  }
}

void processIncoming(uint32_t nowMs) {
  const char* line = nullptr;
  while (gSerialReader.pollLine(line)) {
    handleIncomingLine(line, nowMs);
  }
}

void flushOutbox() {
  core::OutboxMessage message{};
  while (gOutbox.pop(message)) {
    gSerialWriter.writeLine(message.line);
  }
  gSerialWriter.flush();
}

void publishBootEvent(uint32_t nowMs) {
  char line[protocol::kMaxJsonLineLength] = {};
  const int written = snprintf(
      line,
      sizeof(line),
      "{\"v\":\"%s\",\"topic\":\"%s\",\"type\":\"event\","
      "\"seq\":1,\"ts\":%lu,\"payload\":{\"status\":\"ready\"}}",
      protocol::topics::kProtocolVersion,
      protocol::topics::kEvtDevice,
      static_cast<unsigned long>(nowMs));
  if (written > 0 && static_cast<size_t>(written) < sizeof(line)) {
    gOutbox.push(line);
  }
}

} // namespace

void appSetup() {
  Serial.begin(115200);
  const uint32_t startWaitMs = millis();
  while (!Serial && (millis() - startWaitMs) < 1500U) {
    delay(1);
  }

  for (core::IComponent* component : gComponents) {
    component->begin();
  }

  publishBootEvent(millis());
  flushOutbox();
}

void appLoop() {
  const uint32_t nowMs = millis();

  processIncoming(nowMs);

  for (core::IComponent* component : gComponents) {
    component->tick(nowMs);
  }

  flushOutbox();
  delay(config::kMainLoopDelayMs);
}

} // namespace deskstation
