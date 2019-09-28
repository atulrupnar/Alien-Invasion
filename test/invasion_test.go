package invasionTest

import (
    "testing"
    "../src"
)

var mapFile = "../examples/world_map1.txt";
var inv = invasion.New();
func TestDirection(t *testing.T) {
    dir := invasion.OppositeDirection("south");
    if (dir != "north") {
		t.Fatalf("Wrong Opposite Direction: got %s, want: %s", dir, "north");    	
    }
}

func TestReadFile(t *testing.T) {
	data := invasion.ReadFile(mapFile);
	if len(data) == 0 {
		t.Fatalf("Unable to read file : got file length %d, want %d", len(data), 368);
	};
}

func TestBuildMap(t *testing.T) {
	data := invasion.ReadFile(mapFile);
	inv.BuildMap(data);
	cities := inv.GetAllCities();
	if (len(cities) != 10) {
		t.Fatalf("Unable to build map properly, map got %d cities, wants %d cities", len(cities), 10);
	}
}

func TestIsValidInput(t *testing.T) {
	str1 := []string{"Mumbai South=Goa east=Pune north=Delhi"};	
	if err := invasion.IsValidInput(str1); err != nil {
		t.Fatalf("Validation fails for right input");
	};

	str2 := []string{"Mumbai South=Goa east=Pune north="};
	if err := invasion.IsValidInput(str2); err == nil {
		t.Fatalf("Validates wrong input");
	};
}

func TestDestroyCity(t *testing.T) {
	inv.DestroyCity("Mumbai", 0);
	cityMap := inv.GetCityMap();
	if _, ok := cityMap["Mumbai"]; ok {
		t.Fatalf("Destroy city failed to delete city Mumbai from map");
	}
}

func TestDeployAliens(t *testing.T) {
	inv.DeployAliens(3);
	aliens := inv.GetAliens()
	if (len(aliens) == 0) {
		t.Fatalf("Deploy aliens failed, Got 0 aliens, wants at least 1");
	}
}

func TestSimulation(t *testing.T) {
	inv2 := invasion.New();
	inv2.Run(4, "../examples/world_map1.txt", "../examples/output.txt");
	if !inv2.IsSimulationOver() {
		t.Fatalf("Simulation failed...!");
	}
}