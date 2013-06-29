{{#user}}
  {{> _userNav.m}}

  <h1>Submitted</h1>
  <ul class="inline">
   <li><a href="{{Url}}/submitted?order=latest">Latest</a></li>
   <li><a href="{{Url}}/submitted?order=top">Top</a></li>
   <li><a href="{{Url}}/submitted?order=a-z">A-Z</a></li>
  </ul>

 {{#scripts}}
  {{> _script.m}}
 {{/scripts}}
 {{^scripts}}
   <h3>Nothing submitted yet</h3>
 {{/scripts}}
{{/user}}
