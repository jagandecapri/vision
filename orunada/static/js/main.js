$(document).ready(function() {
    console.log("Testing");
    TESTER = document.getElementById('tester');

    Plotly.plot(TESTER, [{
        x: [1, 2, 3, 4, 5],
        y: [1, 2, 4, 8, 16]
    }], {
        margin: {t: 0}
    });

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
        print("RESPONSE: " + evt.data);
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
});