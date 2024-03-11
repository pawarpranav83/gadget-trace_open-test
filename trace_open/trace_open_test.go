package tests

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/inspektor-gadget/inspektor-gadget/integration"
	eventtypes "github.com/inspektor-gadget/inspektor-gadget/pkg/types"
	// IG "github.com/pawarpranav83/ig-testing-framework"
)

type EventType string

const ImageName = "docker.io/pawarpranav83/trace_open:latest"

type Event struct {
	eventtypes.CommonData

	// Type indicates the kind of this event
	Type EventType `json:"type"`

	// Message when Type is ERR, WARN, DEBUG or INFO
	Message string `json:"message,omitempty"`
}

// A user would have to implement their own event type as well, therefore not using this from ig repo
type traceOpenEvent struct {
	Event
	eventtypes.WithMountNsID

	Pid      uint32      `json:"pid,omitempty" column:"pid,minWidth:7"`
	Uid      uint32      `json:"uid,omitempty" column:"uid,minWidth:10,hide"`
	Gid      uint32      `json:"gid" column:"gid,template:gid,hide"`
	Comm     string      `json:"comm,omitempty" column:"comm,maxWidth:16"`
	Ret      int         `json:"ret,omitempty" column:"ret,width:3,fixed,hide"`
	Err      int         `json:"err,omitempty" column:"err,width:3,fixed"`
	Flags    int         `json:"flags,omitempty" column:"flags,width:24,hide"`
	FlagsRaw int32       `json:"flagsRaw,omitempty"`
	Mode     int         `json:"mode,omitempty" column:"mode,width:10,hide"`
	ModeRaw  fs.FileMode `json:"modeRaw,omitempty"`
	FName    string      `json:"fname,omitempty" column:"path,minWidth:24,width:32"`
	FullPath string      `json:"fullPath,omitempty" column:"fullPath,minWidth:24,width:32" columnTags:"param:full-path"`
}

// Need to build the gadget image first before running the test
func TestTraceOpen(t *testing.T) {
	cn := "test-trace-open"
	// handle err
	containerFactory, err := integration.NewContainerFactory("docker")
	if err != nil {
		fmt.Println("error")
	}

	traceOpenCmd := &integration.Command{
		Name:         "TraceOpen",
		Cmd:          "ig run docker.io/pawarpranav83/trace_open:latest --runtimes=docker -o json",
		StartAndStop: true,
		ValidateOutput: func(t *testing.T, output string) {
			expectedEntry := &traceOpenEvent{
				Event: Event{
					Type: "",
					CommonData: eventtypes.CommonData{
						Runtime: eventtypes.BasicRuntimeMetadata{
							RuntimeName:   eventtypes.String2RuntimeName("docker"),
							ContainerName: cn,
						},
					},
				},
				Comm:     "cat",
				Ret:      3,
				Err:      0,
				FName:    "/dev/null",
				FullPath: "",
				Uid:      1000,
				Gid:      1111,
				Flags:    0,
				Mode:     0,
			}

			normalize := func(e *traceOpenEvent) {
				e.MountNsID = 0
				e.Pid = 0

				e.Runtime.ContainerID = ""
				e.Runtime.ContainerImageName = ""
				e.Runtime.ContainerImageDigest = ""
			}

			integration.ExpectEntriesToMatch(t, output, normalize, expectedEntry)

		},
	}
	testSteps := []integration.TestStep{
		traceOpenCmd,
		containerFactory.NewContainer(cn, "setuidgid 1000:1111 cat /dev/null"),
	}

	integration.RunTestSteps(testSteps, t)
}
