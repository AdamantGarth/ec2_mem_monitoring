#!/bin/sh
set -xeu

if test main.go -nt ec2_ram_monitoring; then
	CGO_ENABLED=0 go build -ldflags="-w -s"
fi

if test ec2_ram_monitoring -nt ec2_ram_monitoring.tar.gz || test ec2_ram_monitoring.service -nt ec2_ram_monitoring.tar.gz; then
	tar -czf ec2_ram_monitoring.tar.gz --transform 's|^ec2_ram_monitoring\.service$|lib/systemd/system/&|;s|^ec2_ram_monitoring$|bin/&|' ec2_ram_monitoring{,.service}
fi
