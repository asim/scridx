<h1>Comments</h1>
<ul class="inline">
  <li><a href="/comments?order=latest">Latest</a></li>
  <li><a href="/comments?order=top">Top</a></li>
</ul>

{{#comments}}
 <div class="comment comment0">
  {{> _comment.m}}
 </div>
{{/comments}}
{{^comments}}
<br><h1>Oh Snap! Comments? Look <a href="/comments">Here</a></h1> 
<br><h2>No comments on this page :(</h2>
{{/comments}}
