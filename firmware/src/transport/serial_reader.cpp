#include "transport/serial_reader.hpp"

namespace deskstation {
namespace transport {

SerialReader::SerialReader(Stream& stream, char* buffer, size_t bufferSize)
    : stream_(stream), buffer_(buffer), bufferSize_(bufferSize) {
  reset();
}

bool SerialReader::pollLine(const char*& outLine) {
  outLine = nullptr;

  while (stream_.available() > 0) {
    const int nextByte = stream_.read();
    if (nextByte < 0) {
      break;
    }

    const char c = static_cast<char>(nextByte);
    if (c == '\r') {
      continue;
    }

    if (c == '\n') {
      if (length_ == 0) {
        continue;
      }

      buffer_[length_] = '\0';
      outLine = buffer_;
      length_ = 0;
      return true;
    }

    if (length_ + 1 >= bufferSize_) {
      ++droppedBytes_;
      continue;
    }

    buffer_[length_] = c;
    ++length_;
  }

  return false;
}

size_t SerialReader::droppedBytes() const { return droppedBytes_; }

void SerialReader::reset() {
  length_ = 0;
  droppedBytes_ = 0;
  if (buffer_ != nullptr && bufferSize_ > 0) {
    buffer_[0] = '\0';
  }
}

} // namespace transport
} // namespace deskstation
