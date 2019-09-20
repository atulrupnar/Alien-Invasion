/*
	Package main represents simulation of alien invasion.
	It defines a functions to read a map file, build map, deploy aliens and
	the simulation.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
	"os"
	"log"
)

//City Represents roads connecting to cities in the format direction=>city eg. south=>Newyork
//Alien is the Id of current alien on the city. if no alien, the value is 0.
//Assumption : Alien id is never 0. Its always > 0
type City struct {
	roads map[string]string
	alien int
}

//Invasion structure contains all the data structure and methods required for simulation
//cityMap : maps cityNames to City structure
//aliens : Maps aliens to city
//totalMoves : stores the current move sequence
type Invasion struct {
	cityMap map[string]*City
	aliens map[int]string
	totalMoves int
}

// Defines Global variables for logging information
var (
	MaxMovesAllowed = 10
	f, _ = os.OpenFile("../logs/invasion_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	logger = log.New(f, "", log.LstdFlags)
)

func init() {
	//set seed for random number generation
	rand.Seed(time.Now().UnixNano());
}	

//Build Map reads input world map file and creates internal map(cityMap)
func (inv *Invasion)BuildMap(data string) {
	logger.Println("Build Map");
	content := strings.Split(data, "\n");
	for i:= range content {
		//inv.parseLine(content[i]);
		line := content[i];
		city := strings.Split(line, " ");
		c := &City{roads : make(map[string]string),
					 alien : 0};
		var srcCity string
		for i := range city {
			if i == 0 {
				srcCity = city[0];
				inv.cityMap[srcCity] = c;
			} else {
				road := strings.Split(city[i], "=");
				direction := road[0];
				destCity := road[1];
				inv.cityMap[srcCity].roads[direction] = destCity;
			}
		}
	}
}

//GetAllCities : returns array of cities in the map(does not include destroyed cities)
func (inv *Invasion) GetAllCities() []string {
	cityNames := make([]string, 0);
	for k, _ := range inv.cityMap {		
		cityNames = append(cityNames, k);
	}
	return cityNames;
}

//Get the opposite direction
func OppositeDirection(direction string) string {
	switch direction {
		case "east" : return "west";
		case "west" : return "east";
		case "south" : return "north";
		case "north" : return "south";
	}
	return "";
}

//Destroy city : Destroy cities and connected roads
func (inv *Invasion) DestroyCity(cityTo string, alien int) {
	//destroy roads first
	logger.Println("call to destroy city", cityTo);
	for dir, city := range inv.cityMap[cityTo].roads {
		oppDir := OppositeDirection(dir);
		delete(inv.cityMap[city].roads, oppDir);
	}

	alien2 := inv.cityMap[cityTo].alien;
	fmt.Printf("Sequence : %d => City %s is destroyed by %d and %d \n",
		inv.totalMoves, cityTo, alien, alien2);
	logger.Printf("City %s is destroyed by %d and %d", cityTo, alien, alien2);
	//destroy city
	delete(inv.cityMap, cityTo);
	//destroy aliens
	delete(inv.aliens, alien);
	delete(inv.aliens, alien2);
}

//Move Alien : Moves aliens iteratively to random connected city from current city
func (inv *Invasion) MoveCityTo(alien int, cityFrom string, cityTo string) {
	logger.Printf("Move alien %d from %s to %s \n", alien, cityFrom, cityTo);
	if cityFrom != "" {
		inv.cityMap[cityFrom].alien = 0
	}
	//check if cityTo exists
	if _, ok := inv.cityMap[cityTo]; ok == false  {
		logger.Println("cityTo does not exists", cityTo)
		return;
	}
	if (inv.cityMap[cityTo].alien == 0) {
		inv.aliens[alien] = cityTo;
		inv.cityMap[cityTo].alien = alien
		logger.Println("Alien Moved successfully !!!");
	} else {
		inv.DestroyCity(cityTo, alien);
	}
}

//Get Random Key : get random key out of keys from map object
func GetRandomKey(m map[string]string) string {
	/*keys := reflect.ValueOf(m).MapKeys()
	r := rand.Intn(len(keys));
	direction := keys[r];
	return direction;*/
	r := rand.Intn(len(m))
	for k := range m {
		if r == 0 {
			return k;
		}
		r--
	}
	return "";
}

//Deploy Aliens : Assign alien to random cities on city map
func (inv *Invasion) DeployAliens(totalAliens int) {
	allCities := inv.GetAllCities();
	totalCities := len(allCities);
	logger.Printf("Deploy %d Aliens", totalAliens)
	logger.Println("Alien => City")
	for i := 1; i <= totalAliens; i++ {
		cityId := rand.Intn(totalCities);
		cityName := allCities[cityId];
		inv.MoveCityTo(i, "", cityName);
		/*inv.aliens[i] = cityName;
		inv.cityMap[cityName].alien = i;
		logger.Println(i, " => ", cityName);*/
	}
}

//isSimulationOver : checks whether simulation should be continues or not
func (inv *Invasion)IsSimulationOver() bool {
		//check for whether all aliens finished
		if len(inv.aliens) == 0 {
			fmt.Println("All aliens are destroyed")
			logger.Println("All aliens are destroyed")
			return true;
		}
		//check whether aliens completed maximum allowed moves
		if inv.totalMoves > MaxMovesAllowed {
			fmt.Println("Aliens reached final state");
			logger.Println("Aliens reached final state");
			return true;
		}
		return false;
}

//Move aliens : Move aliens iteratively to random connected city
func (inv *Invasion) MoveAliens() {
	logger.Println("----> moveAliens : ", inv.totalMoves);
	var roads map[string]string;
	for alien, city := range inv.aliens {
		roads = inv.cityMap[city].roads;
		if (len(roads) == 0) {
			continue;
		}
		dir := GetRandomKey(roads);
		cityTo := inv.cityMap[city].roads[dir];
		inv.MoveCityTo(alien, city, cityTo);
	}
}

//Write file : converts internal map to raw input 
//format and writes to output file
func (inv *Invasion)WriteFile() {
	var output = "";
	for k, v := range inv.cityMap {
		output += k + " ";
		for dir, cityTo := range v.roads {
			output += dir + "=" + cityTo + " ";
		}
		output += "\n";
	}

    //f, err := os.OpenFile("myfile.data", os.O_APPEND|os.O_WRONLY, 0600)
    f, err := os.Create("../examples/output.txt")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if _, err = f.WriteString(output); err != nil {
        panic(err)
    }
}

//simulates the alien invasion.
func main() {
	logger.Println("INIT");
	var data, err = ioutil.ReadFile("../examples/world_map1.txt");
	if (err != nil) {
		fmt.Println("Error reading file");
		logger.Println("Error reading map file", err);
	}
	inv := Invasion{
				aliens : make(map[int]string),
				totalMoves : 0,
			};
	inv.cityMap = make(map[string]*City)
	//BUILD MAP
	inv.BuildMap(string(data));

	//Deploy Aliens
	totalAliens := 4;
	inv.DeployAliens(totalAliens);

	//Move aliens
	for {
		inv.totalMoves++;
		logger.Println("Move Sequence No", inv.totalMoves);
		//check for stop condition
		if (inv.IsSimulationOver()) {
			break;
		}
		inv.MoveAliens();
	}
	inv.WriteFile();
}