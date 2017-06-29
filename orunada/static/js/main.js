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
    ws.onmessage = function(evt) {
        //print("RESPONSE: " + evt.data);
        print("RESPONSE: ");
        var data = evt.data;
        createOrUpdatePlot(data)
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
        ws.send("ping");
        return false;
    };

    close = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

    function createDivString(key){
        return '<div class="row">'+'<div id="' + key + '"></div></div>';
    };

    function isGraphDivExist(key){
        if ($('#'+key).length){
            return true;
        } else {
            return false;
        }
    };

    function appendDivToContainer(container_div_key, append_div){
        $('#'+container_div_key).append(append_div);
    };

    function getColumnX(data){
        var columns = _.first(data).Sorter;
        var col_x = columns[0];
        return col_x
    }

    function getColumnY(data){
        var columns = _.first(data).Sorter;
        var col_y = columns[1];
        return col_y;
    }

    function getKey(data){
        var columns = _.first(data).Sorter;
        var key = columns.join('-');
        return key;
    }

    function processData(data, col_x, col_y){
        var x = [];
        var y = [];
        _.forEach(data, function(val, key){
            x.push(val["Norm_vec"][col_x])
            y.push(val["Norm_vec"][col_y])
        });
        return [x,y]
    }

    function processJsonData(raw_data){
        return JSON.parse(raw_data);
    }

    function createNewPlotly(key, data, layout){
        Plotly.newPlot(key, data, layout);
    }

    function updatePlotly(key, data){
        Plotly.restyle(key, data)
    }

    function newPlot(key, data){
        var x_col = getColumnX(data);
        var y_col = getColumnY(data);
        var tmp = processData(data,x_col,y_col)
        var x = tmp[0];
        var y = tmp[1];

        var div = createDivString(key);
        appendDivToContainer("container", div);

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
                range: [0, 1],
                tick0: 0,
                dtick: 0.05
            },
            yaxis: {
                title: y_col,
                range: [0, 1],
                tick0: 0,
                dtick: 0.05
            }
        };

        var data = [trace1];
        createNewPlotly(key, data, layout);
    }

    function updatePlot(key, data){
        var x_col = getColumnX(data);
        var y_col = getColumnY(data);
        var tmp = processData(data,x_col,y_col)
        var x = tmp[0];
        var y = tmp[1];

        var update_data = {
            x: [x],
            y: [y]
        };

        updatePlotly(key, update_data);
    }

    function createOrUpdatePlot(data){
        var tmp = processJsonData(data);
        var data = tmp.Data;
        var key = getKey(data)
        if (!isGraphDivExist(key)){
            newPlot(key, data)
        } else {
            updatePlot(key, data)
        }
    }
});