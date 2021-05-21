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

### deb/rpm/apk:

Download the .deb, .rpm or .apk from the [releases page](https://github.com/hupe1980/ec2connect/releases) and install them with the appropriate tools.

### go install:
```bash
go install github.com/hupe1980/ec2connect
```
### manually:
Download the pre-compiled binaries from the [releases page](https://github.com/hupe1980/ec2connect/releases) and copy to the desired location.

## How to use
```
Usage:
  ec2connect [command]

Available Commands:
  fwd         Port forwarding
  help        Help about any command
  run         Run commands
  scp         Tunnel scp
  session     Start a session
  ssh         Tunnel ssh

Flags:
  -h, --help             help for ec2connect
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
  -v, --version          version for ec2connect

Use "ec2connect [command] --help" for more information about a command.
```

### Start a session
```
Usage:
  ec2connect session [name|ID|IP|DNS|_] [flags]

Flags:
  -h, --help   help for session

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```
### Port forwarding
```
Usage:
  ec2connect fwd [name|ID|IP|DNS|_] [flags]

Flags:
  -h, --help            help for fwd
  -l, --local string    Local port to use (required)
  -r, --remote string   Remote port to forward to (required)

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```

### Run commands
```
Usage:
  ec2connect run [name|ID|IP|DNS|_] [flags]

Flags:
  -c, --cmd string   Command to exceute (required)
  -h, --help         help for run

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```

### Tunnel ssh
```
Usage:
  ec2connect ssh [name|ID|IP|DNS|_] [flags]

Flags:
  -c, --cmd string        (optional)
  -h, --help              help for ssh
  -i, --identity string   Command to exceute (required)
  -p, --port string       SSH port to us (optional) (default "22")
  -l, --user string       SSH user to us (optional) (default "ec2-user")

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```

### Tunnel scp
```
Usage:
  ec2connect scp [name|ID|IP|DNS|_] [flags]

Flags:
  -h, --help              help for scp
  -i, --identity string   (required)
  -p, --port string       SSH port to us (optional) (default "22")
  -s, --source string     Source in the local host (required)
  -t, --target string     Target in the remote host (required)
  -l, --user string       SCP user to us (optional) (default "ec2-user")

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```
## License
[MIT](LICENCE)