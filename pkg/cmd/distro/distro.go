package distro

import (
	"fmt"
	"os"
)

type OS int64

const (
	Linux OS = iota
	Mac
	Windows
)

func (os OS) String() string {
	switch os {
	case Linux:
		return "linux"
	case Mac:
		return "mac"
	case Windows:
		return "windows"
	}
	return "Unknown OS"
}

type LinuxDistro int64

const (
	Debian LinuxDistro = iota
	Ubuntu
	LinuxMint
	CentOS
	Fedora
	ArchLinux
	Manjaro
)

func (ld LinuxDistro) String() string {
	switch ld {
	case Debian:
		return "Debian"
	case Ubuntu:
		return "Ubuntu"
	case LinuxMint:
		return "LinuxMint"
	case CentOS:
		return "Cent OS"
	case Fedora:
		return "Fedora"
	case ArchLinux:
		return "Arch Linux"
	case Manjaro:
		return "Manjaro"
	}
	return "Unknown Linux Distribution"
}

type MacDistroType int64

const (
	MacOS MacDistroType = iota
)

type WindowsDistroType int64

const (
	Windows11 WindowsDistroType = iota
	Windows10
)

func (wd WindowsDistroType) String() string {
	switch wd {
	case Windows11:
		return "Windows 11"
	case Windows10:
		return "Windows 10"
	}
	return "Unknown Windows Distrubtion"
}

type Distro struct {
	OS             OS
	DistroType     interface{}
	Version        string
	ISOSource      string
	ISODestination string
}

func DefaultDir() string {
	if current_dir, err := os.Getwd(); err == nil {
		return fmt.Sprintf("%s/files", current_dir)
	}
	return ""
}
