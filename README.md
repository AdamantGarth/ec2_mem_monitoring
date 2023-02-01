ec2_ram_monitoring
==================

Sends MemAvailable statistic to CloudWatch using instance's IAM role.

To install run:

```sh
curl -sL https://github.com/AdamantGarth/ec2_ram_monitoring/releases/latest/download/ec2_ram_monitoring.tar.gz | sudo tar -xzC /usr/local
sudo systemctl enable --now ec2_ram_monitoring
```

**Note:** instance must have an IAM account connected, with `cloudwatch:PutMetricData` permissions.