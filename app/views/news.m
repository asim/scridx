<h1>News & Misc</h1>
<ul class="inline">
  <li><a href="/news?order=top">Top</a></li>
  <li><a href="/news?order=latest">Latest</a></li>
  <li><a href="/news?order=a-z">A-Z</a></li>
</ul>
{{#news}}
 {{> _news.m}}
{{/news}}
{{^news}}
<br><h1>Oh Snap! News? Look <a href="/news">Here</a></h1> 
<br><h2>No news on this page :(</h2>
{{/news}}
