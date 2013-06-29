{{#news}}
{{#Thing}}
<div class="row-fluid">
 <div><h1>News & Misc</h1>
  <div class="row-fluid">
    <div><h4><a href="{{SourceOrUrl}}">{{PrintTitle}}</a><span class="pull-right">{{> _voteLink.m}}</span></h4>
     <div class="row-fluid font-small">
       <div>Submitted {{PrintStored}} {{#submitter}}by <a href="{{Url}}">{{Name}}</a>{{/submitter}}</div>
     </div>
    </div>
  </div>
  <div class="row-fluid">
   <pre class="text">{{Text}}</pre>
  </div>
 </div>
</div>
{{/Thing}}
{{/news}}
{{> _comments.m}}
