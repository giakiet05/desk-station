#pragma once

namespace deskstation {
namespace protocol {
namespace topics {

constexpr char kProtocolVersion[] = "1.0";

constexpr char kEvtIr[] = "evt.ir";
constexpr char kEvtDht11[] = "evt.dht11";
constexpr char kEvtButton[] = "evt.button";
constexpr char kEvtDevice[] = "evt.device";

constexpr char kCmdLed[] = "cmd.led";
constexpr char kCmdPing[] = "cmd.ping";

constexpr char kAckCommand[] = "ack.command";

} // namespace topics
} // namespace protocol
} // namespace deskstation
