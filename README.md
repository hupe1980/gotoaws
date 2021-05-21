# ec2connect
ec2connect is an interactive CLI tool that you can use to connect to your EC2 instances using the [AWS Systems Manager Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html). It provides secure and auditable instance management without the need to open inbound ports, maintain bastion hosts, or manage SSH keys. You can connect to your instances by name, ID, DNS, IP or or select an instance from a list.

## Prerequisites
- [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) must be installed on your client
- SSM Agent version 2.3.672.0 or later must be installed on the instances you want to connect to through sessions
- An instance profile with proper IAM permissions (e.g AmazonSSMManagedInstanceCore)
- A connection to the AWS System Manager Servive via NAT or better via [VPC Endpoint](https://docs.aws.amazon.com/vpc/latest/privatelink/vpc-endpoints.html) to further reduce the attack surface

## Installing
You can install the pre-compiled binary in several different ways

### homebrew tap:
```bash
brew tap hupe1980/ec2connect
brew install ec2connect
```

### snapcraft:
```bash
sudo snap install --classic ec2connect
```

### scoop:
```bash
scoop bucket add ec2connect https://github.com/hupe1980/ec2connect-bucket.git
scoop install ec2connect
```

## deb/rpm/apk:

Download the .deb, .rpm or .apk from the [releases page](https://github.com/hupe1980/ec2connect/releases) and install them with the appropriate tools.

### go install:
```bash
go install github.com/hupe1980/ec2connect
```
### manually:
Download the pre-compiled binaries from the [releases page](https://github.com/hupe1980/ec2connect/releases) and copy to the desired location.

## License
[MIT](LICENCE)