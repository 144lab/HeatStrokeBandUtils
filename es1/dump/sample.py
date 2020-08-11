from nordicsemi.dfu.dfu import Dfu
from nordicsemi.dfu.dfu_transport_serial import DfuTransportSerial

transport = DfuTransportSerial("/dev/tty/usbmodem14201", 115200, True, True, False)

a = Dfu("", dfu_transport=transport)
print(a)
