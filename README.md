# pshunt
Minimalistic process viewer similar to htop. Features:
- view, search, and kill processes
- vi keybindings

![pshunt_demo](https://github.com/user-attachments/assets/dfad35c9-4725-4510-a569-a0023005fb5f)


### Usage:
- `9`: kill process
- `k | up`: move cursor up
- `j | down`: move cursor down
- `/`: search
- `esc`: exit search mode
- `r`: refetch processes
- `ctrl-f`: next page
- `ctrl-b`: previous page
- `H`: top of page
- `L`: bottom of page
- `G`: move cursor to end
- `g`: move cursor to start



### Installation
Just clone the repo and build from source. Requires Go 1.23 compiler.
```
git clone https://github.com/jamesma100/pshunt
cd pshunt
go build -o ./pshunt ./cmd/pshunt/main.go
./pshunt
```
