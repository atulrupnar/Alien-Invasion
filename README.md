# Alien Invasion

A simulation of fictional aliens invading earth

## About Alien Invasion

This application simulate the alien invasion of earth. 
### Input : 
    1. -map : A world map file
    2. -aliens : Number of Aliens
    3. -out : Output file (world map after the invasion attack)

### Process
* The program reads input file and store the information in an internal data structure
* Randomly deploys aliens on different cities
* Move aliens from city to other connected city(through a road)
* When 2 aliens move to a same city, they fight with each other and destroy themselves, in a process a city is also destroyed
* The simulation stops when all aliens finished or they complete the 10000 moves.
* The application updates the world map

## Key Components and Utility

1. **Read Inputs** : Read inputs such as world map, aliens etc from file/command line
2. **Build Map** : Read World Map and build internal map
3. **Deploy Aliens** : Deploy aliens randomly on different cities
4. **Move Aliens** : Iteratively move aliens from city to connected cities
5. **Destroy Cities** : Destroy city and all its connecting roads and aliens in the city
6. **Validate Map** : Validate input world map file
7. **Stop simulation** : A utility to check for simulation completeion
8. **Update World Map** : A utility to convert internal world map to the original format of input

## Installation
	Prerequisites : Git, Docker, Go

    Method 1 : Build from Source :
	- Clone the repository
    $ git clone https://github.com/atulrupnar/Alien-Invasion.git
	- Run the application (see input section for flag usage)
	$ go run cmd/main --aliens 4
   
	Method 2 : Docker :
	- Clone the repository
    $ git clone https://github.com/atulrupnar/Alien-Invasion.git
    $ sudo docker build .
	$ sudo docker run -it <container-id> -aliens 10

## Sample Output
    $ sudo docker run -it 6728ed5cd516 -aliens 10
    Deploy : => City Mumbai is destroyed by 3 and 1 
    Deploy : => City Hyderabad is destroyed by 9 and 6 
    Move 1 : => City Goa is destroyed by 4 and 2 
    Alien 10 is trapped on Pune island and destroyed himself
    Aliens reached final state

## Testing
	1. Go to test directory : $cd test
	2. Run $go test -v

## Alternative approach :
In the current implementation, Aliens move iteratively one after other. After every move, The program checks whether 2 aliens are in same city. The alien movement can modelled in different ways :
Simalteneous movement : All aliens move independent of each other. At the end of the movement, the program can check for collision of aliens in the city. In that case, there can a multiple aliens collide and destroy at once. The same behaviour can be implemented using concurrent goroutines. Each goroutine will represent one alien and its job is to move from one city to other. The collision check routine(Multiple aliens in same city) will execute after each complete move(movement of all aliens). The current program can be easily converted to the alternative approach with a very little modications to current data structure(and no code modifications only additional code).  

## Todo
- Config file : Define all default variables and file paths to config file such as logfile, output file, map file, maximum number of moves allowed
- Modular code : Common code can be moved to different package
- Optimizations to the code such as convert direction type from string to int
- Implement other approach as mentioned above

## Author
Atul Rupnar

## Assumptions ##
* Alien is a nonzero positive integer
* A city can be a combination of letter, number or _.
* A city can connect to atmost 4 other cities 1 in each direction.

* A city can have only one connection in one direction.
    eg. following example is invalid
    Paris east=Mumbai
    London east=Mumbai

* The roads are bidirectional : The roads connecting to cities are bidirectional meaning that if there is road from city A to city B. City B must have road from B to A in opposite direction

* Aliens moves iteratively from one city to other.
   eg. If 5 aliens are deployed. First alien will make a move first. After the move 2nd alien made a move and so on. The order can be randomized

* A City can be an island. If alien lost on island, it destroy itself.
