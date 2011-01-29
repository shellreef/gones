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
    when AmbrosiaAttrList
        value.attributes.each do |attribute_name, attribute_value|
            node[attribute_name.to_s] = attribute_value 
        end
        expand_node(node, value.content)
    else
        throw "expand_node(#{node}, #{value}): unsupported data type: #{value.class}"
    end
end

def A(attributes, content)
    # TODO: support reverse order too
    return AmbrosiaAttrList.new(attributes, content)
end

class AmbrosiaAttrList
    attr_accessor :attributes, :content
    def initialize(attributes, content)
        @attributes = attributes
        @content = content
    end
end

puts expand(<<HTML, {:x => "Hello, <script>world", :item => [1,2,3], :dead => nil, :link => A({:href => "http://example.com/"}, "example link")})
<p id=x></p>
<ul>
<li id="item">
</ul>
<span id="dead">This will not appear</span>

<a id="link"></a>
HTML

