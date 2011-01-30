#!/usr/bin/ruby
require 'ambrosia'

# Example of to_s
class Person
    attr_accessor :first, :last, :aka

    def initialize(first, last, aka)
        @first = first
        @last = last
        @aka = aka
    end

    def to_s
        return "#{last}, #{first} (also known as: #{aka.join(' / ')})"
    end
end

# Example of to_node
class Technology 
    attr_accessor :name, :language, :url

    def initialize(name, language, url)
        @name = name
        @language = language
        @url = url
    end

    def to_node(root)
        node = Nokogiri::XML::Node.new("span", root)
        node.content = name
    end
end

data = {
    :x => "Hello, <script>world", 
    :item => [1,2,3], 
    :dead => nil, 
    :dead2 => false,
    :alive => true,
    :link => A({:href => "http://example.com/"}, "example link"),
    :link2 => A("another link", {:href => "http://example.com/"}),
    :logo => A({:src => "http://upload.wikimedia.org/wikipedia/commons/3/3c/Ambrosia_salad.jpg"}),
    :person => Person.new("Yukihiro", "Matsumoto", ["Matz", "松本行弘", "まつもとゆきひろ"]),
    :system => [
        Technology.new("Ambrosia", "Ruby", nil),
        Technology.new("Amrita", "Ruby", "http://amrita.sourceforge.jp/"),
        Technology.new("HTML::Seamstress", "Perl", "http://search.cpan.org/~tbone/HTML-Seamstress-5.0h/lib/HTML/Seamstress.pod"),
        Technology.new("All push-style", "Various", "http://www.perlmonks.org/?node_id=674225"),
        ]
    }

puts ambrosia(<<HTML, data)
<p id=x></p>
<ul>
<li id="item">
</ul>
<span id="dead">This will not appear</span><span id="dead2">This either</span><span id="alive">But this will</span>

<a id="link"></a>, <a id="link2"></a>

<img id="logo">

<hr>

<p>Ruby was created by <u id="person"></u>

<p>A few templating systems:
<div id="system"></div>

HTML

