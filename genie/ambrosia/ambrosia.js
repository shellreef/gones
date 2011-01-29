// Created:20110127
// By Jeff Connelly

// Ambrosia: an experimental push-style templating engine for JavaScript

// Inspirations: Amrita http://amrita.sourceforge.jp/
// other systems like HTML::Steamstress http://www.perlmonks.org/?node_id=674225
// and for theory http://www.cs.usfca.edu/~parrt/papers/mvc.templates.pdf

function ambrosia(node, value) {
    if (value instanceof Array) {
        // Clone node for each element of array
        for (var i = 0; i < value.length; i += 1) {
            var new_node = node.cloneNode(true); 
            // TODO: make id unique. or should we use classes for everything instead?

            node.parentNode.insertBefore(new_node, node)

            ambrosia(new_node, value[i]);
        }

        // Remove the template node
        node.parentNode.removeChild(node);
    } else if (value === null) {
        // Remove node
        node.parentNode.removeChild(node);
    } else if (value === undefined) {
        // No operation

    // https://developer.mozilla.org/en/JavaScript/Reference/Operators/Special_Operators/typeof_Operator
    // TODO: boolean, xml
    } else if (typeof value === "string" || typeof value === "number") {
        // Scalar text value
        node.textContent = value;
    } else if (typeof value === "function") {
        ambrosia(node, value.call());
    } else if (value instanceof AmbrosiaAttrList) {
        for (attr in value.attributes) { 
            if (value.attributes.hasOwnProperty(attr)) {
                node.setAttribute(attr, value.attributes[attr]);
            }
        }

        ambrosia(node, value.content);
    } else if (typeof value === "object") {
        // Nested
        for (var id in value) {
            if (value.hasOwnProperty(id)) {
                var next_node = node.querySelector("#" + id);
                if (!next_node) {
                    throw "no such node id: " + id + ", from " + next_node;
                }

                ambrosia(next_node, value[id]);
            }
        }
    } else {
        throw "ambrosia(" + node + ", " + value + "): unsupported data type: " + typeof value;
    }
}

function A(attributes, content) {
    if (typeof attributes === "string" || typeof attributes === "number") {
        // Also accept reversed order of parameters, sometimes more convenient
        return new AmbrosiaAttrList(content, attributes);
    }

    return new AmbrosiaAttrList(attributes, content);
}

function AmbrosiaAttrList(attributes, content) {
    this.attributes = attributes;
    this.content = content;
    return this;
}
