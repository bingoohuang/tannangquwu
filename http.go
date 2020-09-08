package tannangquwu

import (
	"context"
	"database/sql"
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
	args := ctx.QueryArgs()
	num := args.GetUintOrZero("num")
	if num <= 0 {
		num = 1
	}

	db := g.db

	tx, err := db.Begin()
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	run := sqlrun.NewSQLRun(tx, sqlrun.NewMapPreparer(""))
	name := "步兵"
	numSQL := `insert into seq(name, num) values (?, ?) on duplicate key update num = num + ?`
	result := run.DoUpdate(numSQL, name, num, num)
	if result.Error != nil {
		ctx.Error(result.Error.Error(), 500)
		return
	}
	if result.RowsAffected == 0 {
		ctx.Error("result.RowsAffected == 0", 500)
		return
	}

	query := run.DoQuery("select num from seq where name = ?", name)
	numEnd, _ := strconv.ParseInt(query.Rows.([][]string)[0][0], 10, 64)
	//log.Printf("numEnd: %d", numEnd)
	useSQL := `update card set state = 1 where id = ? and state = 0`
	getSQL := `select num from card where id = ?`
	useArgs := []interface{}{numEnd}
	if num > 1 {
		useSQL = `update card set state = 1  where id > ? and id <= ? and state = 0`
		getSQL = `select num from card where id > ? and id <= ? and state = 0`
		useArgs = []interface{}{numEnd - int64(num), numEnd}
	}

	useResult := run.DoUpdate(useSQL, useArgs...)
	if useResult.Error != nil {
		ctx.Error(useResult.Error.Error(), 500)
		return
	}

	if useResult.RowsAffected != int64(num) {
		ctx.Error("Failed", 500)
		return
	}

	getResult := run.DoQuery(getSQL, useArgs...)
	tx.Commit()

	if getResult.RowsCount > 0 {
		//log.Printf("card nums: %s", Join(getResult.Rows.([][]string), ","))
		return
	}

	//log.Printf("no cards found")
	ctx.Error("no cards found", 500)
}

func Join(ss [][]string, sep string) string {
	result := ""

	for _, row := range ss {
		if result != "" {
			result += sep
		}

		result += row[0]
	}

	return result
}
