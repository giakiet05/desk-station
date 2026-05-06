#pragma once

#include <cstddef>
#include <cstdint>

namespace deskstation {
namespace protocol {

enum class MessageType : uint8_t {
  Event = 0,
  Command = 1,
  Ack = 2
};

struct MessageEnvelope {
  const char* version;
  const char* topic;
  MessageType type;
  uint32_t sequence;
  uint32_t timestampMs;
};

struct IrEventPayload {
  uint32_t rawCode;
  uint16_t address;
  uint16_t command;
  bool isRepeat;
};

struct Dht11EventPayload {
  float temperature;
  float humidity;
  bool valid;
};

struct ButtonEventPayload {
  bool pressed;
};

struct LedCommandPayload {
  bool redOn;
  bool yellowOn;
  bool greenOn;
};

struct ParsedLedCommand {
  bool valid;
  uint32_t cmdId;
  LedCommandPayload payload;
};

constexpr size_t kMaxJsonLineLength = 256;

bool encodeIrEventJson(const MessageEnvelope& envelope,
                       const IrEventPayload& payload,
                       char* outBuffer,
                       size_t outBufferSize);

bool encodeDht11EventJson(const MessageEnvelope& envelope,
                          const Dht11EventPayload& payload,
                          char* outBuffer,
                          size_t outBufferSize);

bool encodeButtonEventJson(const MessageEnvelope& envelope,
                           const ButtonEventPayload& payload,
                           char* outBuffer,
                           size_t outBufferSize);

bool encodeAckJson(const MessageEnvelope& envelope,
                   uint32_t ackedCommandId,
                   bool ok,
                   const char* reason,
                   char* outBuffer,
                   size_t outBufferSize);

bool decodeLedCommandJson(const char* jsonLine, ParsedLedCommand& outCommand);

} // namespace protocol
} // namespace deskstation
