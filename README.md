# Package-Hell

# Problem/Motivation:
Before installing a node package, if you would like to know the total install size of
the package, this is for you. The current npmjs.com website shows only the publish size 
of the package and does not give the total size of the package including its external
dependencies and external dependencies of those dependencies.This project aims to 
solve that problem by providing an estimation of the total size of the latest stable 
version of all those external dependencies. 

# To Get Started

Clone the repo:
`git clone https://github.com/CoderParth/package-hell.git`

Build:
`go build`

Or Run Directly:
`go run .`

# Usage:
You will be prompted to enter a package name, and will get results ASAP.

Here is a small demonstration:

`$ go run .`

`Enter a package name:`

- Let us say you entered "react"

You can expect a result like the below:

`List of all the dependant packages and their size`
`react : 318101`
`loeose-envify : 5814`
`jse-tokens : 16067`

`Estimated Total Size: 0.32 MB`



