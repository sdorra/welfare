# Welfare

[![Packagist](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/sdorra/welfare/blob/master/LICENSE)
![Build Status](https://api.travis-ci.org/sdorra/welfare.svg?branch=master)

Welfare is a library for the execution of declarative tasks. 
It is very much inspired by [Ansible](https://www.ansible.com/), but Welfare is designed to be used as embedded component in other applications.

## Usage

```go
copy := files.NewCopyModule("files/issue", "/etc/issue")
copy.FileMode = 0644
changed, _ := copy.Run()
if changed {
    fmt.Println("updated /etc/issue")
}
```
