import { Stack, StackProps } from 'aws-cdk-lib';
import { FargateTaskDefinition, ContainerImage, LinuxParameters, FargatePlatformVersion, CfnService } from 'aws-cdk-lib/aws-ecs'
import { Role, ServicePrincipal, ManagedPolicy, PolicyStatement, Effect } from 'aws-cdk-lib/aws-iam'
import { ApplicationLoadBalancedFargateService } from 'aws-cdk-lib/aws-ecs-patterns'
import { Construct } from 'constructs';

export class ECSStack extends Stack {
  constructor(scope: Construct, id: string, props: StackProps = {}) {
    super(scope, id, props);

    const executionRole = new Role(this, 'ExecutionRole', {
        assumedBy: new ServicePrincipal('ecs-tasks.amazonaws.com'),
        managedPolicies: [ManagedPolicy.fromAwsManagedPolicyName('service-role/AmazonECSTaskExecutionRolePolicy')],
    });

    const taskRole = new Role(this, 'TaskRole', {
        assumedBy: new ServicePrincipal('ecs-tasks.amazonaws.com'),
    });
    
    taskRole.addToPrincipalPolicy(new PolicyStatement({
        effect: Effect.ALLOW,
        actions: [
            'ssmmessages:CreateControlChannel',
            'ssmmessages:CreateDataChannel',
            'ssmmessages:OpenControlChannel',
            'ssmmessages:OpenDataChannel',
        ],
        resources: ['*'],
    }));

    const taskDefinition = new FargateTaskDefinition(this, 'TaskDefinition', {
        cpu: 256,
        memoryLimitMiB: 512,
        executionRole,
        taskRole,
    })

    taskDefinition.addContainer('nginx', {
        image: ContainerImage.fromRegistry('public.ecr.aws/nginx/nginx:stable'),
        portMappings: [{ containerPort: 80 }],
        linuxParameters: new LinuxParameters(this, 'LinuxParameters', {
            initProcessEnabled: true,
        }),
    });

    const svc = new ApplicationLoadBalancedFargateService(this, 'NginxService', {
        taskDefinition,
        platformVersion: FargatePlatformVersion.VERSION1_4,
    });

    // enable EnableExecuteCommand for the service
    const cfnService = svc.service.node.findChild('Service') as CfnService;
    cfnService.addPropertyOverride('EnableExecuteCommand', true);

  }
}