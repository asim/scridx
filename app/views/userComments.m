{{#user}}
  {{> _userNav.m}}

  <h1>Comments</h1>
  <ul class="inline">
   <li><a href="{{Url}}/comments?order=latest">Latest</a></li>
   <li><a href="{{Url}}/comments?order=top">Top</a></li>
  </ul>
 {{#comments}}
  <div class="comment comment0">
   {{> _comment.m}}
  </div>
 {{/comments}}
 {{^comments}}
  <h3>No comments made yet</h3>
 {{/comments}}
{{/user}}
