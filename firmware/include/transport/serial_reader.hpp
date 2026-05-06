#pragma once

#include <Arduino.h>
#include <cstddef>

namespace deskstation {
namespace transport {

class SerialReader {
public:
  SerialReader(Stream& stream, char* buffer, size_t bufferSize);

  bool pollLine(const char*& outLine);
  size_t droppedBytes() const;
  void reset();

private:
  Stream& stream_;
  char* buffer_;
  size_t bufferSize_;
  size_t length_ = 0;
  size_t droppedBytes_ = 0;
};

} // namespace transport
} // namespace deskstation
