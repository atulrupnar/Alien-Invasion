//Package main is a staring point of Alien Invasion simulation
//It reads command line arguments and call a simulation function.
package main

import (
	"flag"
	"fmt"
	"os"
	"../src"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: Alien Invasion [flags]")
		flag.PrintDefaults();
	}

	mapFile := flag.String("map", "world_map1.txt", "map file path")
	numAliens := flag.Int("aliens", 0, "number of aliens to be deployed")
	//logFile := flag.String("log", "logs/log.txt", "log file path");

	flag.Parse()
	if *numAliens == 0 {
		fmt.Println("Aliens is mandatory field, please check the correct usage");
		flag.Usage();
		os.Exit(1)
	}
	//call to simulation function
	invasion.Run(*numAliens, *mapFile);
}