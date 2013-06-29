<h1>Scripts</h1>
<ul class="inline">
  <li><a href="/scripts?order=top">Top</a></li>
  <li><a href="/scripts?order=latest">Latest</a></li>
  <li><a href="/scripts?order=a-z">A-Z</a></li>
</ul>
{{#scripts}}
 {{> _script.m}}
{{/scripts}}
{{^scripts}}
<br><h1>Oh Snap! Scripts? Look <a href="/scripts">Here</a></h1> 
<br><h2>There are no scripts on this page :(</h2>
{{/scripts}}
