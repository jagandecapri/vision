$(document).ready(function() {

    var __lock = false
    var graphs_cache = []
    var startTime = Date.now();

    // Update timer
    setInterval(function(){
        var elapsedTime = (Date.now() - startTime) / 1000;
        $("#timer").text(elapsedTime.toFixed(3));
    }, 100);

        /* Current Plotly.js version */
    console.log(Plotly.BUILD);

    ws = new WebSocket(ws_uri);
    ws.onopen = function(evt) {
        print("OPEN");
        send()
    }
    ws.onclose = function(evt) {
        print("CLOSE");
        ws = null;
    }

    ws.onmessage = function(evt) {
        print("RESPONSE: ");
        var data = evt.data;
        if (data == 'ping'){
            //TODO: How to handle ping messages from server or is it already handled by the browser?
            console.log("Received 'ping' from server")
        } else {
            try{
                var graphs = processJsonData(data)
                if (__lock == false){
                     if (graphs_cache.length > 0){
                        _.forEach(graphs_cache, function(graphs){
                            while (__lock == true){
                                //No-op blocking
                            }
                            __lock = true;
                            createOrUpdatePlot(graphs).then(function () {
                                __lock = false
                            })
                        });
                        graphs_cache = []
                    } else {
                         while (__lock == true){
                             //No-op blocking
                         }
                         __lock = true;
                         createOrUpdatePlot(graphs).then(function () {
                             __lock = false
                         })
                     }
                } else {
                    graphs_cache.push(graphs)
                }
            }
            catch(error){
                console.log(error, evt.data)
            }
        };
    }

    ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }
    var print = function(message) {
        console.log(message)
    };

    send = function(evt) {
        if (!ws) {
            return false;
        }
        ws.send("pong");
        return false;
    };

    close = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

    var graph_ctr = 0;
    var d3 = Plotly.d3;

    var WIDTH_IN_PERCENT_OF_PARENT = 60,
        HEIGHT_IN_PERCENT_OF_PARENT = 80;

    function createRowDivString(key){
        return '<div class="row"><div class="col-md-4"><div id="' + key + '"></div></div></div>';
    };

    function createDivString(key){
        return '<div class="col-md-4"><div id="' + key + '"></div></div>';
    };

    function isGraphDivExist(key){
        if ($('#'+key).length){
            return true;
        } else {
            return false;
        }
    };

    function appendDivToContainer(key){
        if (graph_ctr%3 == 0){
            var str = createRowDivString(key);
            $("#container").append(str);
        } else {
            var str = createDivString(key);
            $("#container").children().last().append(str);
        }
        graph_ctr++;
    };

    function getKey(graph){
        var key = graph.metadata.id;
        return key;
    }

    function getColumnX(graph){
        var column_x = graph.metadata.column_x;
        return column_x;
    }

    function getColumnY(graph){
        var column_y = graph.metadata.column_y;
        return column_y;
    }

    function processData(graph){
        var points_container = graph.points_container
        var traces = [];
        _.forEach(points_container, function(points){
            var x = [];
            var y = [];
            var point_list = points.point_list
            _.forEach(point_list, function(point){
                x.push(point.data.x)
                y.push(point.data.y)
            });
            trace = {
                type: "scatter",
                mode: "markers",
                x: x,
                y: y,
                marker: {
                    color: points.metadata.color
                }
            }
            traces.push(trace)
        });
        return traces
    }

    function processUpdateDate(graph){
        var points_container = graph.points_container
        var x = [];
        var y = [];
        var colors = [];
        _.forEach(points_container, function(points){
            var point_list = points.point_list
            var x_tmp = []
            var y_tmp = []
            _.forEach(point_list, function(point){
                x_tmp.push(point.data.x)
                y_tmp.push(point.data.y)
            });
            x.push(x_tmp)
            y.push(y_tmp)
            colors.push([points.metadata.color])
        });
        var traces = {
            type: "scatter",
            mode: "markers",
            x: x,
            y: y,
            "marker.color": colors
        }
        return traces
    }

    function processJsonData(raw_data){
        return JSON.parse(raw_data);
    }

    function createNewPlotly(node, data, layout){
        return Plotly.newPlot(node, data, layout)
    }

    function updatePlotly(key, data, layout){
        return Plotly.update(key, data, layout)
    }

    function newPlot(key, graph){
        var x_col = getColumnX(graph);
        var y_col = getColumnY(graph);
        var traces = processData(graph);
        appendDivToContainer(key);

        var gd3 = d3.select('#'+key)
            .style({
                width: WIDTH_IN_PERCENT_OF_PARENT + '%',
                height: HEIGHT_IN_PERCENT_OF_PARENT + '%',
                'margin-top': 0,
                'margin-left': 0,
                'margin-bottom': 0,
                'margin-right': 0
            });

        var layout = {
            title: key,
            xaxis: {
                title: x_col,
                range: [-0.05, 1.05],
                tick0: 0,
                dtick: 0.05
            },
            yaxis: {
                title: y_col,
                range: [-0.05, 1.05],
                tick0: 0,
                dtick: 0.05
            }
        };

        var node = gd3.node();
        var promise = createNewPlotly(node, traces, layout);
        return promise;
    }

    function updatePlot(key, graph){
        var x_col = getColumnX(graph);
        var y_col = getColumnY(graph);
        var traces = processUpdateDate(graph);

        var layout = {
            title: key,
            xaxis: {
                title: x_col,
                range: [-0.05, 1.05],
                tick0: 0,
                dtick: 0.05
            },
            yaxis: {
                title: y_col,
                range: [-0.05, 1.05],
                tick0: 0,
                dtick: 0.05
            }
        };

        var promise = updatePlotly(key, traces, layout);
        return promise;
    }

    function createOrUpdatePlot(graphs){
        var promises = []
        _.forEach(graphs, function(graph){
            var key = getKey(graph)
            var promise;
            if (!isGraphDivExist(key)){
                promise = newPlot(key, graph)
            } else {
                promise = updatePlot(key, graph)
            }
            promises.push(promise)
        });
        var ret_promise = $.when(promises)
        return ret_promise
    }
});