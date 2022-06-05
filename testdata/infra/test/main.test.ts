import { Template } from 'aws-cdk-lib/assertions';
import { App } from 'aws-cdk-lib';

import { EKSStack } from '../src/main';

test('EKS Snapshot', () => {
  const app = new App();
  const stack = new EKSStack(app, 'eks-test');

  const template = Template.fromStack(stack);

  expect(template.toJSON()).toMatchSnapshot();
});