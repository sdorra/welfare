package packages

import (
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

const keyListing = `/etc/apt/trusted.gpg
--------------------
pub   1024D/437D05B5 2004-09-12
uid                  Ubuntu Archive Automatic Signing Key <ftpmaster@ubuntu.com>
sub   2048g/79164387 2004-09-12

pub   4096R/C0B21F32 2012-05-11
uid                  Ubuntu Archive Automatic Signing Key (2012) <ftpmaster@ubuntu.com>

pub   4096R/EFE21092 2012-05-11
uid                  Ubuntu CD Image Automatic Signing Key (2012) <cdimage@ubuntu.com>

pub   1024D/FBB75451 2004-12-30
uid                  Ubuntu CD Image Automatic Signing Key <cdimage@ubuntu.com>

pub   2048R/D742B261 2014-06-18
uid                  Sebastian Sdorra <s.sdorra@gmail.com>
sub   2048R/B3C06235 2014-06-18`

func TestContainsKey(t *testing.T) {
	contains, err := containsKey(strings.NewReader(keyListing), "D742B261")
	assert.Nil(t, err)
	assert.True(t, contains)

	contains, err = containsKey(strings.NewReader(keyListing), "C0B21F32")
	assert.Nil(t, err)
	assert.True(t, contains)

	contains, err = containsKey(strings.NewReader(keyListing), "C0B24332")
	assert.Nil(t, err)
	assert.False(t, contains)
}

func TestAptKeyModule_Run(t *testing.T) {
	sys := &testKeySystem{
		isPresent: false,
	}

	key := NewAptKeyModule("D742B261", Present)
	key.system = sys

	change, err := key.Run()
	assert.Nil(t, err)
	assert.True(t, change)

	assert.Equal(t, "D742B261", sys.add)
}

func TestAptKeyModule_RunAlreadyPresent(t *testing.T) {
	sys := &testKeySystem{
		isPresent: true,
	}

	key := NewAptKeyModule("D742B261", Present)
	key.system = sys

	change, err := key.Run()
	assert.Nil(t, err)
	assert.False(t, change)
}

func TestAptKeyModule_RunAbsent(t *testing.T) {
	sys := &testKeySystem{
		isPresent: true,
	}

	key := NewAptKeyModule("D742B261", Absent)
	key.system = sys

	change, err := key.Run()
	assert.Nil(t, err)
	assert.True(t, change)

	assert.Equal(t, "D742B261", sys.remove)
}

func TestAptKeyModule_RunAbsentAlreadyAbsent(t *testing.T) {
	sys := &testKeySystem{
		isPresent: false,
	}

	key := NewAptKeyModule("D742B261", Absent)
	key.system = sys

	change, err := key.Run()
	assert.Nil(t, err)
	assert.False(t, change)
}

type testKeySystem struct {
	add       string
	remove    string
	isPresent bool
}

func (sys *testKeySystem) Add(server string, id string) error {
	sys.add = id
	return nil
}

func (sys *testKeySystem) Remove(id string) error {
	sys.remove = id
	return nil
}

func (sys *testKeySystem) IsPresent(id string) (bool, error) {
	return sys.isPresent, nil
}
