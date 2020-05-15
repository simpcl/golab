package strategy

import "testing"

func TestAddData16MB(t *testing.T) {
	dobj := NewDataObject()
	dobj.AddData(1024 * 1024 * 16)
	t.Logf("dataobject: %v", dobj)
	for dpg_index, dpg := range dobj.dpgs {
		t.Logf("datapartgroup: %d %v", dpg_index, dpg)
		for dp_index, dp := range dpg.parts {
			t.Logf("datapart: %d %v", dp_index, dp)
		}
	}
}

func TestAddData64MB(t *testing.T) {
	dobj := NewDataObject()
	dobj.AddData(1024 * 1024 * 64)
	t.Logf("dataobject: %v", dobj)
	for dpg_index, dpg := range dobj.dpgs {
		t.Logf("datapartgroup: %d %v", dpg_index, dpg)
		for dp_index, dp := range dpg.parts {
			t.Logf("datapart: %d %v", dp_index, dp)
		}
	}
}

func TestAddData2GB(t *testing.T) {
	dobj := NewDataObject()
	dobj.AddData(1024 * 1024 * 1024 * 2)
	t.Logf("dataobject: %v", dobj)
	for dpg_index, dpg := range dobj.dpgs {
		t.Logf("datapartgroup: %d %v", dpg_index, dpg)
		for dp_index, dp := range dpg.parts {
			t.Logf("datapart: %d %v", dp_index, dp)
		}
	}
}

func TestAddData5GB(t *testing.T) {
	dobj := NewDataObject()
	dobj.AddData(1024 * 1024 * 1024 * 5)
	t.Logf("dataobject: %v", dobj)
	for dpg_index, dpg := range dobj.dpgs {
		t.Logf("datapartgroup: %d %v", dpg_index, dpg)
		for dp_index, dp := range dpg.parts {
			t.Logf("datapart: %d %v", dp_index, dp)
		}
	}
}

func TestAddData1TB(t *testing.T) {
	dobj := NewDataObject()
	dobj.AddData(1024 * 1024 * 1024 * 1024 * 1)
	t.Logf("dataobject: %v", dobj)
	for dpg_index, dpg := range dobj.dpgs {
		t.Logf("datapartgroup: %d %v", dpg_index, dpg)
		for dp_index, dp := range dpg.parts {
			t.Logf("datapart: %d %v", dp_index, dp)
		}
	}
}
