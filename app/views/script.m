{{#script}}
{{#Thing}}
<div class="row-fluid">
 <div><h1>Script</h1>
  <div class="row-fluid">
   <div><h4><a href="{{Source}}">{{PrintTitle}}</a><span class="pull-right">{{> _voteLink.m}}</span></h4>
<div class="row-fluid font-small">
 <div class="span8">
  <div>Writers: {{PrintWriters}}</div>
  <div>Draft Date: {{PrintDraftdate}}, Version: {{PrintVersion}}</div>
  <div>Submitted {{PrintStored}} {{#submitter}}by <a href="{{Url}}">{{Name}}</a>{{/submitter}}</div>
 </div>
 <div class="span3 pull-right things">{{{MetaLinks}}}</div>
</div>
   </div>
  </div>
  <div class="row-fluid">
   <div class="span10 text">{{Logline}}</div>
  </div>
 </div>
</div>
{{/Thing}}
{{/script}}

{{> _comments.m}}
