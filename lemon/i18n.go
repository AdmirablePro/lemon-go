package main

import (
	"golang.org/x/text/language"
)

type languageBundle struct {
	LemonStarting           string
	MetricsEnabled          string
	MetricsInLog            string
	GlobalReportEnabled     string
	FetchingTaskError       string
	FetchingTaskNon200      string
	FetchingTaskDecodeError string
}

var (
	currentLangBundle *languageBundle
)

func init() {
	// i18n
	languageSet := make(map[string]languageBundle)

	bundleEN := languageBundle{
		LemonStarting:           "Lemon (Go %s) is starting...",
		MetricsEnabled:          "Metrics is enabled.",
		GlobalReportEnabled:     "Global report is enabled.",
		MetricsInLog:            "Metrics in last %ds: ",
		FetchingTaskError:       "Error when fetching task: %s",
		FetchingTaskNon200:      "Non-200 status code when fetching task: [%s] %s",
		FetchingTaskDecodeError: "Decode error when fetching task: %s"}
	languageSet[language.English.String()] = bundleEN

	bundleZH := languageBundle{
		LemonStarting:           "Lemon (Go %s) 正在启动...",
		MetricsEnabled:          "指标统计已启用",
		GlobalReportEnabled:     "全局报告已启用",
		MetricsInLog:            "过去%d秒中的指标统计: ",
		FetchingTaskError:       "获取任务时发生错误: %s",
		FetchingTaskNon200:      "获取任务时状态码异常: [%s] %s",
		FetchingTaskDecodeError: "获取任务时解码错误: %s"}
	languageSet[language.Chinese.String()] = bundleZH

	clb := languageSet[lang]
	currentLangBundle = &clb
}
