package system

import (
	"github.com/n0rad/go-erlog/logs"
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"io/ioutil"
	"os"
	"testing"
)

func TestRsyncInit(t *testing.T) {
	r := Rsync{}

	if err := r.Init(); err == nil {
		t.Errorf("should fail if no mountpoints")
	}

	//

	r.TargetFilesystem.Mountpoint = "/tmp/target"
	r.SourceFilesystem.Mountpoint = "/tmp/source"
	r.SourceInFilesystemPath = "yopla"
	r.TargetInFilesystemPath = "yopla"

	if err := r.Init(); err != nil {
		t.Errorf("should be fully usable")
	}
}

func TestRsyncSync(t *testing.T) {
	source, _ := ioutil.TempDir("/tmp", "hdmtest-source")
	target, _ := ioutil.TempDir("/tmp", "hdmtest-source")

	_ = os.Mkdir(source+"/ss ss", 0777)
	_ = os.Mkdir(target+"/tt tt", 0777)
	_ = ioutil.WriteFile(source+"/ss ss/file", []byte("yopla"), 0644)

	server := Server{runner: runner.LocalRunner{/*UnSudo: true*/}}
	r := Rsync{
		SourceFilesystem: BlockDevice{
			Mountpoint: source,
			server:     &server,
		},
		TargetFilesystem: BlockDevice{
			Mountpoint: target,
			server:     &server,
		},
		SourceInFilesystemPath: "ss ss",
		TargetInFilesystemPath: "tt tt",
	}
	if err := r.Init(); err != nil {
		t.Errorf("%s", err)
	}

	targetSize, err := r.TargetSize()
	if err != nil {
		logs.WithE(err).Warn("fail")
		t.Errorf("Should get target size, %s", err)
	}
	if targetSize > 0 {
		t.Errorf("target size should be 0, %d", targetSize)
	}

	sourceSize, err := r.SourceSize()
	if err != nil {
		logs.WithE(err).Warn("fail")
		t.Errorf("Should get target size, %s", err)
	}
	if sourceSize != 4 {
		t.Errorf("target size should be 0, %d", sourceSize)
	}

	// TODO: df on target
	//if why, err := r.Rsyncable(); err != nil || why != nil {
	//	t.Errorf("not rsyncable %s %s", why, err)
	//}

}

//if err := r.; err == nil {
//t.Errorf("should fail if no mountpoints")
//}
//
//if ko.MaxServerConnectionAge != keepalive.Infinity {
//t.Errorf("%s maximum connection age %v", t.HandlerName(), ko.MaxServerConnectionAge)
//}
//t.Errorf("should fail if no mountpoints %s %s", t.HandlerName(), err.Error())
