package silo

import (
	"fmt"
	"sparrow/utils"
	"testing"
)

func estPatchReplaceSimpleAts(t *testing.T) {
	initSilo()

	rs := insertRs(patchDevice)
	pr := getPr(`{"Operations":[{"op":"rplace", "value":{"installedDate": "2016-06-18T14:19:14Z"}}]}`, deviceType)

	updatedRs, err := sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	assertIndexVal(deviceType.Name, "installedDate", utils.GetTimeMillis("2016-05-17T14:19:14Z"), false, t)
	assertIndexVal(deviceType.Name, "installedDate", utils.GetTimeMillis("2016-06-18T14:19:14Z"), true, t)

	// apply the same patch on the already updated resource, resource should not get modified
	notUpdatedRs, err := sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	originalMeta := updatedRs.GetMeta().GetFirstSubAt()
	newMeta := notUpdatedRs.GetMeta().GetFirstSubAt()

	assertEquals(t, "meta.created", notUpdatedRs, originalMeta["created"].Values[0])
	assertEquals(t, "meta.version", notUpdatedRs, fmt.Sprint(originalMeta["lastmodified"].Values[0]))
	if originalMeta["lastmodified"].Values[0] != newMeta["lastmodified"].Values[0] {
		t.Errorf("Patch operation modified though the attribute data is unchanged")
	}

	pr = getPr(`{"Operations":[{"op":"replace", "path": "location.latitude", "value": "20°10'45.4\"N"}]}`, deviceType)

	updatedRs, err = sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	assertIndexVal(deviceType.Name, "location.latitude", "19°10'45.4\"N", false, t)
	assertIndexVal(deviceType.Name, "location.latitude", "20°10'45.4\"N", true, t)

	pr = getPr(`{"Operations":[{"op":"replace", "path": "macId", "value": "6A"}]}`, deviceType)

	updatedRs, err = sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	macId := updatedRs.GetAttr("macId").GetSimpleAt().Values[0].(string)
	if macId != "6A" {
		t.Error("macId attribute was not added")
	}
}

func TestAddMultiValuedSubAt(t *testing.T) {
	initSilo()

	rs := insertRs(patchDevice)
	pr := getPr(`{"Operations":[{"op":"replace", "path": "photos.display", "value": "photo display"}]}`, deviceType)

	updatedRs, err := sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	photos := updatedRs.GetAttr("photos").GetComplexAt()
	for _, subAtMap := range photos.SubAts {
		val := subAtMap["display"].Values[0].(string)
		if val != "photo display" {
			t.Errorf("Failed to add display value of photos attribute")
		}
	}
}

func TestReplaceSingleCA(t *testing.T) {
	initSilo()

	rs := insertRs(patchDevice)
	pr := getPr(`{"Operations":[{"op":"replace", "path": "location", "value": {"latitude": "20°10'45.4\"N", "desc": "kodihalli"}}]}`, deviceType)

	updatedRs, err := sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	assertIndexVal(deviceType.Name, "location.latitude", "19°10'45.4\"N", false, t)
	assertIndexVal(deviceType.Name, "location.latitude", "20°10'45.4\"N", true, t)

	desc := updatedRs.GetAttr("location.desc").GetSimpleAt().Values[0].(string)
	if desc != "kodihalli" {
		t.Error("desc attribute was not added")
	}

}

/*
   o  If the target location is a multi-valued attribute and a value
      selection ("valuePath") filter is specified that matches one or
      more values of the multi-valued attribute, then all matching
      record values SHALL be replaced.
*/
func TestReplaceMultiCA(t *testing.T) {
	initSilo()

	rs := insertRs(patchDevice)
	pr := getPr(`{"Operations":[{"op":"replace", "path": "photos[value pr]", "value": {"value": "1.jpg", "display": "added display"}}]}`, deviceType)

	updatedRs, err := sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	assertIndexVal(deviceType.Name, "photos.value", "abc.jpg", false, t)
	assertIndexVal(deviceType.Name, "photos.value", "xyz.jpg", false, t)
	assertIndexVal(deviceType.Name, "photos.value", "1.jpg", true, t)

	displayAt := updatedRs.GetAttr("photos").GetComplexAt()
	for _, subAtMap := range displayAt.SubAts {
		display := subAtMap["display"].Values[0].(string)
		if display != "added display" {
			t.Error("display attribute was not added")
		}
	}
}

/*
   o  If the target location is a complex multi-valued attribute with a
      value selection filter ("valuePath") and a specific sub-attribute
      (e.g., "addresses[type eq "work"].streetAddress"), the matching
      sub-attribute of all matching records is replaced.
*/
func TestMultival(t *testing.T) {
	initSilo()

	rs := insertRs(patchDevice)
	pr := getPr(`{"Operations":[{"op":"replace", "path": "photos[value pr].display", "value": "this is a photo"}]}`, deviceType)

	updatedRs, err := sl.Patch(rs.GetId(), pr, deviceType)
	if err != nil {
		t.Errorf("Failed to apply patch req")
	}

	photos := updatedRs.GetAttr("photos").GetComplexAt()
	for _, saMap := range photos.SubAts {
		display := saMap["display"].Values[0].(string)
		if display != "this is a photo" {
			t.Error("display attribute was not added")
		}
	}
}