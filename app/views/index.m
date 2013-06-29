<!--
<div class="jumbotron">
<h1>Welcome to Scridx!</h1>
<p class="lead">
  Looking for a screenplay? Written one and want some feedback? Scridx is dedicated to 
  collecting all the great scripts out there, with your help. Signup and contribute.
<p>
</div>
-->
<h1>Scripts</h1>
<ul class="inline">
  <li><a href="/?order=top">Top</a></li>
  <li><a href="/?order=latest">Latest</a></li>
  <li><a href="/?order=a-z">A-Z</a></li>
</ul>
{{#scripts}}
{{> _script.m}}
{{/scripts}}
{{^scripts}}
<h4>We seem to be missing something. Uh Reload?</h4>
{{/scripts}}
