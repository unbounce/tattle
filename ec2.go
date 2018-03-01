package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
)

const (
	EC2_INSTANCE_ARN_TEMPLATE = "arn:aws:ec2:%s:%s:instance/%s"
)

type EC2Instance struct {
	Id	string
	Arn	string
}

func NewEC2Instance(doc ec2metadata.EC2InstanceIdentityDocument) EC2Instance {
	return EC2Instance{
		Id: doc.InstanceID,
		Arn: fmt.Sprintf(EC2_INSTANCE_ARN_TEMPLATE, doc.Region, doc.AccountID, doc.InstanceID),
	}
}

func getIdentityDoc(cfg aws.Config) ec2metadata.EC2InstanceIdentityDocument {
	md_svc := ec2metadata.New(cfg)
	doc, err := md_svc.GetInstanceIdentityDocument()
	if err != nil {
		panic("Could not retrieve instance identity document. Error: " + err.Error())
	}

	return doc
}
