{{#request}}
{{#Thing}}
<div class="row-fluid">
 <div><h1>Request #{{Id}}</h1>
  <div class="row-fluid">
    <div>
     <h4>
      <span class="pull-right"><ul class="inline font-small meta-links">{{{MetaStatus}}}<li>{{> _voteLink.m}}</li></ul></span>
      <a href="{{SourceOrUrl}}">{{PrintTitle}}</a>
     </h4>
<div class="row-fluid font-small">
 <div>Writers: {{PrintWriters}}</div>
 <div>Draft Date: {{PrintDraftdate}}, Version: {{PrintVersion}}</div>
 <div>Requested: {{PrintStored}} {{#submitter}}by <a href="{{Url}}">{{Name}}</a>{{/submitter}}</div>
  {{#fulfiller}}<div>Fulfilled by <a href="{{Url}}">{{Name}}</a></div>{{/fulfiller}}
 <div class="meta-links">
  <span class="things">{{{MetaLinks}}}</span>
 </div>
</div>

    </div>
  </div>
  <div class="row-fluid">
   <div class="top-buffer text span10">{{Logline}}</div>
  </div>
 </div>
</div>
{{/Thing}}
{{/request}}
{{> _comments.m}}

