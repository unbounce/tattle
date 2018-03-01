package main

import (
	"testing"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
)

func MockEventEmitter(cfg aws.Config, params cloudwatchevents.PutEventsInput) {
	fmt.Println("Event mocked")
	fmt.Println("- Detail: " + *params.Entries[0].Detail)
	fmt.Println("- DetailType: " + *params.Entries[0].DetailType)
	fmt.Println("- Source: " + *params.Entries[0].Source)
	fmt.Println("- Resource: " + params.Entries[0].Resources[0])
}

func buildPAMEvent() PAMEvent {
	p := PAMEvent{
		Username: "johndoe",
		SessionType: "open_session",
		RemoteHost: "home.example.com",
	}
	return p
}

func buildIdentityDocument() ec2metadata.EC2InstanceIdentityDocument {
	doc := ec2metadata.EC2InstanceIdentityDocument{
		AvailabilityZone: "ca-west-1a",
		Region: "ca-west-1",
		InstanceID: "i-123abc",
		AccountID: "012345678901",
	}
	return doc
}

func TestEmitEvent(t *testing.T) {
	e := NewEmitter(MockEventEmitter)

	e.pam = buildPAMEvent()
	e.doc = buildIdentityDocument()
	e.session_type = "mocked"

	emit(e)
}

