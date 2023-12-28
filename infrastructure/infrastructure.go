package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfrastructureStackProps struct {
	awscdk.StackProps
}

func NewInfrastructureStack(scope constructs.Construct, id string, props *InfrastructureStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	/* create the s3 bucket for storing or stock data*/
	bucket := awss3.NewBucket(stack, jsii.String("data-bucket"), &awss3.BucketProps{
		BucketName:    jsii.String("awesome-stock-data-bucket"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Versioned:     jsii.Bool(true),
	})

	lambdaRole := awsiam.NewRole(stack, jsii.String("LambdaRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
	})

	bucketArn := *bucket.BucketArn()
	lambdaRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:GetObject", "s3:PutObject"),
		Resources: jsii.Strings(bucketArn + "/*"),
	}))

	lambdaCron := awslambda.NewFunction(stack, jsii.String("Stock-downloader"), &awslambda.FunctionProps{
		FunctionName: jsii.String("StockDownloader"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(300)),
		Code: awslambda.AssetCode_FromDockerBuild(jsii.String("../stock-downloader-lambda"), &awslambda.DockerBuildAssetOptions{
			Platform: jsii.Sprintf("linux/amd64"),
		}),
		Handler: jsii.String("main.handler"),
		Runtime: awslambda.Runtime_PYTHON_3_10(),
		Role:    lambdaRole,
	})

	rule := awsevents.NewRule(stack, jsii.String("Rule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Expression(jsii.String("cron(0 18 ? * MON-FRI *)")),
	})

	rule.AddTarget(awseventstargets.NewLambdaFunction(lambdaCron, nil))

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewInfrastructureStack(app, "InfrastructureStack", &InfrastructureStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
