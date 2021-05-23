# awsconnect
awsconnect is an interactive CLI tool that you can use to connect to your AWS resources (EC2, ECS container) using the [AWS Systems Manager Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html). It provides secure and auditable resource management without the need to open inbound ports, maintain bastion hosts, or manage SSH keys.

## Prerequisites
- [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) must be installed on your client
- SSM Agent version 2.3.672.0 or later must be installed on the instances you want to connect to through sessions
- An instance profile with proper IAM permissions (e.g AmazonSSMManagedInstanceCore)
- A connection to the AWS System Manager Servive via NAT or better via [VPC Endpoint](https://docs.aws.amazon.com/vpc/latest/privatelink/vpc-endpoints.html) to further reduce the attack surface
- [Prerequisites for using ECS Exec](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-exec.html)
## Installing
You can install the pre-compiled binary in several different ways

### homebrew tap:
```bash
brew tap hupe1980/awsconnect
brew install awsconnect
```

### scoop:
```bash
scoop bucket add awsconnect https://github.com/hupe1980/awsconnect-bucket.git
scoop install awsconnect
```

### deb/rpm/apk:

Download the .deb, .rpm or .apk from the [releases page](https://github.com/hupe1980/awsconnect/releases) and install them with the appropriate tools.

### go install:
```bash
go install github.com/hupe1980/awsconnect
```
### manually:
Download the pre-compiled binaries from the [releases page](https://github.com/hupe1980/awsconnect/releases) and copy to the desired location.

## How to use
```
Usage:
  awsconnect [command]

Available Commands:
  completion  Prints shell autocompletion scripts for awsconnect
  ec2         Connect to ec2
  ecs         Connect to ecs
  help        Help about any command

Flags:
  -h, --help             help for awsconnect
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
  -v, --version          version for awsconnect

Use "awsconnect [command] --help" for more information about a command.
```

## EC2
You can connect to your instances by name, ID, DNS, IP or select an instance from a list.
```
Usage:
  awsconnect ec2 [command]

Available Commands:
  fwd         Port forwarding
  run         Run commands
  scp         SCP over Session Manager
  session     Start a session
  ssh         SSH over Session Manager

Flags:
  -h, --help   help for ec2

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)

Use "awsconnect ec2 [command] --help" for more information about a command.
```
#### Start a session
```
Usage:
  awsconnect ec2 session [name|ID|IP|DNS|_] [flags]

Examples:
awsconnect ec2 session myserver

Flags:
  -h, --help   help for session

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```
#### Port forwarding
```
Usage:
  awsconnect ec2 fwd [name|ID|IP|DNS|_] [flags]

Examples:
awsconnect ec2 fwd run myserver -l 8080 -r 8080

Flags:
  -h, --help            help for fwd
  -l, --local string    local port to use (required)
  -r, --remote string   remote port to forward to (required)

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```

#### Run commands
```
Usage:
  awsconnect ec2 run [name|ID|IP|DNS|_] [flags]

Examples:
awsconnect ec2 run myserver -c 'cat /etc/passwd'

Flags:
  -c, --cmd string   command to exceute (required)
  -h, --help         help for run

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```

#### SSH over Session Manager
```
Usage:
  awsconnect ec2 ssh [name|ID|IP|DNS|_] [flags]

Examples:
awsconnect ec2 ssh myserver -i key.pem

Flags:
  -c, --cmd string        command to exceute (optional)
  -h, --help              help for ssh
  -i, --identity string   file from which the identity (private key) for public key authentication is read (required)
  -p, --port string       SSH port to us (optional) (default "22")
  -l, --user string       SSH user to us (optional) (default "ec2-user")

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```

#### SCP over Session Manager
```
Usage:
  awsconnect ec2 scp [name|ID|IP|DNS|_] [flags]

Examples:
awsconnect ec2 scp myserver -i key.pem -s file.txt -t /opt/

Flags:
  -h, --help              help for scp
  -i, --identity string   file from which the identity (private key) for public key authentication is read (required)
  -p, --port string       SSH port to us (optional) (default "22")
  -s, --source string     source in the local host (required)
  -t, --target string     target in the remote host (required)
  -l, --user string       SCP user to us (optional) (default "ec2-user")

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```

## ECS
```
Usage:
  awsconnect ecs [command]

Available Commands:
  exec        Exec into container

Flags:
  -h, --help   help for ecs

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)

Use "awsconnect ecs [command] --help" for more information about a command.
```

### Exec into container
```
Usage:
  awsconnect ecs exec [flags]

Flags:
      --cluster string     arn or name of the cluster (optional)
  -c, --cmd string         command to exceute (optional) (default "/bin/sh")
      --container string   name of the container. A container name only needs to be specified for tasks containing multiple containers. (optional)
  -h, --help               help for exec
      --task string        arn or id of the task (optional)

Global Flags:
      --profile string   AWS profile (optional) (default "default")
      --region string    AWS region (optional)
```
## License
[MIT](LICENCE)