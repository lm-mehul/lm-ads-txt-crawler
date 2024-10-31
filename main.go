package main

import (
	"log"
	"os"
	"runtime"

	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/server"

	_ "net/http/pprof"
	"runtime/pprof"
)

func main() {

	// cpu profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // Ensure file is closed after profiling

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile() // Stop profiling when main ends

	// Add memory profiling
	ff, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer ff.Close()

	// Force a garbage collection to get up-to-date statistics
	runtime.GC()
	if err := pprof.WriteHeapProfile(ff); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	runtime.SetBlockProfileRate(1) // Enables block profiling
	bf, err := os.Create("block.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer bf.Close()

	pprof.Lookup("block").WriteTo(bf, 0)

	gf, err := os.Create("goroutine.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer gf.Close()

	pprof.Lookup("goroutine").WriteTo(gf, 0)

	db, err := models.SetupSQLConn()
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
		db.Close() // Example of closing a database connection
		os.Exit(1) // Exit after cleanup with a non-zero status code
	}
	defer db.Close()

	s := server.NewService(db)
	s.Start()

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

}
