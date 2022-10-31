# Mazes for Programmers

This book is a Go implementation for the book [Mazes for Programmers](http://www.mazesforprogrammers.com/)
by James Buck.

It will be used mainly as a tool to learn Go, this book was chosen for
several reasons:

* It is very well written and provides fun challenges, thus keeping the reader interested
* The mazes are implemented using classic algorithms and data structures,
thus making you think
* It delves into graphically representing the mazes, hence a graphics library is needed. This is always fun
* Code samples are written in Ruby, being the high-level language that it is, Ruby provides a lot of "magic" built-in
that Go simply hasn't got. Thus, it forces us to find the most idiomatic way in Go of implementing these features; this
provides a very hands-on way of getting to know the language.

Code samples will be implemented pure Go, whereas for the graphics the Go bindings of [raylib](https://github.com/gen2brain/raylib-go)
will be used. Raylib is a simple library written in C (so it's written in a declarative way) that wraps the OpenGL API, it's made mainly to teach game development;
thus its API is very simplified while at the same time providing everything that is necessary for a 2D (and even 3D) application.

This module is structure like this (work in progress):

* `mfp` contains all the logic for the mazes, abstracted from graphic representation where possible
* `main` contains the application that calls the code defined in the above package
