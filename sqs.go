package main

import (
    "fmt"
	"log"
//	"time"
	"encoding/json"

	// amazone stuff
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
)


func subToSqs() {


	sqsconfig := &aws.Config{
		Region:           config.AwsSqs.Region,
        Credentials:      credentials.NewStaticCredentials(*config.AwsSqs.AccessId, *config.AwsSqs.AccessKey, ""),
    }

	// Do connect and session code here
	sess, err := session.NewSession(sqsconfig)
	if err != nil {
		fmt.Println("SQS:  failed to create sqs session,", err)
		log.Println("SQS:  failed to create sqs session,", err)
		return
	} else {
		fmt.Println("SQS: Session established to ", *config.AwsSqs.Region, " / ", *config.AwsSqs.Url)
		log.Println("SQS: Session established to ", *config.AwsSqs.Region, " / ", *config.AwsSqs.Url)
	}
	

	svc := sqs.New(sess)

    for {
		// Do long poll here
		params := &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(*config.AwsSqs.Url),
			AttributeNames: []*string{
				aws.String("ApproximateNumberOfMessages"), // Required
				aws.String("ApproximateNumberOfMessagesNotVisible"), // Required
				aws.String("DelaySeconds"),
				aws.String("CreatedTimestamp"),
				aws.String("ReceiveMessageWaitTimeSeconds"),
			},
			MaxNumberOfMessages: aws.Int64(*config.AwsSqs.Chunksize),
			MessageAttributeNames: []*string{
				aws.String("All"), // Required
			},
			VisibilityTimeout: aws.Int64(10),
			WaitTimeSeconds:   aws.Int64(*config.AwsSqs.Waitseconds),
		}
		resp, err := svc.ReceiveMessage(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println("SQS Error: ", err.Error())
			log.Println("SQS Error: ", err.Error())
			fmt.Println("SQS Response: ", resp)
			log.Println("SQS Response: ", resp)
		} else {

			var msg Capmsg
			for i, v := range resp.Messages  {
				strvalue := *v.Body
				if jerr := json.Unmarshal([]byte(strvalue), &msg); jerr != nil {
					fmt.Println("SQS:  Error marshalling JSON: ", jerr)
					log.Println("SQS:  Error marshalling JSON: ", jerr)
				} else {
					fmt.Println("SQS:  Got capmsg: ", msg.Bpf, " index: ", i)
					log.Println("SQS:  Got capmsg: ", msg.Bpf, " index: ", i)
					if(len(msg.Interface) > 0)  {
						for _, v := range msg.Interface  {
							if _, ok := ifmap[v]; ok  {
								log.Println("SQS:  Interface " + v + " exists in interface map")
								fmt.Println("SQS:  Interface " + v + " exists in interface map")
								go captureToBuffer(msg, v);
							} else {
								log.Println("SQS:  Interface " + v + " does not exist in interface map")
								fmt.Println("SQS:  Interface " + v + " does not exist in interface map")
							}
						}
					} else if(len(msg.Alias) > 0)  {
						for _,v := range msg.Alias  {
							if _, ok := almap[v]; ok  {
                                for _, dname := range almap[v]  {
                                    log.Println("SQS:  Alias " + v + " exists in alias map for device " + dname)
                                    fmt.Println("SQS:  Alias " + v + " exists in alias map for device " + dname)
                                    go captureToBuffer(msg, dname);
                                }
							} else {
								log.Println("SQS:  Alias " + v + " does not exist in alias map")
								fmt.Println("SQS:  Alias " + v + " does not exist in alias map")
							}
						}
					}
				}
				params := &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(*config.AwsSqs.Url),
					ReceiptHandle: aws.String(*v.ReceiptHandle),
				}

				dresp, derr := svc.DeleteMessage(params)
				if derr != nil {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println("SQS:  Delete Error: ", derr.Error())
					log.Println("SQS:  Delete Error: ", derr.Error())
				} else {
					fmt.Println("SQS:  Successfully deleted message: ", *v.ReceiptHandle, ", ", dresp.String())
					log.Println("SQS:  Successfully deleted message: ", *v.ReceiptHandle, ", ", dresp.String())
				}
				
			}		
			// sleep for a couple seconds
			// fmt.Println("Sleeping after message loop")
			// time.Sleep(time.Second * 3)
		}
	}
}
