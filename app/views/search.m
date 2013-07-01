<h1>Search Results</h1>
<h3>"{{query}}"</h3>
<h5>{{total}} record(s) found</h5>
{{#results}}
<div class="row-fluid top-border">
 <span>
   <span class="font-large"><a href="{{Source}}">{{PrintTitle}}</a></span>
   <div class="row-fluid font-small meta-links">
    <div>Writers: {{PrintWriters}}</div>
    <div>Draft Date: {{PrintDraftdate}}, Version: {{PrintVersion}}</div>
    <div><a href="{{Url}}">Comments</a></div>
   </div>
 </span>
</div>
{{/results}}
