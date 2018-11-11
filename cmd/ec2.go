// Copyright © 2018 sg
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	//	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/koron/go-dproxy"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
)

// ec2Cmd represents the prd command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("prd called")
		/*
			if name, err := cmd.PersistentFlags().GetString("region"); err == nil {
				fmt.Println("region:", name)
			}
		*/

		if region == "" {
			region = viper.GetString("region")
		}
		if region == "" {
			region = "us-east-1"
		}

		if len(args) > 0 {
			//			fmt.Println(args[0])
			getProductsEc2(args[0])
		} else {
			fmt.Println("type instance type")
		}
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ec2Cmd.Flags().StringVarP(&region, "region", "r", "", "region")
}

func getProductsEc2(itype ...string) {
	input := &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []pricing.Filter{
			{
				Field: aws.String("ServiceCode"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("AmazonEC2"),
			},
			{
				Field: aws.String("instanceType"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String(itype[0]),
			},
			{
				Field: aws.String("location"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String(regions[region]),
			},
			{
				Field: aws.String("operatingSystem"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("Linux"),
			},
			{
				Field: aws.String("preInstalledSw"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("NA"),
			},
			{
				Field: aws.String("tenancy"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("Shared"),
			},
			{
				Field: aws.String("capacitystatus"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("used"),
			},
		},
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(1),
	}

	req := pricingsvc.GetProductsRequest(input)
	result, err := req.Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case pricing.ErrCodeInternalErrorException:
				fmt.Println(pricing.ErrCodeInternalErrorException, aerr.Error())
			case pricing.ErrCodeInvalidParameterException:
				fmt.Println(pricing.ErrCodeInvalidParameterException, aerr.Error())
			case pricing.ErrCodeNotFoundException:
				fmt.Println(pricing.ErrCodeNotFoundException, aerr.Error())
			case pricing.ErrCodeInvalidNextTokenException:
				fmt.Println(pricing.ErrCodeInvalidNextTokenException, aerr.Error())
			case pricing.ErrCodeExpiredNextTokenException:
				fmt.Println(pricing.ErrCodeExpiredNextTokenException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	//	fmt.Println(result)

	p := result.PriceList[0]
	p2 := p["product"]
	at := dproxy.New(p2).M("attributes")
	cpu, _ := at.M("clockSpeed").String()
	mem, _ := at.M("memory").String()
	net, _ := at.M("networkPerformance").String()
	ecu, _ := at.M("ecu").String()
	vcpu, _ := at.M("vcpu").String()
	processor, _ := at.M("physicalProcessor").String()

	t, _ := dproxy.New(p["terms"]).M("OnDemand").Map()
	price := ""
	unit := ""
	description := ""
	for _, s := range t {
		v, _ := dproxy.New(s).M("priceDimensions").Map()

		for _, s2 := range v {

			v2 := dproxy.New(s2)
			price, _ = v2.M("pricePerUnit").M("USD").String()
			unit, _ = v2.M("unit").String()
			description, _ = v2.M("description").String()
		}

	}

	/*
		p3 := p2.(interface{}).(map[string]interface{})["attributes"]
		cpu := p3.(interface{}).(map[string]interface{})["clockSpeed"]
		mem := p3.(interface{}).(map[string]interface{})["memory"]
		net := p3.(interface{}).(map[string]interface{})["networkPerformance"]
		ecu := p3.(interface{}).(map[string]interface{})["ecu"]
		vcpu := p3.(interface{}).(map[string]interface{})["vcpu"]
		processor := p3.(interface{}).(map[string]interface{})["physicalProcessor"]
	*/

	usd, _ := strconv.ParseFloat(price, 64)

	en := exrate / (float64(1) / float64(usd))
	fmt.Printf("CPU: %s MEM: %s NETWORK: %s ecu: %s vcpu: %s processor: %s \n", cpu, mem, net, ecu, vcpu, processor)
	fmt.Printf("PRICE: OnDemand %s USD (%f JP) / %s \n", price, en, unit)
	fmt.Println("DESCRIPTION:", description)
	/*

		for _, s := range p2 {
			fmt.Println("---")
			r, _ := json.Marshal(s)
			fmt.Println(string(r))
			fmt.Println("---")
		}
	*/
}
