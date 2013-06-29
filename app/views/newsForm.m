<div class="form-center">
<h1>Submit Content</h1>
<form id="itemForm" action="/submit/news" method="POST">
 <fieldset>
  <legend>Details</legend>
  {{#form}}

  <div class="control-group">
   <label class="control-label">Title</label>
   <div class="controls">
    <input class="span12" type=text name="Title" placeholder="Title" value="{{Title}}" autofocus>  
   </div>
  </div>

   <div class="control-group">
    <label class="control-label">Url</label>
    <div class="controls">
     <input class="span12" type=url name="Source" placeholder="Url" value="{{Source}}"> 
    </div>
   </div>
   <div class="control-group">
    <label class="control-label">Text</label>
    <div class="controls">
     <textarea name=Text rows="5" class="span12" placeholder="Text...">{{Text}}</textarea> <br>
    </div>
   </div>
   <button type=submit class="btn btn-primary">Submit</button>
  {{/form}}
 </fieldset>
 <input type="hidden" name="_csrf" value="{{csrf}}">
</form>
</div>
