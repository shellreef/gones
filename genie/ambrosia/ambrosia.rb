#!/usr/bin/ruby
# Created:20110129
# By Jeff Connelly

# Ambrosia: an experimental push-style templating engine

require 'rubygems'
require 'nokogiri'


def expand(html, data)
    root = Nokogiri::HTML::DocumentFragment.parse(html)

    expand_node(root, data)

    return root.to_html
end

def expand_node(node, value)
    case value
    when Hash
        value.each do |key, next_value|
            next_node = node.at_css("\##{key}")
            if next_node.nil? 
                throw "expand_node(#{node}, #{value}): no such id: #{key}"
            end
            expand_node(next_node, next_value)
        end
    when String
        node.content = value
    when Fixnum
        node.content = value.to_s
    when Array
        value.each do |item|
            new_node = node.clone       # TODO: uniquify id
            node.parent.add_child(new_node)

            expand_node(new_node, item)
        end
        node.remove
    when NilClass
        node.remove
    else
        throw "expand_node(#{node}, #{value}): unsupported data type: #{value.class}"
    end
end

puts expand(<<HTML, {:x => "Hello, <script>world", :item => [1,2,3], :dead => nil})
<p id=x></p>
<ul>
<li id="item">
</ul>
<span id="dead">This will not appear</span>
HTML

