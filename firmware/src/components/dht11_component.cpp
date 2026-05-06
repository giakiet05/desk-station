#include "components/dht11_component.hpp"

#include <DHT.h>
#include <math.h>

#include "core/outbox.hpp"
#include "protocol/message_schema.hpp"
#include "protocol/topics.hpp"

namespace deskstation {
namespace components {

DHT11Component::DHT11Component(uint8_t dataPin, uint32_t publishIntervalMs, core::Outbox& outbox)
    : dataPin_(dataPin), publishIntervalMs_(publishIntervalMs), outbox_(outbox) {}

const char* DHT11Component::name() const { return "dht11"; }

void DHT11Component::begin() {
  if (dht_ == nullptr) {
    dht_ = new DHT(dataPin_, DHT11);
  }
  dht_->begin();
}

void DHT11Component::tick(uint32_t nowMs) {
  if (dht_ == nullptr) {
    return;
  }

  if (lastPublishMs_ != 0 && (nowMs - lastPublishMs_) < publishIntervalMs_) {
    return;
  }
  lastPublishMs_ = nowMs;

  const float humidity = dht_->readHumidity();
  const float temperatureC = dht_->readTemperature();
  const bool valid = !(isnan(temperatureC) || isnan(humidity));

  protocol::Dht11EventPayload payload{};
  payload.temperature = valid ? temperatureC : 0.0F;
  payload.humidity = valid ? humidity : 0.0F;
  payload.valid = valid;

  protocol::MessageEnvelope envelope{
      protocol::topics::kProtocolVersion,
      protocol::topics::kEvtDht11,
      protocol::MessageType::Event,
      ++sequence_,
      nowMs,
  };

  char line[protocol::kMaxJsonLineLength] = {};
  if (protocol::encodeDht11EventJson(envelope, payload, line, sizeof(line))) {
    outbox_.push(line);
  }
}

} // namespace components
} // namespace deskstation
