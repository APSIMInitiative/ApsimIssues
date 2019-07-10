package main

import "time"

// options provides a class to store command line arguments.
type options struct {
	Username string `short:"u" default:"hol430" long:"username" description:"github username"`
	Date     string `short:"s" long:"since" default:"1/1/1970" description:"Only show data after this date"`
	Quiet    bool   `short:"q" long:"quiet" description:"Suppress progress reporting"`
	UseCache bool   `short:"c" long:"use-cache" description:"Use cache - do not fetch live data"`
	DryRun   bool   `short:"d" long:"dry-run" description:"Update cache with live data and immediately exit"`
}

// sinceDate returns the 'since' option passed by the user. Defaults to
// 1/1/1970. Panics if provided option is not a valid Time.
func (o options) Since() time.Time {
	t, err := time.Parse("2/1/2006", o.Date)
	if err != nil {
		panic(err)
	}
	return t
}
