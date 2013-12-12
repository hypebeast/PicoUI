package main

var cmdPublish = &Command{
	UsageLine: "publish [host]",
	Short:     "publish a PicoUi application",
	Long: `
Publish the PicoUi application.

For example:

	picoui publish raspberry.local
`,
	Run: publishApp,
}

func publishApp(args []string) {

}
