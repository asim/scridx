<hr/>
<form id="commentForm" action="/submit/comment" method="POST">
 <div class="control-group">
  <div class="controls">
   <textarea name=Text rows="5" class="span12" placeholder="Comment..."></textarea>
  </div>
  <button type=submit class="btn btn-primary btn-mini">Respond</button>
 </div>
 {{#comment}}
 <input type="hidden" name="ParentId" value="{{ParentId}}">
 <input type="hidden" name="EntityId" value="{{EntityId}}">
 <input type="hidden" name="EntityType" value="{{EntityType}}">
 {{/comment}}
 <input type="hidden" name="_csrf" value="{{csrf}}">
</form>
