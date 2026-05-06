#include "protocol/message_schema.hpp"

#include <cctype>
#include <cstdio>
#include <cstdlib>
#include <cstring>

namespace deskstation {
namespace protocol {

namespace {

const char* typeToString(MessageType type) {
  switch (type) {
  case MessageType::Event:
    return "event";
  case MessageType::Command:
    return "command";
  case MessageType::Ack:
    return "ack";
  default:
    return "unknown";
  }
}

bool formatOneDecimal(float value, char* outBuffer, size_t outBufferSize) {
  if (outBuffer == nullptr || outBufferSize == 0) {
    return false;
  }

  const bool negative = value < 0.0F;
  const float absValue = negative ? -value : value;
  int whole = static_cast<int>(absValue);
  int frac = static_cast<int>((absValue - static_cast<float>(whole)) * 10.0F + 0.5F);
  if (frac >= 10) {
    ++whole;
    frac = 0;
  }

  const int written = std::snprintf(outBuffer,
                                    outBufferSize,
                                    negative ? "-%d.%d" : "%d.%d",
                                    whole,
                                    frac);
  return written > 0 && static_cast<size_t>(written) < outBufferSize;
}

const char* findKey(const char* jsonLine, const char* key) {
  if (jsonLine == nullptr || key == nullptr) {
    return nullptr;
  }
  return std::strstr(jsonLine, key);
}

const char* findValueStart(const char* keyPosition) {
  if (keyPosition == nullptr) {
    return nullptr;
  }

  const char* colon = std::strchr(keyPosition, ':');
  if (colon == nullptr) {
    return nullptr;
  }

  ++colon;
  while (*colon != '\0' && std::isspace(static_cast<unsigned char>(*colon)) != 0) {
    ++colon;
  }
  return colon;
}

bool parseBoolValue(const char* jsonLine, const char* key, bool& outValue) {
  const char* keyPos = findKey(jsonLine, key);
  const char* value = findValueStart(keyPos);
  if (value == nullptr) {
    return false;
  }

  if (std::strncmp(value, "true", 4) == 0) {
    outValue = true;
    return true;
  }

  if (std::strncmp(value, "false", 5) == 0) {
    outValue = false;
    return true;
  }

  if (*value == '1') {
    outValue = true;
    return true;
  }

  if (*value == '0') {
    outValue = false;
    return true;
  }

  return false;
}

bool parseUint32Value(const char* jsonLine, const char* key, uint32_t& outValue) {
  const char* keyPos = findKey(jsonLine, key);
  const char* value = findValueStart(keyPos);
  if (value == nullptr) {
    return false;
  }

  char* end = nullptr;
  const unsigned long parsed = std::strtoul(value, &end, 10);
  if (end == value) {
    return false;
  }

  outValue = static_cast<uint32_t>(parsed);
  return true;
}

} // namespace

bool encodeIrEventJson(const MessageEnvelope& envelope,
                       const IrEventPayload& payload,
                       char* outBuffer,
                       size_t outBufferSize) {
  if (outBuffer == nullptr || outBufferSize == 0) {
    return false;
  }

  const int written = std::snprintf(
      outBuffer,
      outBufferSize,
      "{\"v\":\"%s\",\"topic\":\"%s\",\"type\":\"%s\","
      "\"seq\":%lu,\"ts\":%lu,\"payload\":{\"raw_code\":\"0x%08lX\","
      "\"address\":%u,\"command\":%u,\"is_repeat\":%s}}",
      envelope.version,
      envelope.topic,
      typeToString(envelope.type),
      static_cast<unsigned long>(envelope.sequence),
      static_cast<unsigned long>(envelope.timestampMs),
      static_cast<unsigned long>(payload.rawCode),
      payload.address,
      payload.command,
      payload.isRepeat ? "true" : "false");

  return written > 0 && static_cast<size_t>(written) < outBufferSize;
}

bool encodeDht11EventJson(const MessageEnvelope& envelope,
                          const Dht11EventPayload& payload,
                          char* outBuffer,
                          size_t outBufferSize) {
  if (outBuffer == nullptr || outBufferSize == 0) {
    return false;
  }

  char temperatureBuffer[16] = {};
  char humidityBuffer[16] = {};
  if (!formatOneDecimal(payload.temperature, temperatureBuffer, sizeof(temperatureBuffer))) {
    return false;
  }
  if (!formatOneDecimal(payload.humidity, humidityBuffer, sizeof(humidityBuffer))) {
    return false;
  }

  const int written = std::snprintf(
      outBuffer,
      outBufferSize,
      "{\"v\":\"%s\",\"topic\":\"%s\",\"type\":\"%s\","
      "\"seq\":%lu,\"ts\":%lu,\"payload\":{\"temperature\":%s,"
      "\"humidity\":%s,\"valid\":%s}}",
      envelope.version,
      envelope.topic,
      typeToString(envelope.type),
      static_cast<unsigned long>(envelope.sequence),
      static_cast<unsigned long>(envelope.timestampMs),
      temperatureBuffer,
      humidityBuffer,
      payload.valid ? "true" : "false");

  return written > 0 && static_cast<size_t>(written) < outBufferSize;
}

bool encodeButtonEventJson(const MessageEnvelope& envelope,
                           const ButtonEventPayload& payload,
                           char* outBuffer,
                           size_t outBufferSize) {
  if (outBuffer == nullptr || outBufferSize == 0) {
    return false;
  }

  const int written = std::snprintf(
      outBuffer,
      outBufferSize,
      "{\"v\":\"%s\",\"topic\":\"%s\",\"type\":\"%s\","
      "\"seq\":%lu,\"ts\":%lu,\"payload\":{\"pressed\":%s}}",
      envelope.version,
      envelope.topic,
      typeToString(envelope.type),
      static_cast<unsigned long>(envelope.sequence),
      static_cast<unsigned long>(envelope.timestampMs),
      payload.pressed ? "true" : "false");

  return written > 0 && static_cast<size_t>(written) < outBufferSize;
}

bool encodeAckJson(const MessageEnvelope& envelope,
                   uint32_t ackedCommandId,
                   bool ok,
                   const char* reason,
                   char* outBuffer,
                   size_t outBufferSize) {
  if (outBuffer == nullptr || outBufferSize == 0 || reason == nullptr) {
    return false;
  }

  const int written = std::snprintf(outBuffer,
                                    outBufferSize,
                                    "{\"v\":\"%s\",\"topic\":\"%s\",\"type\":\"%s\","
                                    "\"seq\":%lu,\"ts\":%lu,"
                                    "\"payload\":{\"cmd_id\":%lu,\"ok\":%s,"
                                    "\"reason\":\"%s\"}}",
                                    envelope.version,
                                    envelope.topic,
                                    typeToString(envelope.type),
                                    static_cast<unsigned long>(envelope.sequence),
                                    static_cast<unsigned long>(envelope.timestampMs),
                                    static_cast<unsigned long>(ackedCommandId),
                                    ok ? "true" : "false",
                                    reason);

  return written > 0 && static_cast<size_t>(written) < outBufferSize;
}

bool decodeLedCommandJson(const char* jsonLine, ParsedLedCommand& outCommand) {
  outCommand = {};
  if (jsonLine == nullptr) {
    return false;
  }

  const bool hasTopic = std::strstr(jsonLine, "\"topic\":\"cmd.led\"") != nullptr ||
                        std::strstr(jsonLine, "\"topic\": \"cmd.led\"") != nullptr;
  if (!hasTopic) {
    return false;
  }

  bool redOn = false;
  bool yellowOn = false;
  bool greenOn = false;
  uint32_t commandId = 0;

  if (!parseBoolValue(jsonLine, "\"red_on\"", redOn)) {
    return false;
  }
  if (!parseBoolValue(jsonLine, "\"yellow_on\"", yellowOn)) {
    return false;
  }
  if (!parseBoolValue(jsonLine, "\"green_on\"", greenOn)) {
    return false;
  }
  if (!parseUint32Value(jsonLine, "\"cmd_id\"", commandId)) {
    return false;
  }

  outCommand.valid = true;
  outCommand.cmdId = commandId;
  outCommand.payload.redOn = redOn;
  outCommand.payload.yellowOn = yellowOn;
  outCommand.payload.greenOn = greenOn;
  return true;
}

} // namespace protocol
} // namespace deskstation
