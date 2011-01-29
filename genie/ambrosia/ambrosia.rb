#!/usr/bin/ruby
# Created:20110129
# By Jeff Connelly

# Ambrosia: an experimental push-style templating engine

require 'rubygems'
require 'nokogiri'


def expand(html, data)
    root = Nokogiri::HTML::DocumentFragment.parse(html)

    expandNode(root, data)

    return root.to_html
end

def expandNode(node, value)
    case value
    when Hash
        value.each do |key, next_value|
            next_node = node.at_css("\##{key}")
            if next_node.nil? 
                throw "expandNode(#{node}, #{value}): no such id: #{key}"
            end
            expandNode(next_node, next_value)
        end
    when String
        node.content = value
    else
        throw "expandNode(#{node}, #{value}): unsupported data type: #{value.class}"
    end
end

puts expand("<p id=x></p>", {:x => "Hello, <script>world"})
