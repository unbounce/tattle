package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.Version(APP_VERSION)
	kingpin.Parse()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load SDK config, " + err.Error())
	}

	md_svc := ec2metadata.New(cfg)
	if !md_svc.Available() {
		panic("this is not an ec2 instance")
	}

	pam := NewPAMEvent()

	doc := getIdentityDoc(cfg)

	switch pam.SessionType {
	case "open_session":
		EmitOpenEvent(cfg, doc, pam)
	case "close_session":
		EmitCloseEvent(cfg, doc, pam)
	default:
		kingpin.Errorf("Unsupported PAM_TYPE (%s)", pam.SessionType)
	}
}
