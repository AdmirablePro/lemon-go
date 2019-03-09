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
	FetchTaskCount          string
	ConsumingHTTPDoError    string
	SubmitResultNon200      string
	SubmitResultError       string
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
		FetchingTaskNon200:      "Non-200 status code when fetching task: [%d] %s",
		FetchingTaskDecodeError: "Decode error when fetching task: %s",
		FetchTaskCount:          "Received %d tasks from server",
		ConsumingHTTPDoError:    "Error when consuming task: %s",
		SubmitResultNon200:      "Non-200 status code when submitting task: [%d] %s",
		SubmitResultError:       "Error when posting result to server: %s"}
	languageSet[language.English.String()] = bundleEN

	bundleZH := languageBundle{
		LemonStarting:           "Lemon (Go %s) 正在启动...",
		MetricsEnabled:          "指标统计已启用",
		GlobalReportEnabled:     "全局报告已启用",
		MetricsInLog:            "过去%d秒中的指标统计: ",
		FetchingTaskError:       "获取任务时发生错误: %s",
		FetchingTaskNon200:      "获取任务时状态码异常: [%d] %s",
		FetchingTaskDecodeError: "获取任务时解码错误: %s",
		FetchTaskCount:          "从服务器获取了%d条任务",
		ConsumingHTTPDoError:    "执行任务时异常: %s",
		SubmitResultNon200:      "向服务器提交结果时状态码异常: [%d] %s",
		SubmitResultError:       "向服务器提交结果时异常: %s"}
	languageSet[language.Chinese.String()] = bundleZH

	clb := languageSet[lang]
	currentLangBundle = &clb
}
