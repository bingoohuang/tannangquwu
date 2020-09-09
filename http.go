package tannangquwu

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bingoohuang/ginx/pkg/sqlrun"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/valyala/fasthttp"
)

// HTTPCmd is the struct representing generate sub-command.
type HTTPCmd struct {
	addr     string
	poolSize int64
	ctx      context.Context
	db       *sql.DB
	pool     chan int64
}

// nolint:gochecknoinits
func init() {
	c := HTTPCmd{}
	cmd := &cobra.Command{
		Use:   "http",
		Short: "http service to consume cards",
		Long:  `探囊取物HTTP服务`,
		Run:   c.run,
	}
	c.ctx = cmd.Context()

	rootCmd.AddCommand(cmd)
	c.initFlags(cmd.Flags())
}

func (g *HTTPCmd) initFlags(f *pflag.FlagSet) {
	f.StringVar(&g.addr, "addr", ":8000", "address to listen")
	f.Int64Var(&g.poolSize, "pool", 1000, "poolSize size")
}

func (g *HTTPCmd) run(cmd *cobra.Command, args []string) {
	log.Printf("探囊取物服务开启 %+v", g)
	g.db = setup(0)
	g.pool = make(chan int64, g.poolSize)

	go g.poolPump()

	// pass plain function to fasthttp
	if err := fasthttp.ListenAndServe(g.addr, g.tannangquwu); err != nil {
		panic(err)
	}
}

func (g *HTTPCmd) tannangquwu(ctx *fasthttp.RequestCtx) {
	select {
	case num := <-g.pool:
		ctx.Write([]byte(fmt.Sprintf("%d", num)))
	case <-time.After(1 * time.Second):
		ctx.Write([]byte("timeout"))
		ctx.SetStatusCode(500)
	}
}

func (g *HTTPCmd) poolPump() {
	run := sqlrun.NewSQLRun(g.db, sqlrun.NewMapPreparer(""))
	numSQL := fmt.Sprintf(`
		update seq set num = num + %d where name = '步兵';
		set @num = (select num from seq where name = '步兵');
		update card set state = 1 where id > @num - %d and id <= @num and state = 0;
		select @num;`, g.poolSize, g.poolSize)

	for {
		result := run.DoQuery(numSQL)
		if result.Error != nil {
			log.Printf(result.Error.Error())
			continue
		}

		numEnd, _ := strconv.ParseInt(result.Rows.([][]string)[0][0], 10, 64)
		for i := numEnd - g.poolSize + 1; i <= numEnd; i++ {
			g.pool <- i
		}
	}
}
