PROGRAM_NAME='test'

DEFINE_DEVICE

dvTP = 10001:1:0

DEFINE_EVENT

data_event[dvTP] {
    online: {
        send_string data.device, "'Touch Panel Online'"
    }
}
