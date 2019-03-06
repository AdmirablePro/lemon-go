package main

import (
	"golang.org/x/text/language"
)

type languageBundle struct {
	LemonStarting       string
	MetricsEnabled      string
	MetricsInLog        string
	GlobalReportEnabled string
}

var (
	currentLangBundle *languageBundle
)

func init() {
	// i18n
	languageSet := make(map[string]languageBundle)

	bundleEN := languageBundle{
		LemonStarting:       "Lemon (Go %s) is starting...",
		MetricsEnabled:      "Metrics is enabled.",
		GlobalReportEnabled: "Global report is enabled.",
		MetricsInLog:        "Metrics in last %ds: "}
	languageSet[language.English.String()] = bundleEN

	bundleZH := languageBundle{
		LemonStarting:       "Lemon (Go %s) 正在启动...",
		MetricsEnabled:      "指标统计已启用",
		GlobalReportEnabled: "全局报告已启用",
		MetricsInLog:        "过去%d秒中的指标统计: "}
	languageSet[language.Chinese.String()] = bundleZH

	clb := languageSet[lang]
	currentLangBundle = &clb
}
