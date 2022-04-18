#!/bin/sh
mkdir -p /opt/licenser/{bin,licenses}
cp licenses/* /opt/licenser/licenses
cd src/
go build -o /opt/licenser/bin/licenser
chmod +x /opt/licenser/bin/licenser
