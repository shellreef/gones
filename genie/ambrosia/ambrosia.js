// Created:20110127
// By Jeff Connelly

// Ambrosia: an experimental push-style templating engine for JavaScript

// Inspirations: Amrita http://amrita.sourceforge.jp/
// other systems like HTML::Steamstress http://www.perlmonks.org/?node_id=674225
// and for theory http://www.cs.usfca.edu/~parrt/papers/mvc.templates.pdf

function expand(template, data) {
    for (var id in data) {
        if (data.hasOwnProperty(id)) {
            var node = template.getElementById(id);
            if (!node) {
                throw "no such node id: " + id;
            }

            var value = data[id];

            expandValue(node, value);
        }
    }
}

function expandValue(node, value) {
    if (value instanceof Array) {
        // Clone node for each element of array
        for (var i = 0; i < value.length; i += 1) {
            var new_node = node.cloneNode();

            node.parentNode.insertBefore(new_node, node)

            expandValue(new_node, value[i]);
        }

        // Remove the template node
        node.parentNode.removeChild(node);
    } else if (value === null) {
        // Remove node
        node.parentNode.removeChild(node);
    } else { 
        node.textContent = value;
    }
    // TODO: nesting, very important
    // TODO: attributes, probably through a special object
}
