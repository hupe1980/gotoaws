import { App, Stack, StackProps } from 'aws-cdk-lib';
import { Cluster, EndpointAccess, KubernetesVersion, NodegroupAmiType } from 'aws-cdk-lib/aws-eks'
import { InstanceType, Vpc, SubnetType } from 'aws-cdk-lib/aws-ec2'
import { Role, AccountRootPrincipal, ManagedPolicy, PolicyStatement, Effect } from 'aws-cdk-lib/aws-iam'
import { Construct } from 'constructs';

export class EKSStack extends Stack {
  constructor(scope: Construct, id: string, props: StackProps = {}) {
    super(scope, id, props);

    const vpc = new Vpc(this, 'Vpc', {
      maxAzs: 3,
      natGateways: 1,
      cidr: '10.0.0.0/22', // 1024 IPs
      subnetConfiguration: [
        {
          name: 'Public',
          subnetType: SubnetType.PUBLIC,
          cidrMask: 26 // 3 x 64 IPs
        },
        {
          name: 'Private',
          subnetType: SubnetType.PRIVATE_WITH_NAT,
          cidrMask: 24, // 3 x 256 IPs
        }
      ] 
    })

    const clusterAdminRole = new Role(this, 'ClusterAdminRole', {
      roleName: 'cluster-admin',
      assumedBy: new AccountRootPrincipal(),
    });

    const cluster = new Cluster(this, 'Cluster', {
      version: KubernetesVersion.V1_21,
      endpointAccess: EndpointAccess.PUBLIC_AND_PRIVATE,
      vpc,
      defaultCapacity: 0,
      clusterName: 'gotoaws',
      mastersRole: clusterAdminRole,
    });

    clusterAdminRole.addToPrincipalPolicy(new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ['eks:DescribeCluster'],
      resources: [cluster.clusterArn]
    }))
    
    const ng = cluster.addNodegroupCapacity('NodeGroup', {
      nodegroupName: 'gotoaws',
      instanceTypes: [new InstanceType('t3.small')],
      maxSize: 3,
      diskSize: 20,
      amiType: NodegroupAmiType.BOTTLEROCKET_X86_64,
    });

    ng.role.addManagedPolicy(ManagedPolicy.fromAwsManagedPolicyName('AmazonSSMManagedInstanceCore'))
  }
}

// for development, use account/region from cdk cli
const devEnv = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: process.env.CDK_DEFAULT_REGION,
};

const app = new App({
  analyticsReporting: false,
});

new EKSStack(app, 'gotoaws-eks-dev', { env: devEnv });

app.synth();