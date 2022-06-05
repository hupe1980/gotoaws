const { awscdk } = require('projen');

const project = new awscdk.AwsCdkTypeScriptApp({
  cdkVersion: '2.27.0',
  defaultReleaseBranch: 'main',
  name: 'infra',
  //devDeps: ['ts-node@^10'],
  // cdkDependencies: undefined,  /* Which AWS CDK modules (those that start with "@aws-cdk/") this app uses. */
  // deps: [],                    /* Runtime dependencies of this module. */
  // description: undefined,      /* The description is just a string that helps people understand the purpose of the package. */
  // devDeps: [],                 /* Build dependencies for this module. */
  // packageName: undefined,      /* The "name" in package.json. */
  // release: undefined,          /* Add release management to this project. */
  //no-github-workflow
  github: false,

  license: 'MIT',
  copyrightOwner: 'Frank Hübner',
});

project.synth();