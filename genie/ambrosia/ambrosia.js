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
            // TODO: make id unique

            node.parentNode.insertBefore(new_node, node)

            ambrosia(new_node, value[i]);
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
    } else if (value instanceof AmbrosiaAttrList) {
        for (attr in value.attributes) { 
            if (value.attributes.hasOwnProperty(attr)) {
                node.setAttribute(attr, value.attributes[attr]);
            }
        }

        ambrosia(node, value.value);
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

    // TODO: nesting, very important
    // TODO: attributes, probably through a special object
    } else {
        throw "ambrosia(" + node + ", " + value + "): unsupported data type: " + typeof value;
    }
}

function A(attributes, value) {
    if (typeof attributes === "string" || typeof attributes === "number") {
        // Also accept reversed order of parameters, sometimes more convenient
        return new AmbrosiaAttrList(value, attributes);
    }

    return new AmbrosiaAttrList(attributes, value);
}

function AmbrosiaAttrList(attributes, value) {
    this.attributes = attributes;
    this.value = value;
    return this;
}
