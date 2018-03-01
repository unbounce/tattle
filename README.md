# tattle

Tattle is a simple, self-contained binary that will report to Amazon Web
Services (AWS) that a login event has occurred via a Pluggable
Authentication Module (PAM).

The event is sent through CloudWatch Events as it is intended to be an
asynchronous event for security monitoring services.  Consumers can
choose to react to this event by listening with a particular CloudWatch
Event rule pattern.

The producer logic is as follows:
* This binary is installed into the system.
* PAM configuration is extended to call the binary on every sshd event
* User successfully authenticates to server via SSH
* PAM module calls binary
* Binary publishes login event to CloudWatch Events
* User logs out from server
* PAM module calls binary
* Binary publishes logout event to CloudWatch Events

As for consumers, they can listen for either login or logout events (or
both) and process the event with myriad AWS services.  For instance, the
event could be sent to an SNS topic which alerts an OpsGenie group, or
to a Lambda function which publishes only login events to a Slack
channel to warn administrators.

## Event Payload

The event pattern is a specific structure so that consumers can read it
as they wish.

Below is someone completing an SSH connection to a server:

```json
{
  "detail-type": "SSH connection event detected",
  "source": "com.unbounce.tattle",
  "resources": [
    "arn:aws:ec2:us-east-1:01234567890:instance/i-abc123abc123"
  ],
  "detail": {
    "username": "ubuntu",
    "remote_host": "1.2.3.4",
    "session": "opened"
  }
}
```

Below is someone closing an SSH connection to a server:

```json
{
  "detail-type": "SSH connection event detected",
  "source": "com.unbounce.tattle",
  "resources": [
    "arn:aws:ec2:us-east-1:01234567890:instance/i-abc123abc123"
  ],
  "detail": {
    "username": "ubuntu",
    "remote_host": "1.2.3.4",
    "session": "closed"
  }
}
```

Note that these events are emitted whenever the SSH daemon is used for
connecting to a server, which includes `scp` connections.

## Configuration

At this time there is no way to configure the binary.  This is on the
roadmap.  The idea is that the source and detail-type values in the event
payload will be configurable.

### PAM

Add the following line to the end of `/etc/pam.d/sshd`:

```
session    optional    pam_exec.so /usr/local/bin/tattle
```

This has been tested on Ubuntu 14.04 and 16.04 systems.

### EC2

You will need to update the EC2 server's IAM role to allow `events:PutEvent`
actions so that tattle can interact with the CloudWatchEvents service.
There is no Resource scoping on this action.  Here is a sample policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowCloudWatchEventEmit",
            "Effect": "Allow",
            "Action": "events:PutEvents",
            "Resource": "*"
        }
    ]
}
```

## What PAM Gives Us

Upon SSH authentication
```
PAM_SERVICE=sshd
PAM_RHOST=10.0.2.2
PAM_USER=vagrant
PAM_TYPE=open_session
PAM_TTY=ssh
```

Upon SSH logout
```
PAM_SERVICE=sshd
PAM_RHOST=10.0.2.2
PAM_USER=vagrant
PAM_TYPE=close_session
PAM_TTY=ssh
```

## Installation

Download the latest release and put it in the `/usr/local/bin/` directory
on the target machine.  You can use any directory you like, so long as you
change the PAM configuration to reference the correct installation path.

The binary is self-contained, all libraries and dependencies are
statically-linked.

Before configuring PAM, you can test that everything is working by running
the following command on the server to send a test event:

```
PAM_TYPE="open_session" PAM_USER="test" PAM_RHOST="test" tattle
```

Make sure that you have a CloudWatch Event Rule setup or else the event
will not be consumed and processed by anything.  The simplest Event Rule
is to add a target to an SNS topic that sends an email to you.  The
broadest event pattern to use when creating the rule is:

```json
{
  "source": [
    "com.unbounce.tattle"
  ]
}
```

## Failure Conditions

As this application relies on the distributed nature of AWS, there are several ways in which it can fail.

### Tattle run outside of EC2

This code requires the use of EC2's metadata service in order to discover
the machine's identity.  It cannot be used without being on EC2 and will
throw an error.

### Metadata service is unavailable

The Metadata service could be unavailable due to downtime of the API endpoint, or networking issues with the EC2 instance.  In this event, tattle cannot retrieve the identity document to discover information about itself.

### EC2 has no IAM Role

It is difficult to know if the EC2 instance has permission to emit events to CloudWatch Events until the request is made.  In this case, tattle will panic and emit an error back to the logs.

```
panic: Error sending event, EC2RoleRequestError: no EC2 instance role found
caused by: EC2MetadataError: failed to make EC2Metadata request
caused by: <?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
         "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>404 - Not Found</title>
 </head>
 <body>
  <h1>404 - Not Found</h1>
 </body>
</html>


goroutine 1 [running]:
main.main()
        /usr/local/bin/tattle/main.go:72 +0x8db
```

A panic was chosen so that the caller can easily see where the error is
coming from, without custom formatting or interpretation by tattle.  This
may change in the future.

## License

MIT License.  See [LICENSE](LICENSE) for details.

