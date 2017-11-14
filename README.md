# gaemux

** this repository's state is  pre-alpha **

mux for google app engine

# what this do

Many web framework written in go uses build tags to specify build go version.
```
// +build go1.8
```
but it could not  deploy on GAE environment with go1.8.

gaemux slightly wrap`github.com/julienschmidt/httprouter` that doesn't use build tags.



