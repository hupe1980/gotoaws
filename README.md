# awsconnect
awsconnect is an interactive CLI tool that you can use to connect to your AWS resources (EC2, ECS container) using the [AWS Systems Manager Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html). It provides secure and auditable resource management without the need to open inbound ports, maintain bastion hosts, or manage SSH keys.

![summry](summary.png)

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
  -h, --help               help for awsconnect
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)
  -v, --version            version for awsconnect

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
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)

Use "awsconnect ec2 [command] --help" for more information about a command.
```
#### Start a session
```
Usage:
  awsconnect ec2 session [flags]

Examples:
awsconnect ec2 session -t myserver

Flags:
  -h, --help            help for session
  -t, --target string   name|ID|IP|DNS of the instance (optional)

Global Flags:
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)
```
#### Port forwarding
```
Usage:
  awsconnect ec2 fwd [flags]

Examples:
awsconnect fwd run -t myserver -l 8080 -r 8080

Flags:
  -h, --help            help for fwd
  -l, --local string    local port to use (required)
  -r, --remote string   remote port to forward to (required)
  -t, --target string   name|ID|IP|DNS of the instance (optional)

Global Flags:
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)
```

#### Run commands
```
Usage:
  awsconnect ec2 run [flags]

Examples:
awsconnect ec2 run -t myserver -c 'cat /etc/passwd'

Flags:
  -c, --cmd string      command to exceute (required)
  -h, --help            help for run
  -t, --target string   name|ID|IP|DNS of the instance (optional)

Global Flags:
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)
```

#### SSH over Session Manager
```
Usage:
  awsconnect ec2 ssh [command] [flags]

Examples:
awsconnect ssh -t myserver -i key.pem

Flags:
  -h, --help              help for ssh
  -i, --identity string   file from which the identity (private key) for public key authentication is read (required)
  -L, --lforward string   local port forwarding (optional)
  -p, --port string       SSH port to us (optional) (default "22")
  -t, --target string     name|ID|IP|DNS of the instance (optional)
  -l, --user string       SSH user to us (optional) (default "ec2-user")

Global Flags:
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)
```

#### SCP over Session Manager
```
Usage:
  awsconnect ec2 scp [source(s)] [target] [flags]

Examples:
awsconnect ec2 scp file.txt /opt/ -t myserver -i key.pem

Flags:
  -h, --help              help for scp
  -i, --identity string   file from which the identity (private key) for public key authentication is read (required)
  -p, --port string       SSH port to us (optional) (default "22")
  -R, --recv              receive files from target (optional)
  -t, --target string     name|ID|IP|DNS of the instance (optional)
  -l, --user string       SCP user to us (optional) (default "ec2-user")

Global Flags:
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)
```

## ECS
You can directly interact with containers without needing to first interact with the host container operating system, open inbound ports, or manage SSH keys.
 
```
Usage:
  awsconnect ecs [command]

Available Commands:
  exec        Exec into container

Examples:
awsconnect ecs exec --cluster demo-cluster

Flags:
  -h, --help   help for ecs

Global Flags:
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)

Use "awsconnect ecs [command] --help" for more information about a command.
```

### Exec into container
```
Usage:
  awsconnect ecs exec [flags]

Flags:
      --cluster string     arn or name of the cluster (optional) (default "default")
  -c, --cmd string         command to exceute (optional) (default "/bin/sh")
      --container string   name of the container. A container name only needs to be specified for tasks containing multiple containers. (optional)
  -h, --help               help for exec
      --task string        arn or id of the task (optional)

Global Flags:
      --profile string     AWS profile (optional) (default "default")
      --region string      AWS region (optional)
      --timeout duration   timeout for network requests (default 15s)
```
## License
[MIT](LICENCE)