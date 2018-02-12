$(document).ready(function() {
        /* Current Plotly.js version */
    console.log(Plotly.BUILD);

    var graphs_cache = []
    ws = new WebSocket(ws_uri);
    ws.onopen = function(evt) {
        print("OPEN");
        (function interval(graphs_cache){
            if (graphs_cache.length > 0){
                _.forEach(graphs_cache, function(graphs){
                    createOrUpdatePlot(graphs)
                })
                graphs_cache = []
            }
            _.delay(interval, 1000, graphs_cache)
        })(graphs_cache)
        send()
    }
    ws.onclose = function(evt) {
        print("CLOSE");
        ws = null;
    }

    var throttled = _.throttle(createOrUpdatePlot, 200);
    var graphs_cache = []

    ws.onmessage = function(evt) {
        print("RESPONSE: ");
        var data = evt.data;
        if (data == 'ping'){
            //TODO: How to handle ping messages from server or is it already handled by the browser?
            console.log("Received 'ping' from server")
        } else {
            try{
                var graphs = processJsonData(data)
                graphs_cache.push(graphs)
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
            var trace = {};
            var x = [];
            var y = [];
            var point_list = points.point_list
            _.forEach(point_list, function(point){
                x.push(point.data.x)
                y.push(point.data.y)
            });
            trace.x = x
            trace.y = y
            trace.color = points.metadata.color
            traces.push(trace)
        });
        return traces
    }

    function processJsonData(raw_data){
        return JSON.parse(raw_data);
    }

    function createNewPlotly(node, data, layout){
        Plotly.newPlot(node, data, layout);
    }

    function updatePlotly(key, traces){
        Plotly.update(key, traces)
    }

    function newPlot(key, x_col, y_col, traces){

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

        var node = gd3.node()
        createNewPlotly(node, traces, layout);
    }

    function updatePlot(key, graph){
        var traces = processData(graph)
        updatePlotly(key, traces);
    }

    function createOrUpdatePlot(graphs){
        _.forEach(graphs, function(graph){
            var key = getKey(graph)
            if (!isGraphDivExist(key)){
                var x_col = getColumnX(graph);
                var y_col = getColumnY(graph);
                var traces = processData(graph);
                newPlot(key, x_col, y_col, traces)
            } else {
                var traces = processData(graph);
                updatePlot(key, traces)
            }
        })

    }
});