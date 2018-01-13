# Welfare

Welfare is a library for the execution of declarative tasks. 
It is very much inspired by [Ansible](https://www.ansible.com/), but Welfare is designed to be used as embedded component in other applications.

## Usage

```go
copy := modules.NewCopyModule("files/issue", "/etc/issue")
copy.FileMode = 0644
changed, _ := copy.Run()
if changed {
    fmt.Println("updated /etc/issue")
}
```