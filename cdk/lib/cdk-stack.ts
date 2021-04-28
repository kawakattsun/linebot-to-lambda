import * as cdk from '@aws-cdk/core';
import { CfnFunction } from '@aws-cdk/aws-sam'
import { PolicyStatement } from '@aws-cdk/aws-iam'

export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new CfnFunction(this, 'linebot-to-lambda-function', {
      codeUri: 'deploy/lambda/linebot-to-lambda',
      handler: 'main',
      runtime: 'go1.x',
      environment: {
        variables: {
          'GOOGLE_CALENDAR_ID': '/line-to-lambda/GOOGLE_CALENDAR_ID',
          'LINE_CHANNEL_ACCESS_TOKEN': '/line-to-lambda/LINE_CHANNEL_ACCESS_TOKEN',
          'LINE_CHANNEL_SECRET': '/line-to-lambda/LINE_CHANNEL_SECRET',
        }
      },
      events:{
        GetMethod: {
          type: 'Api',
          properties: {
            path: '/linebot-to-lambda',
            method: 'GET',
          }
        }
      },
      functionName: 'linebot-to-lambda',
      policies: [
        {
          statement: new PolicyStatement({
            actions: [
              'sts:AssumeRole',
              'ssm:GetParameters',
            ],
            resources: [
              '*'
            ]
          })
        }
      ]
    })
  }
}
