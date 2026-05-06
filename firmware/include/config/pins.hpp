#pragma once

#include <cstdint>

#ifndef DS_PIN_IR
#define DS_PIN_IR 3
#endif

#ifndef DS_PIN_DHT11
#define DS_PIN_DHT11 4
#endif

#ifndef DS_PIN_BUTTON
#define DS_PIN_BUTTON 15
#endif

#ifndef DS_PIN_LED_RED
#define DS_PIN_LED_RED 28
#endif

#ifndef DS_PIN_LED_YELLOW
#define DS_PIN_LED_YELLOW 27
#endif

#ifndef DS_PIN_LED_GREEN
#define DS_PIN_LED_GREEN 26
#endif

namespace deskstation {
namespace config {

constexpr uint8_t kIrReceivePin = DS_PIN_IR;
constexpr uint8_t kDht11Pin = DS_PIN_DHT11;
constexpr uint8_t kButtonPin = DS_PIN_BUTTON;

constexpr uint8_t kRedLedPin = DS_PIN_LED_RED;
constexpr uint8_t kYellowLedPin = DS_PIN_LED_YELLOW;
constexpr uint8_t kGreenLedPin = DS_PIN_LED_GREEN;

} // namespace config
} // namespace deskstation
