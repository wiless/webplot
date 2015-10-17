# webplot
Matlab like interface for plotting, uses c3js

![Screenshot](http://github.com/wiless/webplot/blob/refs/images/snapshot.jpg)

## Sample output figure
![Screenshot](https://github.com/wiless/webplot/blob/refs/images/snapshot1.jpg)


**Installation**
```
go get github.com/wiless/webplot
```

**Start Session**
- Start the Server

- Open the dashboard


**Run the example**
##Usage 
```
s := wm.NewSession("HETNET")
s.Plot(vlib.RandUFVec(10), "holdon", "title=CDF Plot of received signal", "LineWidth=2")
```

![Screenshot](https://github.com/wiless/webplot/blob/refs/images/snapshot.jpg)

## Sample output figure
![Screenshot2](https://github.com/wiless/webplot/blob/refs/images/snapshot1.jpg)


