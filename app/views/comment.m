<h1>Comment</h1>
{{#comment}}
{{#submitter}}
<h4><a href="{{parent_url}}">Re: {{parent_text}}</a></h4>
{{#Thing}}
<div class="comment comment0">
  <div class="row-fluid">
   <span class="pull-right">
    {{> _voteComment.m}}
   </span>
   <span><a href="{{Url}}"><strong>{{Name}}</strong></a><span>
   <div><pre class="text comment-text">{{Text}}</pre></div>
   <div class="font-small pull-right">{{PrintStored}}</div>
  </div>
</div>
{{/Thing}}
{{/submitter}}
<hr/>

<form id="commentForm" action="/submit/comment" method="POST">
 <div class="control-group">
  <div class="controls">
   <textarea name=Text rows="5" class="span12" placeholder="Comment..."></textarea>
  </div>
  <button type=submit class="btn btn-primary btn-mini">Reply</button>
 </div>
 <input type="hidden" name="ParentId" value="{{Id}}">
 <input type="hidden" name="EntityId" value="{{EntityId}}">
 <input type="hidden" name="EntityType" value="{{EntityType}}">
 <input type="hidden" name="_csrf" value="{{csrf}}">
</form>
{{/comment}}

{{#comments}}
  <div class="comment comment{{Depth}} offset{{Depth}}">
   {{> _comment.m}}
  </div>
{{/comments}}
