import { App } from 'aws-cdk-lib';

import { EKSStack } from './eks-stack';

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