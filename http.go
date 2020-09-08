package tannangquwu

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bingoohuang/ginx/pkg/sqlrun"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

// HTTPCmd is the struct representing generate sub-command.
type HTTPCmd struct {
	addr string
	ctx  context.Context
	db   *sql.DB
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
	c.ctx = cmd.Context()

	rootCmd.AddCommand(cmd)
	c.initFlags(cmd.Flags())
}

func (g *HTTPCmd) initFlags(f *pflag.FlagSet) {
	f.StringVar(&g.addr, "addr", ":8000", "address to listen")
}

func (g *HTTPCmd) run(cmd *cobra.Command, args []string) {
	log.Printf("探囊取物 %+v", g)
	g.db = setup(0)

	// pass plain function to fasthttp
	if err := fasthttp.ListenAndServe(g.addr, g.tannangquwu); err != nil {
		panic(err)
	}
}

func (g *HTTPCmd) tannangquwu(ctx *fasthttp.RequestCtx) {
	run := sqlrun.NewSQLRun(g.db, sqlrun.NewMapPreparer(""))
	numSQL := `
		update seq set num = num + 1 where name = '步兵';
		update card set state = 1 where id = (select num from seq where name = '步兵') and state = 0;
		select num from seq where name = '步兵';
		`
	result := run.DoQuery(numSQL)
	if result.Error != nil {
		ctx.Error(result.Error.Error(), 500)
		return
	}

	numEnd, _ := strconv.ParseInt(result.Rows.([][]string)[0][0], 10, 64)
	ctx.Write([]byte(fmt.Sprintf("%d", numEnd)))
}
