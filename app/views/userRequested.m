{{#user}}
  {{> _userNav.m}}
  <h1>Requested</h1>
  <ul class="inline">
   <li><a href="{{Url}}/requested?order=latest">Latest</a></li>
   <li><a href="{{Url}}/requested?order=top">Top</a></li>
   <li><a href="{{Url}}/requested?order=a-z">A-Z</a></li>
  </ul>

 {{#requests}}
  {{> _request.m}}
 {{/requests}}
 {{^requests}}
  <h3>Nothing requested yet</h3>
 {{/requests}}
{{/user}}
