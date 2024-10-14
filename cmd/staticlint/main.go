package main

// // Config — имя файла конфигурации.
// const Config = `config.json`

// // ConfigData описывает структуру файла конфигурации.
// type ConfigData struct {
// 	Staticcheck []string
// }

// func main() {
// 	appfile, err := os.Executable()

// 	if err != nil {
// 		panic(err)
// 	}
// 	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
// 	if err != nil {
// 		panic(err)
// 	}
// 	var cfg ConfigData
// 	if err = json.Unmarshal(data, &cfg); err != nil {
// 		panic(err)
// 	}
// 	mychecks := []*analysis.Analyzer{
// 		checker.ExitCheckAnalyzer,
// 		printf.Analyzer,
// 		shadow.Analyzer,
// 		structtag.Analyzer,
// 		nilness.Analyzer,
// 		errcheck.Analyzer,
// 		unused.Analyzer.Analyzer,
// 	}
// 	checks := make(map[string]bool)
// 	for _, v := range cfg.Staticcheck {
// 		checks[v] = true
// 	}
// 	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
// 	for _, v := range staticcheck.Analyzers {
// 		if checks[v.Analyzer.Name] {
// 			mychecks = append(mychecks, v.Analyzer)
// 		}
// 	}
// 	multichecker.Main(
// 		mychecks...,
// 	)
// }
