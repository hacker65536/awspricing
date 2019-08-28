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
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	//"strconv"
	//"reflect"
	"syscall"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awspricing",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var (
	pricingsvc *pricing.Client
	cfg        aws.Config
	err        error
	ratecnf    = viper.New()
	home       string
)

var exrate float64
var region string
var regions = map[string]string{
	"ap-northeast-1": "Asia Pacific (Tokyo)",
	"ap-northeast-2": "Asia Pacific (Seoul)",
	"ap-south-1":     "Asia Pacific (Mumbai)",
	"ap-southeast-1": "Asia Pacific (Singapore)",
	"ap-southeast-2": "Asia Pacific (Sydney)",
	"ca-central-1":   "Canada (Central)",
	"eu-central-1":   "EU (Frankfurt)",
	"eu-west-1":      "EU (Ireland)",
	"eu-west-2":      "EU (London)",
	"eu-west-3":      "EU (Paris)",
	"sa-east-1":      "South America (Sao Paulo)",
	"us-east-1":      "US East (N. Virginia)",
	"us-east-2":      "US East (Ohio)",
	"us-west-1":      "US West (N. California)",
	"us-west-2":      "US West (Oregon)",
}

var rateurl = "http://free.currencyconverterapi.com/api/v5/convert?q=USD_JPY&compact=y"

func chkerr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "ErrOutput: \n", e)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	home, err = homedir.Dir()
	chkerr(err)

	cfg, err = external.LoadDefaultAWSConfig()
	cfg.Region = endpoints.UsEast1RegionID
	chkerr(err)
	pricingsvc = pricing.New(cfg)

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.awspricing.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initRate()

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.

		// Search config in home directory with name ".awspricing" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".awspricing")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initRate() {
	// check rate conf
	ratecnf.AddConfigPath(home)
	ratecnf.SetConfigName(".awspricing_rate")

	if err := ratecnf.ReadInConfig(); err != nil {
		// if nofile
		getRate()

	} else {
		// chkexpiration
		chkf()

	}
	if ratecnf.Get("USD_JPY.val") == nil {
		//	fmt.Println("nodata")
		getRate()
	}

	//	exrate, _ = strconv.ParseFloat(ratecnf.Get("USD_JPY.val").(string), 64)
	//	fmt.Println("usd", exrate)
	exrate, _ = ratecnf.Get("USD_JPY.val").(float64)

}

func getRate() {
	resp, err := http.Get(rateurl)
	chkerr(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	chkerr(err)

	//fmt.Printf("%s", body)
	w2f(j2y(body))
	ratecnf.ReadInConfig()
	fmt.Printf("\x1b[32m%s\x1b[0m\n", "updata rate file:"+home+"/"+".awspricing_rate.yaml")

}

func j2y(b []byte) []byte {

	m := make(map[string]interface{})
	yaml.Unmarshal(b, &m)
	//	fmt.Println("y", m)
	o, _ := yaml.Marshal(&m)
	//	fmt.Printf("%s", o)
	return o
}

func w2f(o []byte) {
	err = ioutil.WriteFile(home+"/"+".awspricing_rate.yaml", o, 0644)
	chkerr(err)
	/*
		file, err := os.OpenFile(home+"/"+".awspricing_rate.yaml", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(file, o)
	*/
}

func chkf() {
	var s syscall.Stat_t
	syscall.Stat(home+"/"+".awspricing_rate.yaml", &s)

	c, _ := s.Ctim.Unix()

	t := time.Now()
	if t.Unix()-c > 172800 {
		getRate()
	}
}
