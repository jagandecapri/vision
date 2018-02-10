$(document).ready(function() {
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
    var throttled = _.throttle(createOrUpdatePlot, 200);
    ws.onmessage = function(evt) {
        print("RESPONSE: ");
        var data = evt.data;
        if (evt.data == 'ping'){
            //TODO: How to handle ping messages from server or is it already handled by the browser?
            console.log("Received 'ping' from server")
        } else {
            throttled(data)
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
            _.forEach(points, function(points_list){
                var x = [];
                var y = [];
                _.forEach(points_list.data, function(data){
                    x.push(data.x)
                    y.push(data.y)
                })
                trace.x = x
                trace.y = y
                trace.color = points.metadata.color
            });
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

    function updatePlotly(key, data){
        Plotly.restyle(key, data)
    }

    function newPlot(key, graph){
        var x_col = getColumnX(graph);
        var y_col = getColumnY(graph);
        var traces = processData(graph)

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

        var trace1 = {
            x: x,
            y: y,
            mode: 'markers',
            type: 'scatter'
        };

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

    function createOrUpdatePlot(data){
        var tmp = processJsonData(data);
        var graphs = tmp.Data;
        _.forEach(graphs, function(graph){
            var key = getKey(graph)
            if (!isGraphDivExist(key)){
                newPlot(key, data)
            } else {
                updatePlot(key, data)
            }
        })

    }
});