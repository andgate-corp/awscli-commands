package commands

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2コマンドの定義
const (
	DescribeInstances = "describe-instances"
	StartInstances    = "start-instances"
	StopInstances     = "stop-instances"
)

// EC2Command EC2関連コマンド用インターフェース
type EC2Command struct {
	OutStream, ErrStream io.Writer
	Result               CommandResult
}

// GetResult コマンドの結果を取得する
func (c *EC2Command) GetResult() CommandResult {
	return c.Result
}

// Run コマンドを実行する
func (c *EC2Command) Run(argv []string) error {

	c.Result = CommandResult{}

	switch argv[0] {
	case DescribeInstances:
		return c.DescribeInstances(argv[1:])
	case StartInstances:
		return c.StartInstances(argv[1:])
	case StopInstances:
		return c.StopInstances(argv[1:])
	default:
		msg := fmt.Sprintf("[EC2] Command %s is not supported.", argv[0])
		return fmt.Errorf(msg)
	}
}

// TagNamesValue Nameタグ取得用のエイリアス
type TagNamesValue []string

// Set flag.Valueインターフェース実装
func (v *TagNamesValue) Set(s string) error {
	strs := strings.Split(s, ",")
	*v = append(*v, strs...)
	return nil
}

func (v *TagNamesValue) String() string {
	return strings.Join(([]string)(*v), ",")
}

// DescribeInstances インスタンス情報を取得する
func (c *EC2Command) DescribeInstances(argv []string) error {

	var (
		input    = &ec2.DescribeInstancesInput{}
		tagnames []string
		region   string
	)

	flags := flag.NewFlagSet(DescribeInstances, flag.ContinueOnError)
	flags.StringVar(&region, "Region", "", "Region")
	flags.Var((*TagNamesValue)(&tagnames), "Name", "Set comma separates 'tag:Names' (ex: A,B,C)")

	if err := flags.Parse(argv); err != nil {
		return err
	}

	sess := session.Must(session.NewSession())
	svc := ec2.New(sess, aws.NewConfig().WithRegion(region))

	fmt.Println(tagnames)

	if len(tagnames) > 0 {
		input.Filters = []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: aws.StringSlice(tagnames),
			},
		}
	}

	output, err := svc.DescribeInstances(input)

	if err != nil {
		return err
	}

	for _, resv := range output.Reservations {
		for _, inst := range resv.Instances {
			c.Result.Attachments = append(c.Result.Attachments, createAttachment(inst, region))
		}
	}

	if len(c.Result.Attachments) > 0 {
		msg := fmt.Sprintf("Found %d Instances.", len(c.Result.Attachments))
		c.Result.Text = msg
	} else {
		msg := "Instance is not found."
		c.Result = CommandResult{
			Text: msg,
		}
	}

	return nil
}

// createAttachment SlackのAttachment形式のデータ構造を作成する
func createAttachment(instance *ec2.Instance, region string) interface{} {

	att := ButtonActionAttachment{
		Fields: []AttachmentField{
			{Title: "InstanceID", Value: *instance.InstanceId, Short: true},
			{Title: "VpcId", Value: *instance.VpcId, Short: true},
			{Title: "Region", Value: region, Short: true},
			{Title: "tag:Name", Value: func() string {
				for _, t := range instance.Tags {
					if *t.Key == "Name" {
						return *t.Value
					}
				}
				return "<null>"
			}(), Short: true},
			{Title: "State", Value: *instance.State.Name, Short: true},
			{Title: "InstanceType", Value: *instance.InstanceType, Short: true},
		},
	}
	if *instance.State.Name == "stopped" {
		att.Fallback = "Start instance"
		att.CallbackID = "start_instance"
		att.Actions = []ButtonActionItem{
			{
				Name:  "action",
				Type:  "button",
				Text:  "Start Instance",
				Value: fmt.Sprintf("ec2 start-instances -InstanceID %s -Region %s", *instance.InstanceId, region),
				Style: "primary",
			},
		}
	} else {
		att.Fallback = "Stop instance"
		att.CallbackID = "stop_instance"
		att.Actions = []ButtonActionItem{
			{
				Name:  "action",
				Type:  "button",
				Text:  "Stop Instance",
				Value: fmt.Sprintf("ec2 stop-instances -InstanceID %s -Region %s", *instance.InstanceId, region),
				Style: "danger",
			},
			{
				Name:  "action",
				Type:  "button",
				Text:  "Stop Instance(Force)",
				Value: fmt.Sprintf("ec2 stop-instances -Force -InstanceID %s -Region %s", *instance.InstanceId, region),
			},
		}
	}
	return att
}

// StartInstances 指定のインスタンスを実行状態にする
func (c *EC2Command) StartInstances(argv []string) error {

	var (
		region     string
		instanceID string
	)

	flags := flag.NewFlagSet(StartInstances, flag.ContinueOnError)
	flags.StringVar(&region, "Region", "", "Region")
	flags.StringVar(&instanceID, "InstanceID", "", "InstanceID")

	if err := flags.Parse(argv); err != nil {
		return err
	}

	if instanceID == "" {
		return fmt.Errorf("InstanceID is not define.")
	}

	sess := session.Must(session.NewSession())
	svc := ec2.New(sess, aws.NewConfig().WithRegion(region))

	input := &ec2.StartInstancesInput{
		InstanceIds: aws.StringSlice([]string{instanceID}),
	}

	output, err := svc.StartInstances(input)

	if err != nil {
		return err
	}

	inst := output.StartingInstances[0]

	c.Result.Text = fmt.Sprintf(
		"Starting instances %s (Prev: %s -> Current: %s",
		*inst.InstanceId,
		*inst.PreviousState.Name,
		*inst.CurrentState.Name)

	return nil
}

// StopInstances 指定のインスタンスを停止状態にする
func (c *EC2Command) StopInstances(argv []string) error {

	var (
		region     string
		instanceID string
		force      bool
	)

	flags := flag.NewFlagSet(StopInstances, flag.ContinueOnError)
	flags.StringVar(&region, "Region", "", "Region")
	flags.StringVar(&instanceID, "InstanceID", "", "InstanceID")
	flags.BoolVar(&force, "Force", false, "Force stop")

	if err := flags.Parse(argv); err != nil {
		return err
	}

	if instanceID == "" {
		return fmt.Errorf("InstaneID is not define.")
	}

	sess := session.Must(session.NewSession())
	svc := ec2.New(sess, aws.NewConfig().WithRegion(region))

	input := &ec2.StopInstancesInput{
		Force:       aws.Bool(force),
		InstanceIds: aws.StringSlice([]string{instanceID}),
	}

	output, err := svc.StopInstances(input)

	if err != nil {
		return err
	}

	inst := output.StoppingInstances[0]
	c.Result.Text = fmt.Sprintf(
		"Stopping instances %s (Prev: %s -> Current: %s)",
		*inst.InstanceId,
		*inst.PreviousState.Name,
		*inst.CurrentState.Name)

	return nil
}
