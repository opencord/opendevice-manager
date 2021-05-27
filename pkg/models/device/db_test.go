/*
 * Copyright 2020-present Open Networking Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package modifiablecomponent stores ModifiableComponent methods and functions
package device

import (
	"context"
	"fmt"
	"strings"
	"testing"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/opencord/device-management-interface/go/dmi"
	"github.com/opencord/opendevice-manager/pkg/db"
)

// mockModifiableComp refers to mocking of modifiable component req
func mockModifiableComp(id string) *dmi.ModifiableComponent {
	req := new(dmi.ModifiableComponent)
	req.Name = "olt-name-" + id
	req.Alias = "olt-1-alias-" + id
	req.AssetId = "olt-1-assetid-" + id
	req.Uri = new(dmi.Uri)
	req.Uri.Uri = "127.0.0." + id
	req.Parent = new(dmi.Component)
	req.AdminState = dmi.ComponentAdminState_COMP_ADMIN_STATE_UNLOCKED
	return req
}

func mockHardware(id string) *dmi.Hardware {
	hw := new(dmi.Hardware)
	hw.LastChange = new(timestamp.Timestamp)
	hw.LastBooted = new(timestamp.Timestamp)
	hw.Root = new(dmi.Component)
	hw.Root.ModelName = "model-name-" + id
	return hw
}

func runTcase(tcaseName string, t *testing.T, tcaseFunc func() (bool, error)) (string, bool) {
	fmt.Println("\n#======= TESTCASE STARTED : " + tcaseName + " ========#")
	if ok, err := tcaseFunc(); !ok {
		fmt.Println("#======= TESTCASE FAILED : "+tcaseName+"  ========#", err)
		return tcaseName, false
	}
	fmt.Println("#======= TESTCASE PASSED : " + tcaseName + " ========#")
	return tcaseName, true
}

// _Test_PositiveTcaseNewDeviceRecord refers to the positive tcase defined for testing func NewDeviceRecord
func Test_PositiveTcaseNewDeviceRecord(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord

	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for NewDeviceRecord
	tcase1 := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for NewDeviceRecord-1", t, tcase1); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}
}

// _Test_NegativeTcaseDBGetByName refers to the negative tcase defined for testing func DBGetByName
func Test_NegativeTcaseDBGetByName(t *testing.T) {
	req := mockModifiableComp("1")

	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Negative Testcase for DBGetByName
	tcase2 := func() (bool, error) {
		if rec, err := DBGetByName(ctx, req.Name); rec != nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Negative Testcase for DBGetByName-1", t, tcase2); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}
}

// _Test_Suite refers to all component testcases belongs to all packages
func Test_Suite(t *testing.T) {
}

// _Test_NegativeTcaseDBAddByName refers to the negative tcase defined for testing func DBAddByName
func Test_NegativeTcaseDBAddByName(t *testing.T) {

	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Negative Testcase for DBAddByName
	tcase3 := func() (bool, error) {
		emptyDevRec := new(DeviceRecord)
		if err := emptyDevRec.DBAddByName(ctx); err == nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Negative Testcase for DBAddByName-1", t, tcase3); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_PositiveTcaseDBAddByName refers to the positive tcase defined for testing func DBAddByName
func Test_PositiveTcaseDBAddByName(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord
	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for DBAddByName
	tcase4 := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		if err = rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for DBAddByName-1", t, tcase4); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_PositiveTcaseDBGetByName refers to the positive tcase defined for testing func DBGetByName
func Test_PositiveTcaseDBGetByName(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord
	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for DBGetByName
	tcase5 := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		if err = rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		if rec, err := DBGetByName(ctx, rec.Name); rec == nil || err != nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for DBGetByName-1", t, tcase5); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_PositiveTcaseDBGetByNameWithCacheMiss refers to the positive tcase defined for testing func DBGetByName with cache miss
func Test_PositiveTcaseDBGetByNameWithCacheMiss(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord
	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for DBGetByName with cache miss
	tcase5 := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		if err = rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		if rec, err := DBGetByName(ctx, rec.Name); rec == nil || err != nil {
			return false, err
		}
		ClearCacheEntry(ctx, rec.Name, "")
		if rec, err := DBGetByName(ctx, rec.Name); rec == nil || err != nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for DBGetByName-2 with cache miss", t, tcase5); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_PositiveTcaseDBAddUuidLookup refers to the positive tcase defined for testing func DBAddUuidLookup
func Test_PositiveTcaseDBAddUuidLookup(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord
	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for DBAddUuidLookup
	tcase6 := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		if err = rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		rec.Uuid = strings.Replace(rec.Name, "name", "uuid", 1)
		if err := rec.DBAddUuidLookup(ctx); err != nil {
			return false, err
		}
		if err := rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		if rec, err := DBGetByName(ctx, rec.Name); rec == nil || err != nil || rec.Uuid == "" {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for DBAddUuidLookup-1", t, tcase6); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_PositiveTcaseDBGetAll refers to the positive tcase defined for testing func DBGetAll
func Test_PositiveTcaseDBGetAll(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord
	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for DBGetAll
	tcase8 := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		if err = rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		if list, err := DBGetAll(ctx); list == nil || err != nil || len(list) != 1 {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for DBGetAll-1", t, tcase8); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_NegativeTcaseDBGetByUuid refers to the negative tcase defined for testing func DBGetByUuid
func Test_NegativeTcaseDBGetByUuid(t *testing.T) {

	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Negative Testcase for DBGetByUuid
	tcase10 := func() (bool, error) {
		if rec, err := DBGetByUuid(ctx, "invalid-uuid-1"); rec != nil || err == nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Negative Testcase for DBGetByUuid-1", t, tcase10); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_PositiveTcaseDBDelRecord refers to the positive tcase defined for testing func DBDelRecord
func Test_PositiveTcaseDBDelRecord(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord
	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for DBDelRecord
	tcase11 := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		if err = rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		if err := rec.DBDelRecord(ctx); err != nil {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for DBDelRecord-1", t, tcase11); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}

// _Test_PositiveTcaseDBSaveHwInfo refers to the positive tcase defined for testing func DBSaveHwInfo
func Test_PositiveTcaseDBSaveHwInfo(t *testing.T) {
	req := mockModifiableComp("1")
	var rec *DeviceRecord
	db.MockKVClient()
	ctx := context.Background()
	defer db.ClearCache()

	// Positive Testcase for DBSaveHwInfo
	tcase := func() (bool, error) {
		var err error
		if rec, err = NewDeviceRecord(ctx, req); rec == nil || err != nil {
			return false, err
		}
		if err = rec.DBAddByName(ctx); err != nil {
			return false, err
		}
		rec.Uuid = strings.Replace(rec.Name, "name", "uuid", 1)
		if err := rec.DBAddUuidLookup(ctx); err != nil {
			return false, err
		}
		hwInfo := mockHardware("1")
		if err := rec.DBSaveHwInfo(ctx, hwInfo); rec == nil || err != nil {
			return false, err
		}
		if rec, err := DBGetByName(ctx, rec.Name); rec == nil || err != nil || rec.ModelName != hwInfo.Root.ModelName {
			return false, err
		}
		return true, nil
	}

	if name, ok := runTcase("Positive Testcase for DBSaveHwInfo-1", t, tcase); !ok {
		t.Errorf("#======= FAILED :  Testcase " + name + "  ========#")
		return
	}

}
