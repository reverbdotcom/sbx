package help

const Help = `
Usage: sbx <up|down|name>
`

func Run() (string, error) {
	return Help, nil
}
