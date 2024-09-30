# What is this?
small attemp of a tui tool to change the mouse sensibility and acceleration using xinput Coordinate Transformation Matrix property

![example](./demo.gif)
gif made with [vhs](https://github.com/charmbracelet/vhs/)

# Why?
I just wanted to change my mouse sensibility and found out that the only way was to use xinput and messing around with a matrix
I tried piper and my mouse was not supported, I did not found any ui or tui that handle this (maybe there are) so I decided to make one

# How to install?

## With Go
```bash
go install github.com/jhonnyV-V/xinput-mouse-matrix@latest
```

## Compile it yourself

```bash
git clone git@github.com:jhonnyV-V/xinput-mouse-matrix && cd xinput-mouse-matrix
```
```bash
go build -o xinput-mouse-matrix
```
```bash
sudo mv ./xinput-mouse-matrix /usr/local/bin
```
