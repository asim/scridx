{{#user}}
  {{> _userNav.m}}

  <h1>Feedback Requests</h1>
  <ul class="inline">
   <li><a href="{{Url}}/feedback?order=latest">Latest</a></li>
   <li><a href="{{Url}}/feedback?order=top">Top</a></li>
   <li><a href="{{Url}}/feedback?order=a-z">A-Z</a></li>
  </ul>

 {{#feedback}}
  {{> _feedback.m}}
 {{/feedback}}
 {{^feedback}}
   <h3>No feedback requested yet</h3>
 {{/feedback}}
{{/user}}
