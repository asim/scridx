{{#user}}
  {{> _userNav.m}}

  <h1>News & Misc</h1>
  <ul class="inline">
   <li><a href="{{Url}}/news?order=latest">Latest</a></li>
   <li><a href="{{Url}}/news?order=top">Top</a></li>
   <li><a href="{{Url}}/news?order=a-z">A-Z</a></li>
  </ul>

 {{#news}}
  {{> _news.m}}
 {{/news}}
 {{^news}}
   <h3>Nothing submitted yet</h3>
 {{/news}}
{{/user}}
