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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
)

// svcsCmd represents the svcs command
var svcsCmd = &cobra.Command{
	Use:   "svcs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("svcs called")

		if len(args) > 0 {
			//			fmt.Println(args[0])

			describeServices(args[0])
			return
		}
		describeServices()
	},
}

func init() {
	rootCmd.AddCommand(svcsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// svcsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// svcsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func describeServices(svc ...string) {
	input := &pricing.DescribeServicesInput{
		FormatVersion: aws.String("aws_v1"),
		//      MaxResults:    aws.Int64(1),
		//ServiceCode:   aws.String("AmazonEC2"),
	}
	if len(svc) > 0 {

		input.ServiceCode = aws.String(svc[0])
		input.MaxResults = aws.Int64(1)

	}

	req := pricingsvc.DescribeServicesRequest(input)
	resp, err := req.Send()

	chkerr(err)

	fmt.Println(resp)
	for _, s := range resp.Services {

		fmt.Printf("%s %s\n", aws.StringValue(s.ServiceCode), s.AttributeNames)
	}

}
