# Alien Invasion

A simulation of fictional aliens invading earth

## About Alien Invasion

This application simulate the alien invasion of earth. 
### Input : 
    1. A world map of city and roads connecting to other cities
    2. Number of Aliens
   
### Process
* The program reads inpput file and store the information in an internal data structure
* Randomly deploys aliens on different cities
* Move aliens from city to other connected city(through a road)
* When 2 aliens move to a same city, they fight with each other and destroy themselves, in a process a city is also destroyed
* The simulation stops when all aliens finished or they complete the 10000 moves.
* The application updates the world map

## Key Components and Utility

1. **Read Inputs** : Read inputs such as world map, aliens etc from file/command line
2. **Build Map** : Read World Map from input file and build internal map
3. **Deploy Aliens** : Deploy aliens randomly on different cities
4. **Move Aliens** : Iteratively move aliens from city to connected cities
5. **Destroy Cities** : Destroy city and all its connecting roads
6. **Validate Map** : Validate input world map file
7. **Stop simulation** : A utility to check for simulation completeion
8. **Update World Map** : A utility to convert internal world map to the original format of input

## Assumptions ##

* A city can connect to atmost 4 other cities 1 in each direction.

* A city can have only one connection in one direction.
    eg. following example is invalid
    Paris east=Mumbai
    London east=Mumbai

* The roads are bidirectional : The roads connecting to cities are bidirectional meaning that if there is road from city A to city B. City B must have road from B to A in opposite direction

* Aliens moves iteratively from one city to other.
   eg. If 5 aliens are deployed. First alien will make a move first. After the move 2nd alien made a move and so on. The order can be randomized