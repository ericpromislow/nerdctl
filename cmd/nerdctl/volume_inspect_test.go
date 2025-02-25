/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/containerd/nerdctl/pkg/inspecttypes/native"
	"github.com/containerd/nerdctl/pkg/testutil"
	"gotest.tools/v3/assert"
)

func TestVolumeInspectContainsLabels(t *testing.T) {
	t.Parallel()
	testVolume := testutil.Identifier(t)

	base := testutil.NewBase(t)
	base.Cmd("volume", "create", "--label", "tag=testVolume", testVolume).AssertOK()
	defer base.Cmd("volume", "rm", "-f", testVolume).Run()

	inspect := base.InspectVolume(testVolume)
	inspectNerdctlLabels := (*inspect.Labels)
	expected := make(map[string]string, 1)
	expected["tag"] = "testVolume"
	assert.DeepEqual(base.T, expected, inspectNerdctlLabels)
}

func TestVolumeInspectSize(t *testing.T) {
	testutil.DockerIncompatible(t)
	t.Parallel()
	testVolume := testutil.Identifier(t)
	base := testutil.NewBase(t)
	base.Cmd("volume", "create", testVolume).AssertOK()
	defer base.Cmd("volume", "rm", "-f", testVolume).Run()

	cmdResult := base.Cmd("volume", "inspect", "--size", testVolume).Run()
	var dc []native.Volume
	err := json.Unmarshal([]byte(cmdResult.Stdout()), &dc)
	assert.NilError(t, err)
	var mountpoint = dc[0].Mountpoint

	token := make([]byte, 1028)
	rand.Read(token)
	err = os.WriteFile(filepath.Join(mountpoint, "test-file"), token, 0644)
	assert.NilError(t, err)
	cmdResult = base.Cmd("volume", "inspect", "--size", testVolume).Run()
	assert.NilError(t, err)
	err = json.Unmarshal([]byte(cmdResult.Stdout()), &dc)
	assert.NilError(t, err)
	assert.Equal(t, dc[0].Size, int64(1028))
}
