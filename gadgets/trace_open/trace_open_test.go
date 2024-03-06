package tests

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/inspektor-gadget/inspektor-gadget/integration"
	eventtypes "github.com/inspektor-gadget/inspektor-gadget/pkg/types"
)

// A user would have to implement their own event type as well, therefore not using this from ig repo
type traceOpenEvent struct {
	eventtypes.Event
	eventtypes.WithMountNsID

	Pid      uint32      `json:"pid,omitempty" column:"pid,minWidth:7"`
	Uid      uint32      `json:"uid,omitempty" column:"uid,minWidth:10,hide"`
	Gid      uint32      `json:"gid" column:"gid,template:gid,hide"`
	Comm     string      `json:"comm,omitempty" column:"comm,maxWidth:16"`
	Fd       int         `json:"fd,omitempty" column:"fd,minWidth:2,width:3"`
	Ret      int         `json:"ret,omitempty" column:"ret,width:3,fixed,hide"`
	Err      int         `json:"err,omitempty" column:"err,width:3,fixed"`
	Flags    []string    `json:"flags,omitempty" column:"flags,width:24,hide"`
	FlagsRaw int32       `json:"flagsRaw,omitempty"`
	Mode     string      `json:"mode,omitempty" column:"mode,width:10,hide"`
	ModeRaw  fs.FileMode `json:"modeRaw,omitempty"`
	Path     string      `json:"path,omitempty" column:"path,minWidth:24,width:32"`
	FullPath string      `json:"fullPath,omitempty" column:"fullPath,minWidth:24,width:32" columnTags:"param:full-path"`
}

func TestTraceOpen(t *testing.T) {
	cn := "test-trace-open"
	// handle err
	containerFactory, err := integration.NewContainerFactory("docker")
	if err != nil {
		fmt.Println("error")
	}

	// utils.WrapperForIg()

	traceOpenCmd := &integration.Command{
		Name: "TraceOpen",
		Cmd:  "IG_EXPERIMENTAL=true sudo -E ig run docker.io/pawarpranav83/traceopen:latest",
		ValidateOutput: func(t *testing.T, output string) {
			fmt.Println("hello")
			expectedEntry := &traceOpenEvent{
				Event: eventtypes.Event{
					Type: eventtypes.NORMAL,
					CommonData: eventtypes.CommonData{
						Runtime: eventtypes.BasicRuntimeMetadata{
							RuntimeName:   eventtypes.String2RuntimeName("docker"),
							ContainerName: cn,
						},
					},
				},
				Comm:     "cat",
				Fd:       3,
				Ret:      3,
				Err:      0,
				Path:     "/dev/null",
				FullPath: "",
				Uid:      1000,
				Gid:      1111,
				Flags:    []string{"O_RDONLY"},
				Mode:     "----------",
			}

			normalize := func(e *traceOpenEvent) {
				e.Pid = 0
				e.MountNsID = 0

				e.Runtime.ContainerID = ""
				// TODO: Handle once we support getting ContainerImageName from Docker
				e.Runtime.ContainerImageName = ""
				e.Runtime.ContainerImageDigest = ""
			}

			integration.ExpectEntriesToMatch(t, output, normalize, expectedEntry)

		},
	}
	testSteps := []integration.TestStep{
		containerFactory.NewContainer(cn, "setuidgid 1000:1111 cat /dev/null"),
		traceOpenCmd,
	}

	integration.RunTestSteps(testSteps, t)
}
