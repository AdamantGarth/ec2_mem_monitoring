#!/bin/sh
set -xeu

CGO_ENABLED=0 go build -ldflags="-w -s"
tar -czf ec2_mem_monitoring.tar.gz --transform 's|^ec2_mem_monitoring\.service|lib/systemd/system/&|;s|^ec2_mem_monitoring$|bin/&|' ec2_mem_monitoring{,.service{,.d}}

CGO_ENABLED=0 GOARCH=arm64 go build -ldflags="-w -s" -o ec2_mem_monitoring.arm64
tar -czf ec2_mem_monitoring.arm64.tar.gz --transform 's|^ec2_mem_monitoring\.service|lib/systemd/system/&|;s|^ec2_mem_monitoring.arm64$|bin/ec2_mem_monitoring|' ec2_mem_monitoring{.arm64,.service{,.d}}
