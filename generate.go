package tannangquwu

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// GenerateCmd is the struct representing generate sub-command.
type GenerateCmd struct {
	NumRecs    int
	BatchSize  int
	LogSeconds int
}

// nolint:gochecknoinits
func init() {
	c := GenerateCmd{LogSeconds: 10}
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generate records to benchmark against",
		Long:  `生成百万大军`,
		Run:   c.run,
	}

	rootCmd.AddCommand(cmd)
	c.initFlags(cmd.Flags())
}

func (g *GenerateCmd) initFlags(f *pflag.FlagSet) {
	// Here you will define your flags and configuration settings.
	f.IntVarP(&g.NumRecs, "records", "r", 100000, "number of records to generate")
	f.IntVarP(&g.BatchSize, "batch", "b", 1000, "batch size to insert at one time")
}

func (g *GenerateCmd) run(cmd *cobra.Command, args []string) {
	log.Printf("Generating records by config %+v", g)

	db := setup(0)
	defer db.Close()

	// Preinitialize i so that we can use it in a goroutine to give proper feedback
	var i int
	// Set up logging mechanism. We use a goroutine here which logs the
	// records already generated every two seconds until "done" is signaled
	// via the channel.
	done := make(chan bool)
	start := time.Now()

	if g.NumRecs > 0 {
		go g.inserts(&i, db, done)
		g.progressLogging(start, &i, done)
	}
}

// nolint:gomnd,gosec
func (g GenerateCmd) inserts(i *int, db *sql.DB, done chan bool) {
	// Prepare values needed so that there aren't any allocations done in the loop
	query := "INSERT INTO card(num, code) VALUES" + strings.Repeat(",(?,?)", g.BatchSize)[1:]

	var execFn func(args ...interface{}) (sql.Result, error)

	ps, _ := db.Prepare(query)
	defer ps.Close()

	execFn = ps.Exec
	lastNum := g.NumRecs % g.BatchSize

	// Start generation of actual records
	log.Println("Starting inserts")

	args := make([]interface{}, 0, g.BatchSize*3)

	for *i = 0; *i < g.NumRecs; *i++ {
		// B77AF07D-4617-47A8-ABCF-DA74C809D881
		args = append(args, uuid.New().String(), uuid.New().String())

		if len(args) == g.BatchSize*2 {
			if _, err := execFn(args...); err != nil {
				log.Fatalf("Inserting values into database failed: %s", err)
			}

			args = args[0:0]
		} else if lastNum > 0 && *i+1 == g.NumRecs {
			query := "INSERT INTO v VALUES" +
				strings.Repeat(",(?,?)", lastNum)[1:]
			if _, err := db.Exec(query, args...); err != nil {
				log.Fatalf("Inserting values into database failed: %s", err)
			}
		}
	}

	// Signal the progress log that we are done
	done <- true
}

// nolint:gomnd
func (g GenerateCmd) progressLogging(start time.Time, i *int, done chan bool) {
	log.Println("Starting progress logging")

	l := len(fmt.Sprintf("%d", g.NumRecs))
	// Precalculate the percentage each record represents
	p := float64(100) / float64(g.NumRecs)

	ticker := time.NewTicker(time.Duration(g.LogSeconds) * time.Second)
	defer ticker.Stop()

out:
	for {
		select {
		// Since this is a time consuming process depending on the number of
		// records	created, we want some feedback every 2 seconds
		case <-ticker.C:
			dur := time.Since(start)
			log.Printf("%*d/%*d (%6.2f%%) written in %s, avg: %s/record, %2.2f records/s",
				l, *i, l, g.NumRecs, p*float64(*i), dur,
				time.Duration(dur.Nanoseconds()/int64(*i)), float64(*i)/dur.Seconds())
		case <-done:
			break out
		}
	}

	dur := time.Since(start)
	log.Printf("%*d/%*d (%6.2f%%) written in %s, avg: %s/record, %2.2f records/s",
		l, g.NumRecs, l, g.NumRecs, p*float64(g.NumRecs), dur,
		time.Duration(dur.Nanoseconds()/int64(g.NumRecs)), float64(g.NumRecs)/dur.Seconds())
}
