var windowsIDmap=[];
function RePlot(plotcmd) {
}
function NewPlot(plotcmd) {
	options=plotcmd.Options;
	newid="W"+plotcmd.Handle.toString();
	console.log(windowsIDmap,newid);
	var holdon=plotcmd.HoldOn;

	if (windowsIDmap[newid]== undefined){
//Create a New Window

newwindow=window.open(url="/figure.html",options.Title,"location:0,menubar=no,status=yes,toolbar=no,titlebar=no");
windowsIDmap[newid]={Window:newwindow,Chart:null};
newwindow.addEventListener("resize", windowResized);
function windowResized(obj){
	d=obj.currentTarget.document;
	f=d.getElementById("figurebox");
	wid=f.dataset.wid;
	console.log("I am resized ",obj.currentTarget.name, "ID is ",wid);
	if (windowsIDmap[wid]!=undefined ){
		ch=windowsIDmap[wid].Chart;
		ch.show();
	}
}
newwindow.onbeforeunload = function(){
// alert('Have to Close this '+newid);
windowsIDmap[newid]=null;
// set warning message
};
newwindow.onclose = function()
{
// alert('Have to Close this '+newid);
windowsIDmap[newid]=null;
}

newwindow.onload=function(obj) {
// alert("Loaded Window ");
win=obj.currentTarget;
d=win.document;
d.title=plotcmd.Options.Title;
console.log("Document element",d);
figure=d.getElementById("figurebox");
// "W"+plotcmd.Handle.toString();
console.log("new id avaialbe ",newid)
figure.setAttribute("data-wid",newid);
chartobj=PlotOnFigure(figure,plotcmd);
windowsIDmap[newid]={Window:win,Chart:chartobj};
console.log("Updated Entry ",windowsIDmap);
}
// windowsIDmap[newid]={Window:newwindow,Chart:nil};
}else {
holdon=plotcmd.HoldOn; // ?? May be or rewrite window
console.log("Already Exists Holon = ??",holdon);
wobj=windowsIDmap[newid];
cwindow=wobj.Window;
d=cwindow.document;
d.title=plotcmd.Options.Title;
console.log("Document element",d);
figure=d.getElementById("figurebox");
if (holdon) {
	chartobj=wobj.Chart;
	if (chartobj==null || chartobj==undefined) {
		chartobj=PlotOnFigure(figure,plotcmd);
		wobj.Chart=chartobj;
		windowsIDmap[newid]=wobj;
	}else		{
		chartobj=RePlotOnFigure(figure,plotcmd,chartobj);
		wobj.Chart=chartobj;
		windowsIDmap[newid]=wobj;
	}
}else {
	chartobj=PlotOnFigure(figure,plotcmd);
	wobj.Chart=chartobj;
	windowsIDmap[newid]=wobj;
}
}
}
function RePlotOnFigure(element,plotcmd,chartobj) {
	ycol=plotcmd.Y;
	label='y1'+Math.random();
	// console.log("Can I replot on this ",chartobj,label);
	ycol.unshift(label);
	

if (plotcmd.X!=null)
{
	console.log("Received more");
	xcol=plotcmd.X;
	xcol.unshift('xdata');
	dataobj=[xcol,ycol];
 
	chartobj.load({columns:dataobj});
return chartobj;
}else{
	dataobj=[ycol];
chartobj.load({columns:dataobj});	
}

return chartobj;
}


function PlotOnFigure(element,plotcmd) {
	ycol=plotcmd.Y;
	ycol.unshift('y1');
	console.log("Check plot marker ",plotcmd.Options.Marker);
	if (plotcmd.Options.Marker==""){
		ptype="line";
	}else{
 	console.log("Running scatter plot ");
		ptype="scatter";
	}
	colorfn=function(c,d){
// console.log("Original color ",c,"data ",d);
if (plotcmd.Options.Color=="")
	return c;
else
	return plotcmd.Options.Color;
}
if (plotcmd.X!=null)
{
	console.log("Received more");
	xcol=plotcmd.X;
	xcol.unshift('xdata');
	dataobj={x:'xdata',columns:[xcol,ycol],type:ptype,color:colorfn};
	chart=c3.generate({
		bindto:element,		
		data:dataobj,
		axis:axisobj,
		grid:{x:{show:true},y:{show:true}},
		zoom:{enabled:true}

	});
}
else {
	dataobj={columns:[ycol],type:ptype,color:colorfn};
	chart=c3.generate({
		bindto:element,
		data:dataobj,
		axis:axisobj,
		grid:{x:{show:true},y:{show:true}},
		zoom:{enabled:true}

	});
}

// xcol=xval;
// xcol.unshift('t');
// dataobj={columns:[xcol,ycol],x:'t'};
// label={}

return chart;
}
function PlotSineWave(xval,yval,options) {
	newwindow=window.open(url="/figure.html","","location=no,menubar=no,status=no,toolbar=no,titlebar=no,height=440");
// newwindow=window.open("/figure.html",'name','toolbar=1,scrollbars=1,location=1,statusbar=0,menubar=1,resizable=1,width=800,height=600');
var chart;
newwindow.onload=function(obj) {
// alert("Loaded Window ");
console.log(obj)
win=obj.currentTarget;
d=win.document;
console.log("Document element",d);
figure=d.getElementById("figurebox");
console.log("Figure element",figure);
// console.log("plot",xvals,yvals);
xcol=xval;
ycol=yval;
xcol.unshift('t');
ycol.unshift('sineWave');
dataobj={columns:[xcol,ycol],x:'t'};
// label={}
console.log(dataobj);
chart=c3.generate({
	bindto:figure,
	data:dataobj,
	axis:axisobj,
	grid:{x:{show:true},y:{show:true}},
	zoom:{enabled:true}
}
);
}
return chart;
}