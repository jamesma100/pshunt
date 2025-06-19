# pshunt
Minimalistic terminal app to hunt and kill processes.

![pshunt_demo](https://github.com/user-attachments/assets/dfad35c9-4725-4510-a569-a0023005fb5f)


### Usage:
- `9`: kill process
- `k | up`: move cursor down
- `j | down`: move cursor down
- `ctrl-f`: next page
- `ctrl-b`: previous page
- `G`: move cursor to end
- `g`: move cursor to start
- `r`: refetch processes
- `/`: search (basic substring search, does not support regex)
- `esc`: exit search mode

### Installation
Just clone the repo and build from source. Requires Go 1.23 compiler.
```
git clone https://github.com/jamesma100/pshunt
cd pshunt
go build -o ./pshunt ./cmd/pshunt/main.go
./pshunt
```
