package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jackc/pgtype"

	_ "github.com/jackc/pgx/v4/stdlib" // register driver
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	tmplDir := os.Getenv("TEMPLATE_DIR")
	if tmplDir == "" {
		tmplDir = "templates"
	}
	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if err := run(debug, dbURL, ":"+port, tmplDir); err != nil {
		log.Fatal(err)
	}
}

// run runs the application, connecting to the database at dbURL and listening for HTTP at the provided address.
func run(debug bool, dbURL, addr, tmplDir string) error {
	// construct an app ctx that we can cancel
	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	ctx = setDebug(ctx, debug)

	// cancel it if interrupted
	cancelOnSignal(ctx, cncl, os.Interrupt)

	db, err := connectDB(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to db: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error while closing DB: %v", err)
		}
	}()

	tmpls, err := loadTmplMgr(tmplDir)
	if err != nil {
		return fmt.Errorf("loadTmplMgr %v: %w", tmplDir, err)
	}

	tzes, err := loadTimeZones(ctx, db)
	if err != nil {
		return fmt.Errorf("loadTimeZone: %v", err)
	}

	app := route(&handler{tmpls, &txMgr{db}, newStringGener(), tzes})

	return listenAndServe(ctx, addr, app)
}

func loadTimeZones(ctx context.Context, db *sql.DB) (s []string, err error) {
	var arr pgtype.TextArray
	if err = db.QueryRowContext(
		ctx,
		`SELECT ARRAY_AGG(name ORDER BY name) FROM pg_timezone_names`,
	).Scan(&arr); err != nil {
		return
	}
	err = arr.AssignTo(&s)
	return
}

// cancelOnSignal cancels the provided context when any of the signals is received
func cancelOnSignal(ctx context.Context, cncl context.CancelFunc, sigs ...os.Signal) {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, sigs...)

	// cancel the context if either a os.Interrupt arrives
	go func() {
		defer cncl()
		select {
		case <-sigC:
		case <-ctx.Done():
		}
	}()
}

func connectDB(ctx context.Context, url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := func() error {
		ctx, cncl := context.WithTimeout(ctx, 10*time.Second)
		defer cncl()
		return db.PingContext(ctx)
	}(); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return db, nil
}

func listenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	srvr := http.Server{
		Addr:        addr,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Handler:     handler,
	}

	errs := make(chan error, 1)
	go func() {
		errs <- srvr.ListenAndServe()
	}()

	select {
	case err := <-errs: // server shutdown on its own
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("srvr.ListenAndServe: %w", err)

	case <-ctx.Done():
		ctx, cncl := context.WithTimeout(context.Background(), time.Second*5)
		defer cncl()

		log.Println("Beginning shutdown.")

		cancelOnSignal(ctx, cncl, os.Interrupt) // if you interrupt again, we'll shutdown immediately

		if err := srvr.Shutdown(ctx); err != nil {
			return fmt.Errorf("srvr.Shutdown: %w", err)
		}
		return nil
	}
}
