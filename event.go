package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"fmt"
)

const (
	EVENT_DETAIL_TYPE = "SSH connection event detected"
	EVENT_SOURCE = "com.unbounce.tattle"
	MAX_RETRIES = 3
	DETAIL_TEMPLATE = "{ \"username\": \"%s\", \"remote_host\": \"%s\", \"session\": \"%s\" }"  // yes, hardcoded JSON.  it's easier
)

type EventEmitter func(cfg aws.Config, params cloudwatchevents.PutEventsInput)

type Emitter struct {
	cfg		aws.Config
	doc		ec2metadata.EC2InstanceIdentityDocument
	pam		PAMEvent
	session_type	string
	emitEvent	EventEmitter
}

func NewDefaultEmitter() Emitter {
	return NewEmitter(CloudWatchEventEmitter)
}

func NewEmitter(ee EventEmitter) Emitter {
	return Emitter{emitEvent: ee}
}

func CloudWatchEventEmitter(cfg aws.Config, params cloudwatchevents.PutEventsInput) {
	retryCount := 0

	svc := cloudwatchevents.New(cfg)
	req := svc.PutEventsRequest(&params)
	for retryCount <= MAX_RETRIES {
		_, err := req.Send()
		if err != nil {
			retryCount++;
		} else {
			break
		}
	}
	if retryCount >= MAX_RETRIES {
		panic("Error sending event")
	}
}

func EmitOpenEvent(cfg aws.Config, doc ec2metadata.EC2InstanceIdentityDocument, pam PAMEvent) {
	e := NewDefaultEmitter()
	e.cfg = cfg
	e.doc = doc
	e.pam = pam
	e.session_type = "opened"
	emit(e)
}

func EmitCloseEvent(cfg aws.Config, doc ec2metadata.EC2InstanceIdentityDocument, pam PAMEvent) {
	e := NewDefaultEmitter()
	e.cfg = cfg
	e.doc = doc
	e.pam = pam
	e.session_type = "closed"
	emit(e)
}

func buildEventParams(doc ec2metadata.EC2InstanceIdentityDocument, sessionType string, pam PAMEvent) cloudwatchevents.PutEventsInput {
	var msg string
	var detail string
	instance := NewEC2Instance(doc)
	source := EVENT_SOURCE
	msg = EVENT_DETAIL_TYPE
	detail = fmt.Sprintf(DETAIL_TEMPLATE, pam.Username, pam.RemoteHost, sessionType)

	params := cloudwatchevents.PutEventsInput{
		Entries: []cloudwatchevents.PutEventsRequestEntry{
			cloudwatchevents.PutEventsRequestEntry{
				Detail: &detail,
				DetailType: &msg,
				Resources: []string{ instance.Arn },
				Source: &source,
			},
		},
	}

	if err := params.Validate(); err != nil {
		panic("CloudWatchEvents request is not valid.  Error: " + err.Error())
	}

	return params
}

func emit(e Emitter) {
	params := buildEventParams(e.doc, e.session_type, e.pam)
	e.cfg.Region = e.doc.Region
	e.emitEvent(e.cfg, params)
}
