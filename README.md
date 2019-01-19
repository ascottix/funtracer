# Funtracer

A simple raytracer in Go.

Funtracer is based on [The Ray Tracer Challenge](https://pragprog.com/book/jbtracer/the-ray-tracer-challenge "The Ray Tracer Challenge") book by Jamis Buck.

Now (slowly) following the [Physically Based Rendering](http://www.pbr-book.org/ "Physically Based Rendering") book by Matt Pharr, Wenzel Jacob and Greg Humphreys.

![Stanford Dragon](https://ascottix.github.io/funtracer/dragon.png)

## Features

The main goal of Funtracer is to have fun and learn Go with a project that can offer interesting and sometimes difficult challenges but nice rewards too.

So far it has worked really well (the way the book is designed helps a lot) and although the program cannot boast any particularly amazing feature, it takes great pride and satisfaction in what it _can_ do:

- Basic shapes: cone, cube, cylinder, plane, sphere, triangle meshes
- Groups
- Constructive Solid Geometry (CSG)
- Bounding Volume Hierarchies (BVH) with the Surface Area Heuristic (SAH)
- Color patterns
- Lights: directional, point, spot
- Import .fun, .ray and .obj files
- Parallel rendering

## How to build

Go takes care of everything, just run:

`go build`

As the program was developed using a test-driven approach, there are a number of tests and examples to try too. To run them all use:

`go test`

Some tests will generate a PNG file, look for files named `test_*.png`.

## How to use

Funtracer takes a scene file and converts it into a PNG image. A number of options may also be specified, run the program without arguments to display usage information.

Note that any option, if present, must be specified **before** the name of the scene file, or it will be ignored.

If no scene file is specified, the program looks for a file named `have.fun` and reads it if present.

The directory `scenes` contains several scene file to try. Try for example:

`./funtracer scenes/teapot`

or if you can wait a bit more add supersampling:

`./funtracer -ss 4 scenes/teapot`

The `-ss 4` option renders the scene 16 times, sampling each pixel at a slightly different position every time. This helps reduce aliasing and other visual artifacts.

Options may also be specified in a `config.json` file, for example:

    {
        "nt":   1
    }

will override the default value for the number of threads and set it to 1.

### Scene file format

Scenes are described in plain text using a simple description language. 

The format is based on the .ray format explained here:

[Working With .ray Files](http://www.cs.cmu.edu/afs/cs.cmu.edu/academic/class/15864-s04/www/assignment4/format.html "Working With .ray Files")

Many features have been added and there could be some minor differences as well, the best way to learn is to open some scenes in a text editor and play with them.

## License

See LICENSE file.
