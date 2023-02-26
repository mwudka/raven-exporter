# Overview

raven-exporter is a prometheus exporter that exposes the metrics from [Rainforest Automation](https://www.rainforestautomation.com/homeowners/) energy monitors like the [EMU-2](https://www.rainforestautomation.com/rfa-z105-2-emu-2-2/) that speak the [raven protocol](https://rainforestautomation.com/wp-content/uploads/2014/02/raven_xml_api_r127.pdf).

# Installation

raven-exporter is a single go binary with no dependencies. To run, copy the appropriate binary to a computer connected to the energy monitor and run. For example, if running on a raspi something like this will work:

1. Connect the energy monitor to your computer via USB. A new device should appear in `/dev` called `/dev/ttyACM0` or similar.
2. Copy the raven-exporter.service file to /etc/systemd/system/raven-exporter.service
3. Edit the file to set the correct username and `--serial-port`. By default, the metrics will be available on localhost:2112/metrics. To change the host, port, or path add `--http-host`, `http-port`, or `metrics-path` arguments to the service file
4. Reload unit files: `sudo systemctl daemon-reload`
5. Start the service: `sudo service raven-exporter start`
6. Monitor the startup process: `sudo journalctl -f -u raven-exporter`

If everything is working, you should see output like
```
Feb 26 15:42:46 hostname.local raven-exporter[105127]: Demand for 0x0000000000000000: 1149 watts
Feb 26 15:43:01 hostname.local raven-exporter[105127]: Demand for 0x0000000000000000: 1162 watts
Feb 26 15:43:16 hostname.local raven-exporter[105127]: Demand for 0x0000000000000000: 1170 watts
Feb 26 15:43:31 hostname.local raven-exporter[105127]: Demand for 0x0000000000000000: 1149 watts
Feb 26 15:43:46 hostname.local raven-exporter[105127]: Demand for 0x0000000000000000: 1161 watts
Feb 26 15:44:01 hostname.local raven-exporter[105127]: Demand for 0x0000000000000000: 1146 watts
Feb 26 15:44:07 hostname.local raven-exporter[105127]: Total delivered for 0x0000000000000000: 100000 watt-hours
Feb 26 15:44:16 hostname.local raven-exporter[105127]: Demand for 0x0000000000000000: 1159 watts
```
