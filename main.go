package glockChain

import (
	"os"

	"github.com/atoyr/glockchain"
)

func main() {
	cli := core.NewCLI()
	cli.App.Run(os.Args)
}
