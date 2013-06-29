<h1>Requests</h1>
<ul class="inline">
  <li><a href="/requests?order=top">Top</a></li>
  <li><a href="/requests?order=latest">Latest</a></li>
  <li><a href="/requests?order=a-z">A-Z</a></li>
</ul>
{{#requests}}
 {{> _request.m}}
{{/requests}}
{{^requests}}
<br><h1>No pending requests</h1>
{{/requests}}
