/*
	Package main represents simulation of alien invasion.
	It defines a functions to read a map file, build map, deploy aliens and
	the simulation.
*/
package invasion

import (
	"fmt"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
	"os"
	"log"
    "regexp"
    "errors"
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
	f, _ = os.OpenFile("./logs/invasion_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	logger = log.New(f, "", log.LstdFlags)
)

func init() {
	//set seed for random number generation
	rand.Seed(time.Now().UnixNano());
}	

func New() *Invasion {
	return	&Invasion{
		aliens : make(map[int]string),
		cityMap : make(map[string]*City),
		totalMoves : 0,
	};
}

//function to check valid city road entries
//city roads are bidirectional. It checks for when road from city A to B exists, 
//road B to A must exist
func isValidRoads(cityMap map[string]*City) error {
	//loop through citmap
	for cityFrom, val := range cityMap {
		//loop through roads
		for dir, cityTo := range val.roads {
			if _, isCityExists := cityMap[cityTo]; !isCityExists {
				logger.Printf("error : Incomplete Map : City %s is not defined\n", cityTo);
				return errors.New(fmt.Sprintf("error : Incomplete Map : City %s is not defined", cityTo));
			}
			oppositeDir := OppositeDirection(dir);
			//Check for city on opposite direction
			//ex. if CityA south=CityB, Then CityB north=CityA 
			if cityFrom != cityMap[cityTo].roads[oppositeDir] {
				logger.Println("error : Invalid map entry : check roads for following cities => ", cityFrom, cityTo);
				return errors.New(fmt.Sprintf("error : Invalid map entry : check roads for %s and %s", cityFrom, cityTo));
			}
		}
	}
	return nil;
}

//Build Map reads input world map file and creates internal map(cityMap)
func (inv *Invasion)BuildMap(data string) error {
	logger.Println("Build Map");
	//split file content to multiple lines
	content := strings.Split(data, "\n");
	if err := isValidInput(content); err != nil {
		logger.Println("error : Invalid Input", err.Error());
		log.Fatal(err.Error());
	}
	for _, line := range content {
		//split line => City [direction=city direction=city ...]
		lineInfo := strings.Split(line, " ");
		srcCity := strings.TrimSpace(lineInfo[0])

		//Check for duplicate city on map
		if _, isCityExists := inv.cityMap[srcCity]; isCityExists {
			logger.Println("error : Duplicate city found, City can not be redefined", srcCity);
			log.Fatalf("error : Duplicate city found, City can not be redefined : %s", srcCity);
		}

		c := &City{
			roads : make(map[string]string),
			alien : 0,
		};
		inv.cityMap[srcCity] = c;

		//cache to store destination cities for each city(to identify invalid entry)
		cityCache := make(map[string]bool);

		//loop through all the roads of city and stores info in cityMap
		for _, dirInfo  := range lineInfo[1:] {
			dirInfo = strings.TrimSpace(dirInfo)
			//split => Direction=City
			road := strings.Split(dirInfo, "=");
			direction := strings.ToLower(road[0]);

			//Check for duplicate road on same direction.
			//ex. south=Mumbai south=Delhi is invalid
			if _, isRoadExists := inv.cityMap[srcCity].roads[direction]; isRoadExists {
				logger.Println("error : Duplicate road found, road can not be redefined %s", lineInfo[1:]);
				log.Fatalf("error : Duplicate road found, road can not be redefined : %s", lineInfo[1:]);
			}

			destCity := road[1];
			if destCity == srcCity {
				logger.Println("error : A city can not have road to itself", srcCity);
				log.Fatalf("error : A city can not have road to itself. Check for %s city", srcCity);
			}

			//Check for duplicate destination city. 
			//ex. south=Mumbai north=Mumbai is invalid
			if _, isCityToExists := cityCache[destCity]; isCityToExists {
				logger.Println(`error : Duplicate destination city found, two roads
					from same city can not have same destination city`, lineInfo[1:]);
				log.Fatalf(`error : Two roads from same source can not same destination : %s`, lineInfo[1:]);

			}
			cityCache[destCity] = true;
			inv.cityMap[srcCity].roads[direction] = destCity;
		}
	}
	if err := isValidRoads(inv.cityMap); err != nil {
		log.Fatal(err.Error());
	}
}

func isValidInput(input []string) error {
    dir := `(south|north|east|west)`;
    pattern := `(?i)^\s*\w+(\s+` + dir +`=[0-9a-zA-Z]+\s*)*\s*$`
    //match, _ := regexp.MatchString(pattern, str);
    r, _ := regexp.Compile(pattern);
    for i := range input {
		line := input[i]
		if r.MatchString(line) == false {
			logger.Println("error : Invalid entry in map file at line", i, line);
			return errors.New(fmt.Sprintf(`error : Invalid entry in map file at line %d :
			%s Correct Usage : City [Direction=City Direction=City]`, i, line));
		}
    }
    return nil;
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
	if inv.totalMoves == 0 {		
		fmt.Printf("Deploy : => ");
	} else {
		fmt.Printf("Move %d : => ", inv.totalMoves);
	}

	fmt.Printf("City %s is destroyed by %d and %d \n",
		cityTo, alien, alien2);
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
	//logger.Println("----> moveAliens : ", inv.totalMoves);
	var roads map[string]string;
	for alien, city := range inv.aliens {
		roads = inv.cityMap[city].roads;
		if (len(roads) == 0) {
			//City is island
			//Alien is trapped and destroyed itself
			delete(inv.aliens, alien);
			fmt.Printf("Alien %d is trapped on %s island and destroyed himself\n", alien, city);
			logger.Printf("Alien %d is trapped on %s island and destroyed himself\n", alien, city);
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
    f, err := os.Create("./examples/output.txt")
    defer f.Close()
    if err != nil {
    	logger.Println("Error opening output file", err);
        panic(err)
    }

    if _, err = f.WriteString(output); err != nil {
    	logger.Println("Error writing to output file", err);
        panic(err)
    }
}

func ReadFile(mapFile string) string {
	var data, err = ioutil.ReadFile(mapFile);
	if (err != nil) {
		log.Fatal(err);
		logger.Println("Error reading map file", err);
	}
	return string(data);
}

//simulates the alien invasion.
func Run(totalAliens int, mapFile string) {
	logger.Println("INIT");
	mapDir := "./examples/";
	data := ReadFile(mapDir + mapFile);
	inv := New();
	//BUILD MAP
	inv.BuildMap(data);

	//Deploy Aliens
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