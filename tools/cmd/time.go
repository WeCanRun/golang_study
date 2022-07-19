package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"study/tools/internal/timer"
	"time"
)

const (
	format = "2006-01-02 15:04:05"
)

var (
	duration      string
	calculateTime string
)

func init() {
	calcCmd.Flags().StringVarP(&duration, "duration", "d", "", "需要计算的时间，有效单位为时间戳或已格式化后的时间")
	calcCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "", `持续时间，有效单位为 "ns" 
																	"us" "ms" "s" "m" "h"`)
	timeCmd.AddCommand(nowTimeCmd)
	timeCmd.AddCommand(calcCmd)

}

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "时间格式处理",
	Long:  "时间格式处理",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var nowTimeCmd = &cobra.Command{
	Use:   "now",
	Short: "获取当前时间",
	Long:  "获取当前时间",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		log.Printf("当前时间: %s, 对应的时间戳: %d", nowTime.Format(format), nowTime.Unix())
	},
}

var calcCmd = &cobra.Command{
	Use:   "calc",
	Short: "推算时间",
	Long:  "推算时间",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var calculate time.Time
		location, _ := time.LoadLocation(timer.Time_Zone)
		if calculateTime == "" {
			calculate = timer.GetNowTime()
		} else {
			calculate, err = time.ParseInLocation(format, calculateTime, location)
			if err != nil {
				log.Fatalf("time.ParseDuration err: %v", err)
			}
		}

		t, err := timer.GetCalculateTime(calculate, duration)
		if err != nil {
			log.Fatalf("timer.GetCalculateTime err: %v", err)
		}
		log.Printf("输出结果: %s, 对应的时间戳: %d", t.In(location).Format(format), t.Unix())
	},
}
