// Created:20110127
// By Jeff Connelly

// Ambrosia: an experimental push-style templating engine for JavaScript

// Inspirations: Amrita http://amrita.sourceforge.jp/
// other systems like HTML::Steamstress http://www.perlmonks.org/?node_id=674225
// and for theory http://www.cs.usfca.edu/~parrt/papers/mvc.templates.pdf

function expand(root, data) {
    for (var id in data) {
        if (data.hasOwnProperty(id)) {
            var node = root.querySelector("#" + id);
            if (!node) {
                console.log(root);
                throw "no such node id: " + id + ", from " + root;
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
            var new_node = node.cloneNode(true); 
            // TODO: make id unique

            node.parentNode.insertBefore(new_node, node)

            expandValue(new_node, value[i]);
        }

        // Remove the template node
        node.parentNode.removeChild(node);
    } else if (value === null) {
        // Remove node
        node.parentNode.removeChild(node);
    // https://developer.mozilla.org/en/JavaScript/Reference/Operators/Special_Operators/typeof_Operator
    // TODO: undefined, boolean, function, xml
    } else if (typeof value === "string" || typeof value === "number") {
        // Scalar text value
        node.textContent = value;
    } else if (typeof value === "object") {
        // Nested
        expand(node, value);
    // TODO: nesting, very important
    // TODO: attributes, probably through a special object
    } else {
        throw "expandValue(" + node + ", " + value + "): unsupported data type: " + typeof value;
    }
}
