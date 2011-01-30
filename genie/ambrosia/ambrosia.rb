#!/usr/bin/ruby
# Created:20110129
# By Jeff Connelly

# Ambrosia: an experimental push-style templating engine

require 'rubygems'
require 'nokogiri'


def ambrosia(html, data)
    root = Nokogiri::HTML::DocumentFragment.parse(html)

    root.ambrosia(data)

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
            next_node.ambrosia(next_value)
        end
    when String
        node.content = value
    when Array
        value.each do |item|
            new_node = node.clone       # TODO: uniquify id
            node.parent.add_child(new_node)

            new_node.ambrosia(item)
        end
        node.remove
    when NilClass, FalseClass
        node.remove
    when TrueClass
        # no operation
    when AmbrosiaAttrList
        value.attributes.each do |attribute_name, attribute_value|
            node[attribute_name.to_s] = attribute_value 
        end
        node.ambrosia(value.content)
    else
        # Not a recognized type, so try by what it responds to
        if value.respond_to? :to_node
            value.to_node(node)
        elsif value.respond_to? :to_s
            node.content = value.to_s
        else
            throw "expand_node(#{node}, #{value}): unsupported data type: #{value.class}"
        end
    end
end

def A(attributes, content=true)
    if attributes.class == String || attributes.class == Fixnum
        # reversed order
        return AmbrosiaAttrList.new(content, attributes)
    end

    return AmbrosiaAttrList.new(attributes, content)
end

class AmbrosiaAttrList
    attr_accessor :attributes, :content
    def initialize(attributes, content)
        @attributes = attributes
        @content = content
    end
end


# Convenience methods to make it more OO
class Nokogiri::HTML::DocumentFragment
    def ambrosia(value)
        expand_node(self, value)
    end
end

class Nokogiri::XML::Element
    def ambrosia(value)
        expand_node(self, value)
    end
end


