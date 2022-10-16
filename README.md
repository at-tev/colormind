# COLORMIND
Colormind is a CLI for [colormind](http://colormind.io) web application - the AI color palette generator. Colormind can help you create color schemes 
with different models - color styles from photographs, movies, and popular art.

# Installing
Use `go install` to install the latest version of the Colormind CLI.
```
go install github.com/at-tev/colormind@latest
```

# Usage
There are few ways to get color scheme:
 1. By command *random* in order to get random color scheme
 2. By command *suggest* with 1-4 color code arguments in order to get color scheme suggestion
 
Color schemes are displayed and input as RGB color code by default. Pass `-H|--hex` flag to use hexadecimal color code.

You can use different models by passing model name as `-m|--model` flag value. The models "default" and "ui" are always available, but others will 
change each day at 07:00 UTC. The service will be down for 30 seconds while the new models are refreshed. Use *models* command to know currently 
available models. The default model of color schemes is 'default'.

# License
Colormind is released under the MIT license. See [LICENSE.txt](/LICENSE.txt)
