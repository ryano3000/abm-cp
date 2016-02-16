###Note that all code is alpha until version 1.0

###Current Version: 0.4.0a

![abm preview](https://giant.gfycat.com/NippyTidyCassowary.gif)

##Context
Software for computing a Predator-Prey Agent Based Model of prey colour polymorphism (CP) in Go. 

Written specifically to assist the research by Dr. James Dale in evaluating hypotheses for the evolutionary maintenance of extreme colour polymorphism.

##Goal

Server-side computation, client-side GUI controls and visualisation of the running ABM in a web browser.

Generalise system so any `abm` package with local model context, agent design, etc can be hot-swapped in (linked via a command line flag) and run.

###How to use:

Install Go from [here](https://golang.org/dl/).

Download these repositories:  run `go get -u github.com/benjamin-rood/abm-cp && go get -u github.com/benjamin-rood/gobr` 

Change current directory to `$GOPATH/src/github.com/benjamin-rood/abm-cp` and run `cd net && go build && ./net`

Point web browser at `localhost:9999`

Current version only tested on Safari on OS X and Chrome on Windows 8.1.

Known issue:

Sometimes when hitting the "Start ABM..." button on the app (in the middle of an existing run) nothing happens.

Error message (server-side): 
` SignalHub.register() failed: receiver signature already exists.
 Clash when registering Model: wo vis: for sync with m.turnSignal `

Meaning: I'm not cleaning up properly inside of an individual session when resetting the modelling contextual parameters, including deleting the registered receiver for `Model.vis()`

Workaround: just reload the browser window. 



## Roadmap

### 0.1.0
Base requirements completed.

* A simple interface for the CP Prey ABM, with the visualisation on the left, with contextual parameter controls on the right. :white_check_mark:

* Use P5js for render viewport. :white_check_mark:

* Browser recieves drawing instructions from ABM which get loaded into the P5 instance's draw loop. :white_check_mark:

* Responsive design, visualisation (render viewport) scales to the dimensions of available browser real estate. :white_check_mark:

* Server handles concurrent bi-directional connections for concurrent ABM modelling individual to each user, with data sent between client and server using Web Sockets. :white_check_mark:

* Server cleans up after user disconnections. :white_check_mark:

* Server receives new submitted contextual parameters as signal to start/restart ABM. :white_check_mark:

* Serialisation of data messages using JSON (prohibatively slow for anything approaching 10,000 agents).  :white_check_mark:

* CP Prey agents implementation:
	 * Rule-Based-Behaviour. :white_check_mark:
	 * Asexual Reproduction. :white_check_mark:
	 * Mutation (colouration). :white_check_mark:

* Visual Predator implementation:
	* Rule-Based-Behaviour. :white_check_mark:
	* Exhaustive Prey Search (very slow). :white_check_mark:
	* Colour Imprinting (needs tweaking, no baseline yet established). :white_check_mark:

* Simple concurrent execution of Predator/Prey RBB. :white_check_mark:

* Unit tests for `geometry`, `calc`, `render` packages. :white_check_mark:

### 0.2.0

* Dispatch errors along channels, use external handlers who receive the errors and process/print them. :white_check_mark:

* Essential unit tests for `abm` package :white_check_mark:

* Show live population and timeline statistics inside P5js viewport. :white_check_mark:
 
* Visual Predator implemenation:
	* Find baseline params for Colour Imprinting. :white_check_mark:
	* Adaptation in response to hunger. :white_check_mark:
	* Starvation ⟹ Death. :white_check_mark:
	* Sexual Reproduction. :white_check_mark:

* General modelling and interactions between agent types in place, with all baseline parameters set for end-use. :white_check_mark:

### 0.3.0

* Toggle Visualisation on/off :white_check_mark:

* Record data of a running model to log files as JSON for statistical analysis as a concurrent process, independent of Visualisation or model computation. Toggleable. :white_check_mark:

* User determined frequency of data logging on client side. :white_check_mark:

* Better model for Predator specialisation through colour-imprinting which directly gives search/attack advantage, rather than being decided randomly. :white_check_mark:

* Complete tests for `abm` package :white_check_mark:

### 0.4.0

* Switch data serialisation to Protocol Buffers (protobuf) ~10X speedup. Marshalling drawing instructions to JSON is currently the single most expensive action!

* Better Prey Search (using grid system). Will be at least a 2X speedup.

* Client-side input validation.

### 0.5.0

* Have complete control over ABM computation, logging, visualisation from command-line, rather than just starting up a web server and controlling through the (limited) in-browser Visualisation GUI.

* Import JSON-formatted text files as pre-defined modelling Context.

* Use JSON-formatted logging for debug only.

* Switch to a compressed binary encoding for log files. 


### 1.0.0

* Switch all public html file to templated/generated ones based on context parameters etc. Send JSON to client, let Javascript API process that into the formatted HTML?

* Use *k-dimensional tree* for spatial partitioning of model environment.

* <s>Live statistical graphing widgets</s>  (No longer a requirement.)

* Allow end-user to switch between different browser layouts: Visualisation only, Standard <s> and Statistical.</s>

*  Offline ABM computation (start the model running, leave it, reconnect, see what it's up to).

* Batch processing.

* Enable use in a distributed environment.

* Complete testing suite.

* Allow hot-swapping of different `abm` packages.

* Store modelling sessions to server database along with statistical data for later retrieval.

* Fluid ABM timescale controls.

* Optional recording of Visualisation to SVG frame sequence. 