#!/usr/bin/ruby
require 'ambrosia'

data = {
    :x => "Hello, <script>world", 
    :item => [1,2,3], 
    :dead => nil, 
    :dead2 => false,
    :alive => true,
    :link => A({:href => "http://example.com/"}, "example link"),
    :link2 => A("another link", {:href => "http://example.com/"}),
    :logo => A({:src => "http://upload.wikimedia.org/wikipedia/commons/3/3c/Ambrosia_salad.jpg"}),
    }

puts ambrosia(<<HTML, data)
<p id=x></p>
<ul>
<li id="item">
</ul>
<span id="dead">This will not appear</span><span id="dead2">This either</span><span id="alive">But this will</span>

<a id="link"></a>, <a id="link2"></a>

<img id="logo">
HTML

