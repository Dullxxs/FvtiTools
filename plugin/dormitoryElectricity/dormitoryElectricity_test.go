package dormitoryElectricity

import (
	"testing"
)

func TestGetRoomQuantity(t *testing.T) {
	t.Log(getRoomQuantity("1号楼#1层#101"))
}
