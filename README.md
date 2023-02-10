ec2_mem_monitoring
==================

Sends MemAvailable statistic to CloudWatch using instance's IAM role.

Optionally, can monitor available disk space as well. To enable the feature, run the tool with `MONITOR_DISK` environment variable set to the path on the disk that you want to monitor. To monitor root partition, for example:

```sh
MONITOR_DISK=/ ./ec2_mem_monitoring
```

To install run:

```sh
curl -sL https://github.com/AdamantGarth/ec2_mem_monitoring/releases/latest/download/ec2_mem_monitoring.tar.gz | sudo tar -xzC /usr/local --exclude='*.d'
sudo systemctl enable --now ec2_mem_monitoring
```

With root partition disk space monitoring enabled:

```sh
curl -sL https://github.com/AdamantGarth/ec2_mem_monitoring/releases/latest/download/ec2_mem_monitoring.tar.gz | sudo tar -xzC /usr/local
sudo systemctl enable --now ec2_mem_monitoring
```

**Note:** instance must have an IAM account connected, with `cloudwatch:PutMetricData` permissions.
