package tannangquwu

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/valyala/fasthttp"
)

// HTTPCmd is the struct representing generate sub-command.
type HTTPCmd struct {
	addr string
}

// nolint:gochecknoinits
func init() {
	c := HTTPCmd{}
	cmd := &cobra.Command{
		Use:   "http",
		Short: "http service to consume cards",
		Long:  `探囊取物`,
		Run:   c.run,
	}

	rootCmd.AddCommand(cmd)
	c.initFlags(cmd.Flags())
}

func (g *HTTPCmd) initFlags(f *pflag.FlagSet) {
	f.StringVar(&g.addr, "addr", ":8000", "address to listen")
}

func (g *HTTPCmd) run(cmd *cobra.Command, args []string) {
	log.Printf("探囊取物 %+v", g)

	db := setup(0)
	defer db.Close()

	// pass plain function to fasthttp
	if err := fasthttp.ListenAndServe(g.addr, func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
	}); err != nil {
		panic(err)
	}
}
