package cli

func Run(args []string) int {
	root := NewRootCommand()
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		return 1
	}

	return 0
}
