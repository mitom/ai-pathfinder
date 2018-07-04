var states =[];
var graph;
var current = 0;
var found = false;
var goal = 0;
var playing = false;
var player;

function addState(current, path, open_list, closed_list) {
    states.push({
        "current": current,
        "path": path.split('->'),
        "open": open_list,
        "closed": closed_list
    })

    if (current.i == goal) {
        found = true
    }
}

function setCave(cave) {
    var g = {
        'nodes': [],
        'edges': []
    };

    for (var i in cave.nodes) {
        var node = cave.nodes[i];
        g.nodes.push({
            "id": 'n' + node.i,
            "cid": node.i,
            "x": node.x,
            "y": node.y,
            "label": node.i + "(" + node.x + "," + node.y + ")",
            "size":10
        });

        goal = Math.max(goal, node.i)
    }

    for (var n in cave.edges) {
        var edge = cave.edges[n];
        g.edges.push({
            "id": "e"+edge.f+"-"+edge.t,
            "source" : "n"+edge.f,
            "target": "n"+edge.t,
            "from": edge.f,
            "to": edge.t,
            "color": edge.f == goal ? colours["pristine"] : null,
        })
    }

    var c = g.nodes.length;
    var labelThreshold = 8;
    if (c > 40) {
        labelThreshold = 20
    } else if (c > 30) {
        labelThreshold = 15
    } else if (c> 15) {
        labelThreshold = 10
    }

    graph = new sigma({
        graph: g,
        renderers: [
            {
                container: document.getElementById('container'),
                type: 'canvas'
            }
        ],


        settings: {
            defaultNodeColor: colours["pristine"],
            defaultEdgeColor: colours["pristine_light"],
            edgeColor: "default",
            defaultEdgeType: "arrow",
            minArrowSize: 5,
            labelThreshold: labelThreshold
        }
    });

}


var colours = {
    "pristine_light": "#91918f",
    "pristine": "#626262",
    "current": "#117ccd",
    "closed": "#aa322a",
    "open": "#aaa90e",
    "success": "#1daa14",
    "goal": "#803faa"
};

function playAction() {
    player = setTimeout(function() {
        if (current == states.length-1) {
            play();
        } else {
            step(1);
            playAction()
        }
    }, 100)
}

function play() {
    if (playing && player) {
        playing = false;
        document.getElementById('play').classList.remove('hidden')
        document.getElementById('pause').classList.add('hidden')
        clearTimeout(player)
    } else {
        playing = true;
        document.getElementById('play').classList.add('hidden')
        document.getElementById('pause').classList.remove('hidden')
        playAction()
    }
}

function reset() {
    playing = false;
    current = 0;
    load(states[0]);
    step(0)
}

function step(dir) {
    if (current + dir < 0 || current+dir > states.length-1) {
        return
    }

    if (current + dir != current) {
        current = current + dir;
        load(states[current]);
    }

    if (current == 0) {
        document.getElementById('back').classList.add('disabled')
        document.getElementById('reset').classList.add('disabled')
    } else {
        document.getElementById('back').classList.remove('disabled')
        document.getElementById('reset').classList.remove('disabled')
    }

    if (current == states.length-1) {
        document.getElementById('forward').classList.add('disabled')
        document.getElementById('play').classList.add('disabled')
    } else {
        document.getElementById('forward').classList.remove('disabled')
        document.getElementById('play').classList.remove('disabled')
    }
}

function init() {
    load(states[0])
    step(0)
}

function in_list(id, list) {
    for (var i in list) {
        if (list[i].i == id) {
            return list[i]
        }
    }

    return null
}

function load(state) {
    var nodes = graph.graph.nodes();
    for (var n in nodes) {
        var node = nodes[n];
        var o;
        node.label = node.cid + "(" + node.x + "," + node.y + ")";
        // debugger;
        if (node.cid == goal) {
            node["color"] = colours["goal"];
        } else if (node.cid == state.current.i) {
            node["color"] = colours["current"]
        } else if (state.path.includes(String(node.cid))) {
            if (current == states.length-1) {
                if (found) {
                    node["color"] = colours["success"]
                } else {
                    node["color"] = colours["closed"]
                }
            } else {
                node["color"] = colours["current"]
            }
        } else if (o = in_list(node.cid, state.open)) {
            node["color"] = colours["open"];
            node.label = node.cid + "(" + node.x + "," + node.y + ")[" + o.c.toFixed(3) + "]"
        } else if (o = in_list(node.cid, state.closed)) {
            node["color"] = colours["closed"]
        } else {
            delete node["color"]
        }
    }

    var edges = graph.graph.edges();

    for (var n in edges) {
        var edge = edges[n];
        var index = state.path.indexOf(String(edge.from));

        if (index != -1 && edge.to == state.path[index+1]) {
            if (current == states.length-1) {
                if (found) {
                    edge["color"] = colours["success"]
                } else {
                    edge["color"] = colours["closed"]
                }
            } else {
                edge["color"] = colours["current"]
            }
        } else {
            delete edge["color"]
        }

    }

    graph.refresh()
}


document.onkeydown = function (e) {
    e = e || window.event;
    switch (e.which || e.keyCode) {
        case 37:
            step(-1);
            break;
        case 39:
            step(1);
            break;
        case 32:
            play();
            break;
    }
}