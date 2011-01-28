// Created:20110127
// By Jeff Connelly

// Ambrosia: an experimental push-style templating engine for JavaScript

// Inspirations: Amrita http://amrita.sourceforge.jp/
// other systems like HTML::Steamstress http://www.perlmonks.org/?node_id=674225
// and for theory http://www.cs.usfca.edu/~parrt/papers/mvc.templates.pdf

function expand(template, data) {
    for (var id in data) {
        if (data.hasOwnProperty(id)) {
            var element = template.getElementById(id);
            if (!element) {
                throw "no such element id: " + id;
            }

            var value = data[id];

            expandValue(element, value);
        }
    }
}

function expandValue(element, value) {
    /* TODO if (value instanceof Array) {
        for (var i = 0; i < value.length; i += 1) {
            var new_node = expandValue(value[i]);

            element.parentNode.insertBefore(new_node, element)
        }
    } else { */
        element.textContent = value;
    //}
}
